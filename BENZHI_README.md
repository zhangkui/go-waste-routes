# go-waste-routes 标准构建说明

本文件用于评测环境构建 Go 后端和 Vue 前端镜像，不包含业务题面、修复说明或答案信息。

## 构建

```bash
./build_benzhi_docker.sh go-waste-routes linux/amd64
```

也支持 `linux/arm64`：

```bash
./build_benzhi_docker.sh go-waste-routes linux/arm64
```

构建入口为 `benzhi.Dockerfile`，镜像安装 Go 依赖和 `frontend/` 的 Node.js 依赖，并执行后端编译和前端构建。
