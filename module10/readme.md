# 1. 第一部分

现在你对 Kubernetes 的控制面板的工作机制是否有了深入的了解呢？
是否对如何构建一个优雅的云上应用有了深刻的认识，那么接下来用最近学过的知识把你之前编写的 http 以优雅的方式部署起来吧，你可能需要审视之前代码是否能满足优雅上云的需求。
作业要求：编写 Kubernetes 部署脚本将 httpserver 部署到 Kubernetes 集群，以下是你可以思考的维度。

- 优雅启动
- 优雅终止
- 资源需求和 QoS 保证
- 探活
- 日常运维需求，日志等级
- 配置和代码分离

代码地址见：httpserver_k8s.go
spec 地址见：httpserver.yaml

## 修改 httpserver

httpserver 代码需要改造为支持以下功能：

1. 配置代码分离，支持从外部读取配置
2. 支持优雅终止

## 编写 K8s spec

## 优雅启动

```yaml
readinessProbe:
  httpGet:
    path: /localhost/healthz
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
```

## 优雅终止

捕获系统终止信号，优雅关闭程序。

```go
	// 优雅终止
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		// 优雅关闭
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.WithError(err).Error("Server shutdown failed")
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.WithError(err).Error("Server closed unexpectedly")
		os.Exit(1)
	}

	<-idleConnsClosed
```

## 资源需求和 QoS 保证

```yaml
resources:
  limits:
    cpu: "1"
    memory: "512Mi"
  requests:
    cpu: "500m"
    memory: "256Mi"
```

## 探活

```yaml
livenessProbe:
  httpGet:
    path: /localhost/healthz
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
```

## 日常运维需求，日志等级

通过变量来控制日志级别，同时日志要输出到文件中去。

```go
  // 设置日志输出到文件
	logger.SetOutput(os.Stdout)
  logLevel := os.Getenv("LOG_LEVEL")
	logger.Debug("loglevel is:", logLevel)
	if logLevel == "DEBUG" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
```

## 配置和代码分离

通过 configMap 加载配置：

```yaml
env:
  - name: VERSION
    value: "1.0"
  - name: LOG_LEVEL
    valueFrom:
      configMapKeyRef:
        name: log-config
        key: LOG_LEVEL
```

# 2. 第二部分

除了将 httpServer 应用优雅的运行在 Kubernetes 之上，我们还应该考虑如何将服务发布给对内和对外的调用方。
来尝试用 Service, Ingress 将你的服务发布给集群外部的调用方吧。
在第一部分的基础上提供更加完备的部署 spec，包括（不限于）：

- Service
- Ingress

可以考虑的细节

- 如何确保整个应用的高可用。
- 如何通过证书保证 httpServer 的通讯安全。

## Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: httpserver-k8s-service
spec:
  selector:
    app: httpserver-k8s
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
```

## Ingress

### 生成证书相关信息

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=cncamp.com/O=cncamp" -addext "subjectAltName = DNS:cncamp.com"
```

### 创建 secret

```bash
kubectl create secret tls httpserver-tls --cert=./tls.crt --key=./tls.key
```

### Spec

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: httpserver-ingress
spec:
  tls:
    - hosts:
        - httpserver.k8s.com
      secretName: httpserver-tls
  rules:
    - host: httpserver.k8s.com
      http:
        paths:
          - pathType: Prefix
            path: "/"
            backend:
              service:
                name: httpserver-k8s-service
                port:
                  number: 80
```

## 验证

ingress 的 service IP 是：

```bash
[howardyuan@node1 ~/my_code/cncamp/module8]$ k get svc -n ingress-nginx
NAME                                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             NodePort    10.109.53.159   <none>        80:31296/TCP,443:32734/TCP   13h
ingress-nginx-controller-admission   ClusterIP   10.101.3.197    <none>        443/TCP                      13h
```

访问测试：

```bash
curl -H "Host: httpserver.k8s.com" https://10.109.53.159 -v -k
```

![](https://s3plus.meituan.net/v1/mss_f32142e8d47149129e9550e929704625/yzz-test-image/e72c74bc7f884094b82a9e0e0a8ac0d3)

