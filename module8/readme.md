# 修改 httpserver

httpserver 代码需要改造为支持以下功能：
1. 配置代码分离，支持从外部读取配置
2. 支持优雅终止

# 编写 K8s spec

# 优雅启动
livenessProbe:
httpGet:
path: /localhost/healthz
port: 8080
initialDelaySeconds: 5
periodSeconds: 10
readinessProbe:
httpGet:
path: /localhost/healthz
port: 8080
initialDelaySeconds: 5
periodSeconds: 10
优雅终止
资源需求和 QoS 保证
探活
日常运维需求，日志等级
配置和代码分离