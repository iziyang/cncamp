# 修改代码并构建镜像

```shell
go mod init github.com/iziyang/cncamp/module12/service0
go mod tidy
go mod init github.com/iziyang/cncamp/module12/service1
go mod tidy
go mod init github.com/iziyang/cncamp/module12/service2
go mod tidy
git add . && git commit -m "go mod"
docker build -t isziyang/service0:v1.0 -f Dokerfile .
docker push isziyang/service0:v1.0
docker build -t isziyang/service1:v1.0 -f Dokerfile .
docker push isziyang/service1:v1.0
docker build -t isziyang/service2:v1.0 -f Dokerfile .
docker push isziyang/service2:v1.0
```
## 安装 jaeger

```shell
kubectl apply -f jaeger.yaml
kubectl edit configmap istio -n istio-system
set tracing.sampling=100
```

## 部署服务

```shell
kubectl create ns module12
kubectl label ns module12 istio-injection=enabled
kubectl -n module12 apply -f service0.yaml
kubectl -n module12 apply -f service1.yaml
kubectl -n module12 apply -f service2.yaml
```

## 生成证书

```shell
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cncamp Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt
kubectl create -n istio-system secret tls cncamp-credential --key=cncamp.io.key --cert=cncamp.io.crt
```

## 修改后的 istio yaml 文件

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: service0
spec:
  gateways:
    - service0
  hosts:
    - 'httpsserver.cncamp.io'
  http:
  - match:
      - port: 443
      - uri:
          exact: /hello/service0
    rewrite:
        uri: "/service0"
    route:
      - destination:
          host: service0
          port:
            number: 80
```

```shell
kubectl apply -f istio-specs.yaml -n module12
```
修改之后，已经支持了安全保证和七层路由以及 traceing 功能。

## 实际访问

```sh
k get svc -nistio-system

istio-ingressgateway   LoadBalancer   $INGRESS_IP
```

### Access the tracing via ingress for 100 times(sampling rate is 1%)

```sh
curl $INGRESS_IP/service0
```

### Check tracing dashboard

```sh
istioctl dashboard jaeger
```