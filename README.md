# IPv6 DDNS 更新工具

这是一个自动更新 DNSPod DNS 记录的工具，主要用于动态 IPv6 地址的 DDNS 更新。

## 功能特性

- 自动检测本地 IPv6 地址
- 自动更新 DNSPod DNS 记录
- 错误重试机制
- 邮件通知功能
- 健康检查
- 反向代理

## 配置说明

配置文件 `config.yaml` 示例：

```yaml
tencent:
  secret_id: "your_secret_id"
  secret_key: "your_secret_key"
domain:
  domain: "example.com"
  sub_domain: "www"
check_interval: 300  # 检查间隔（秒）
email:
  smtp_server: "smtp.example.com"
  smtp_port: 587
  username: "your_email@example.com"
  password: "your_password"
  to_email: "notify@example.com"
```

## 使用方法

1. 准备配置文件
2. 运行程序：`go run main.go`

## 错误处理

- 当连续3次更新失败时，将发送邮件通知
- 使用指数退避算法进行重试

## 开发说明

项目使用 Go 模块管理依赖，主要依赖包括：

- github.com/cenkalti/backoff/v4：用于实现重试机制
- github.com/tencentcloud/tencentcloud-sdk-go：腾讯云 API SDK
- gopkg.in/yaml.v2：用于解析 YAML 配置文件

## 贡献指南

欢迎提交 Issue 和 Pull Request。在提交 PR 之前，请确保：

1. 代码已经格式化（go fmt）
2. 所有测试通过（go test）
3. 已添加必要的注释和文档
