# colorsci

colorsci 是一个色彩科学工具库与命令行工具：解析 CSS 颜色语法（hex、rgb()/hsl() 函数式、W3C 命名颜色），
在 sRGB / 线性 RGB / CIE XYZ(D65) / CIELAB / CIELCh 色彩空间之间转换，支持越界颜色回钳到 sRGB 色域，
并计算 CIE76 / CIE94 / CIEDE2000 色差与 WCAG 2.1 对比度。纯标准库、离线可构建，无外部依赖。

## 构建 / 运行 / 测试

```text
go build ./...        # 编译（含 example/）
go test ./...         # 单元测试（3 个 internal 包，70 个 TestXxx）
go run . parse "hsl(120 100% 50%)"     # CLI 示例
go run . contrast "#ffffff" "#767676"  # 对比度示例
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```

容器内验证：`cd /src && go build ./... && go test ./...`
