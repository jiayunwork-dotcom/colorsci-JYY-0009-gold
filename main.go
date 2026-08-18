// Command colorsci is a color science toolkit: parse CSS colors, convert
// between color spaces and compute color difference / contrast metrics.
//
// Subcommands:
//
//	colorsci parse <color>            parse and re-format a CSS color
//	colorsci convert <color> -to s   convert to sRGB/hex/xyz/lab/lch/hsl
//	colorsci deltae <c1> <c2> -m 2000  color difference (76/94/2000)
//	colorsci contrast <c1> <c2>       WCAG contrast ratio and levels
//	colorsci names [-contains s]      list supported named colors
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"colorsci/internal/colorconv"
	"colorsci/internal/colormetric"
	"colorsci/internal/colorparse"
	"colorsci/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		usage(os.Stderr)
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "parse":
		err = runParse(os.Args[2:])
	case "convert":
		err = runConvert(os.Args[2:])
	case "deltae":
		err = runDeltaE(os.Args[2:])
	case "contrast":
		err = runContrast(os.Args[2:])
	case "names":
		err = runNames(os.Args[2:])
	case "serve":
		addr := ":8080"
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
		fmt.Printf("Starting color science server on %s ...\n", addr)
		srv := server.New(addr)
		err = srv.ListenAndServe()
	case "help", "-h", "--help":
		usage(os.Stdout)
		return
	default:
		fmt.Fprintf(os.Stderr, "colorsci: unknown subcommand %q\n\n", os.Args[1])
		usage(os.Stderr)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "colorsci:", err)
		os.Exit(1)
	}
}

func usage(w *os.File) {
	fmt.Fprintln(w, `usage: colorsci <subcommand> [flags] <args>

subcommands:
  parse <color>                parse and re-format a CSS color
  convert <color> -to <space>  convert to sRGB/hex/xyz/lab/lch/hsl
  deltae <c1> <c2> -m <76|94|2000>  CIE color difference
  contrast <c1> <c2>           WCAG contrast ratio and levels
  names [-contains s]          list supported named colors

colors accept hex, rgb()/hsl() functions and W3C names.`)
}

// reorderFlags moves flags (and their values) ahead of positional
// arguments so flag.FlagSet.Parse sees them first. Go's flag package
// stops at the first non-flag argument.
func reorderFlags(args []string, boolFlags map[string]bool) []string {
	flags := make([]string, 0, len(args))
	rest := make([]string, 0, len(args))
	i := 0
	for i < len(args) {
		a := args[i]
		if !strings.HasPrefix(a, "-") || a == "-" {
			rest = append(rest, a)
			i++
			continue
		}
		name := strings.TrimLeft(a, "-")
		key := name
		if eq := strings.IndexByte(name, '='); eq >= 0 {
			key = name[:eq]
		}
		flags = append(flags, a)
		i++
		if !boolFlags[key] && i < len(args) && !strings.HasPrefix(args[i], "-") {
			flags = append(flags, args[i])
			i++
		}
	}
	return append(flags, rest...)
}

func parseColor(s string) (colorparse.Color, error) {
	c, err := colorparse.Parse(s)
	if err != nil {
		return colorparse.Color{}, fmt.Errorf("bad color %q: %v", s, err)
	}
	return c, nil
}

type parseResult struct {
	Hex    string  `json:"hex"`
	RGB    string  `json:"rgb"`
	RGB8   [3]int  `json:"rgb8"`
	Alpha  float64 `json:"alpha"`
	Opaque bool    `json:"opaque"`
}

func runParse(args []string) (err error) {
	fs := flag.NewFlagSet("parse", flag.ContinueOnError)
	jsonOut := fs.Bool("json", false, "emit JSON")
	if err := fs.Parse(reorderFlags(args, map[string]bool{"json": true})); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("parse needs exactly one color argument")
	}
	c, err := parseColor(fs.Arg(0))
	if err != nil {
		return err
	}
	res := parseResult{
		Hex:    colorparse.FormatHex(c),
		RGB:    colorparse.FormatRGB(c),
		RGB8:   [3]int{int(c.R), int(c.G), int(c.B)},
		Alpha:  c.A,
		Opaque: c.Opaque(),
	}
	if *jsonOut {
		return json.NewEncoder(os.Stdout).Encode(res)
	}
	fmt.Printf("hex:    %s\n", res.Hex)
	fmt.Printf("rgb:    %s\n", res.RGB)
	fmt.Printf("rgb8:   %d %d %d\n", res.RGB8[0], res.RGB8[1], res.RGB8[2])
	fmt.Printf("alpha:  %.4f\n", res.Alpha)
	return nil
}

