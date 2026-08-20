# beam-vmd

beam-vmd 是欧拉梁（Euler–Bernoulli）剪力、弯矩与挠度核算工具：给定跨度 L、抗弯刚度 EI、约束类型与梁上荷载，沿轴线分段平衡给出剪力 V(x) 与弯矩 M(x)，由边界条件解支座反力，再按 EI y''=M 两次积分给出转角与挠度，输出沿 x 的采样、支座反力、极值与交叉校验报告。纯标准库，无第三方依赖，可离线构建。

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
