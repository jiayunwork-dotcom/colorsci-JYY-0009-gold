# colorsci

colorsci is a color science toolkit for Go. It parses CSS color syntax
(hex, `rgb()`/`hsl()` functional forms, and the W3C named color table),
converts between the sRGB, linear RGB, CIE XYZ (D65), CIELAB and CIELCh
color spaces, clamps out-of-gamut colors back into the sRGB gamut, and
computes the CIE76 / CIE94 / CIEDE2000 color difference metrics plus
WCAG 2.1 contrast ratios. It is a pure standard-library module with no
external dependencies, suitable for embedding in rendering, image
analysis, data-viz and accessibility tooling.

## Packages

- `internal/colorparse` — CSS color parsing and formatting
- `internal/colorconv` — color space conversions and gamut clamping
- `internal/colormetric` — color difference and contrast metrics

## CLI

```
colorsci parse <color>                 parse and re-format a CSS color
colorsci convert <color> -to <space>   convert to sRGB/hex/xyz/lab/lch/hsl
colorsci deltae <c1> <c2> -m <76|94|2000>
colorsci contrast <c1> <c2>            WCAG contrast ratio and levels
colorsci names [-contains s]           list supported named colors
```

Example:

```
$ go run . convert "#ff0000" -to lab
lab: L=53.2329 a=80.1094 b=67.2201

$ go run . contrast "#ffffff" "#767676"
contrast = 4.539:1
AA normal: true  AA large: true  AAA normal: false  AAA large: true
```

## Library usage

```go
c, err := colorparse.Parse("hsl(120 100% 50%)")
if err != nil { /* ... */ }

l, a, b := colorconv.ColorToLab(c)
d := colormetric.DeltaE2000(
    colormetric.Lab{L: l, A: a, B: b},
    colormetric.Lab{L: 50, A: 0, B: 0},
)

ratio := colormetric.ContrastRatio(c, colorparse.MustParse("#ffffff"))
```

`MustParse` is provided in the `colorparse` package for inputs that are
known-valid at compile time; prefer `Parse` elsewhere.

## Build and test

```
go build ./...
go test ./...
go vet ./...
```

The module targets Go 1.21 and builds offline.
