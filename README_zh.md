# API Proxy 项目

## 简介

API Proxy 是一个基于 Gin 框架的 Go 语言项目，用于代理和转发 API 请求。它支持多平台的 API 请求代理，并提供 IP 白名单和授权验证功能。

## 功能

- 支持多平台 API 请求代理
- IP 白名单验证
- 授权验证
- 支持 CORS
- 支持流式响应

## 安装

1. 克隆项目到本地：

   ```bash
   git clone https://github.com/heavi715/api-proxy.git
   ```

2. 进入项目目录：

   ```bash
   cd api-proxy
   ```

3. 安装依赖：

   ```bash
   go mod tidy
   ```

4. 配置 `config.json` 文件：

   根据需要修改 `config/config.json` 文件中的配置项，包括 `server_addr`、`source_list`、`server_key_list`、`proxy_url` 和 `platform_list`。

## 使用

1. 启动服务：

   ```bash
   go run main.go
   ```

2. 访问健康检查接口：

   在浏览器中访问 `http://<server_addr>/health`，应返回 `success`。

3. 代理请求：

   通过 `/proxy/:platform/:source/*url` 路径发送请求，其中 `:platform` 是目标平台名称，`:source` 是来源标识，`*url` 是目标 API 的路径。

## 配置说明

- `server_addr`: 服务器监听地址和端口。
- `source_list`: 允许的来源标识列表。如果为空不校验
- `server_key_list`: 允许的服务器密钥列表。
- `proxy_url`: 本地 HTTP 代理配置。如果为空不走代理
- `platform_list`: 支持的目标平台配置，包括 `name`、`url`、`header_key` 和 `header_values`。

## 贡献

欢迎提交问题和请求合并。请确保在提交请求前运行所有测试并遵循代码风格指南。

## 许可证

该项目使用 MIT 许可证。详情请参阅 LICENSE 文件。
