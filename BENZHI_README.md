# beam-vmd — Go 欧拉梁剪力/弯矩/挠度 Web 求解服务与 CLI 工具（Euler-Bernoulli 分段平衡 + 边界积分）
beam-vmd 是一个纯 Go 标准库实现的欧拉梁求解器，通过 HTTP JSON API（POST /api/solve）或命令行输出支座反力、剪力 V(x)、弯矩 M(x) 与挠度 y(x) 沿轴采样，内置交互式 Web 控制台。

## 构建 / 运行 / 测试

```text
go build ./...             # 编译
go run . -solve example/simply.json   # CLI：解算例（JSON 输出）
go run . -solve example/simply.json -table  # CLI：文本表格输出
go run . -http :8080       # Web 求解控制台
go test ./...              # 单元测试（model / beam / sample）
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

容器内可先跑 `go test ./...` 自检，再 `go run . -http :8080` 启动 Web 控制台，或 `go run . -solve example/simply.json` 验证简支跨中集中力算例（跨中 M=PL/4、挠度=PL³/(48EI)）。