func runConvert(args []string) error {
	fs := flag.NewFlagSet("convert", flag.ContinueOnError)
	to := fs.String("to", "lab", "target space: sRGB|hex|xyz|lab|lch|hsl")
	jsonOut := fs.Bool("json", false, "emit JSON")
	if err := fs.Parse(reorderFlags(args, map[string]bool{"json": true})); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("convert needs exactly one color argument")
	}
	c, err := parseColor(fs.Arg(0))
	if err != nil {
		return err
	}
	switch strings.ToLower(*to) {
	case "srgb", "rgb":
		fmt.Println(colorparse.FormatRGB(c))
	case "hex":
		fmt.Println(colorparse.FormatHex(c))
	case "xyz":
		x, y, z := colorconv.ColorToXYZ(c)
		out := map[string]float64{"x": x, "y": y, "z": z}
		if *jsonOut {
			return json.NewEncoder(os.Stdout).Encode(out)
		}
		fmt.Printf("xyz: x=%.6f y=%.6f z=%.6f\n", x, y, z)
	case "lab":
		l, a, b := colorconv.ColorToLab(c)
		if *jsonOut {
			return json.NewEncoder(os.Stdout).Encode(map[string]float64{"l": l, "a": a, "b": b})
		}
		fmt.Printf("lab: L=%.4f a=%.4f b=%.4f\n", l, a, b)
	case "lch":
		l, a, b := colorconv.ColorToLab(c)
		ll, cc, hh := colorconv.LabToLCh(l, a, b)
		if *jsonOut {
			return json.NewEncoder(os.Stdout).Encode(map[string]float64{"l": ll, "c": cc, "h": hh})
		}
		fmt.Printf("lch: L=%.4f C=%.4f h=%.4f\n", ll, cc, hh)
	case "hsl":
		h, s, l := colorconv.RGBToHSL(c.R, c.G, c.B)
		if *jsonOut {
			return json.NewEncoder(os.Stdout).Encode(map[string]float64{"h": h, "s": s, "l": l})
		}
		fmt.Printf("hsl: h=%.2f s=%.4f l=%.4f\n", h, s, l)
	default:
		return fmt.Errorf("unknown target space %q (want sRGB|hex|xyz|lab|lch|hsl)", *to)
	}
	return nil
}

func runDeltaE(args []string) error {
	fs := flag.NewFlagSet("deltae", flag.ContinueOnError)
	metric := fs.String("m", "2000", "metric: 76|94|2000")
	if err := fs.Parse(reorderFlags(args, nil)); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		return fmt.Errorf("deltae needs two color arguments")
	}
	c1, err := parseColor(fs.Arg(0))
	if err != nil {
		return err
	}
	c2, err := parseColor(fs.Arg(1))
	if err != nil {
		return err
	}
	l1, a1, b1 := colorconv.ColorToLab(c1)
	l2, a2, b2 := colorconv.ColorToLab(c2)
	x := colormetric.Lab{L: l1, A: a1, B: b1}
	y := colormetric.Lab{L: l2, A: a2, B: b2}
	var d float64
	switch *metric {
	case "76":
		d = colormetric.DeltaE76(x, y)
	case "94":
		d = colormetric.DeltaE94(x, y)
	case "2000":
		d = colormetric.DeltaE2000(x, y)
	default:
		return fmt.Errorf("unknown metric %q (want 76|94|2000)", *metric)
	}
	fmt.Printf("deltaE%s = %.6f\n", *metric, d)
	return nil
}

func runContrast(args []string) error {
	fs := flag.NewFlagSet("contrast", flag.ContinueOnError)
	if err := fs.Parse(reorderFlags(args, nil)); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		return fmt.Errorf("contrast needs two color arguments")
	}
	c1, err := parseColor(fs.Arg(0))
	if err != nil {
		return err
	}
	c2, err := parseColor(fs.Arg(1))
	if err != nil {
		return err
	}
	r := colormetric.ContrastRatio(c1, c2)
	lv := colormetric.WCAGLevels(r)
	fmt.Printf("contrast = %.3f:1\n", r)
	fmt.Printf("AA normal: %v  AA large: %v  AAA normal: %v  AAA large: %v\n",
		lv.AANormal, lv.AALarge, lv.AAANormal, lv.AAALarge)
	return nil
}

func runNames(args []string) error {
	fs := flag.NewFlagSet("names", flag.ContinueOnError)
	contains := fs.String("contains", "", "only names containing this substring")
	if err := fs.Parse(reorderFlags(args, nil)); err != nil {
		return err
	}
	for _, n := range colorparse.Names() {
		if *contains == "" || strings.Contains(n, *contains) {
			c, _ := colorparse.ParseNamed(n)
			fmt.Printf("%-22s %s\n", n, colorparse.FormatHex(c))
		}
	}
	return nil
}
