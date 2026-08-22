# go-waste-routes 标准构建说明

本文件用于评测环境构建 Go 后端镜像，不包含业务题面、修复说明或答案信息。

## 构建

```bash
DOCKER_PLATFORM=linux/amd64 IMAGE_NAME=go-waste-routes ./build_benzhi_docker.sh
```

也支持 `linux/arm64`：

```bash
DOCKER_PLATFORM=linux/arm64 IMAGE_NAME=go-waste-routes ./build_benzhi_docker.sh
```

构建入口为 `benzhi.Dockerfile`，源码通过项目自身的 `go.mod` 和 `go.sum` 安装依赖，并执行 `go build ./...`。
