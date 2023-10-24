# 修改代码并构建镜像
修改后的代码，已经支持了 tracing

```go
req, err := http.NewRequest("GET", "http://service1", nil)
req, err := http.NewRequest("GET", "http://service2", nil)
```

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

### 将访问到 /hello/service0 转到 / 路径下
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
        uri: "/"
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

```sh
howardyuan@node1:module12$ curl --resolve httpsserver.cncamp.io:443:10.103.249.202 https://httpsserver.cncamp.io/hello/service0 -v -k
* Added httpsserver.cncamp.io:443:10.103.249.202 to DNS cache
* About to connect() to httpsserver.cncamp.io port 443 (#0)
*   Trying 10.103.249.202...
* Connected to httpsserver.cncamp.io (10.103.249.202) port 443 (#0)
* Initializing NSS with certpath: sql:/etc/pki/nssdb
* skipping SSL peer certificate verification
* SSL connection using TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
* Server certificate:
*       subject: CN=*.cncamp.io,O=cncamp Inc.
*       start date: Oct 24 00:19:06 2023 GMT
*       expire date: Oct 23 00:19:06 2024 GMT
*       common name: *.cncamp.io
*       issuer: CN=*.cncamp.io,O=cncamp Inc.
> GET /hello/service0 HTTP/1.1
> User-Agent: curl/7.29.0
> Host: httpsserver.cncamp.io
> Accept: */*
> 
< HTTP/1.1 200 OK
< date: Tue, 24 Oct 2023 23:39:26 GMT
< content-type: text/plain; charset=utf-8
< x-envoy-upstream-service-time: 98
< server: istio-envoy
< transfer-encoding: chunked
< 
===================Details of the http request header:============
HTTP/1.1 200 OK
Content-Length: 915
Content-Type: text/plain; charset=utf-8
Date: Tue, 24 Oct 2023 23:39:26 GMT
Server: envoy
X-Envoy-Upstream-Service-Time: 45

===================Details of the http request header:============
HTTP/1.1 200 OK
Content-Length: 13
Accept: */*
Accept-Encoding: gzip,gzip
Content-Type: text/plain; charset=utf-8
Date: Tue, 24 Oct 2023 23:39:26 GMT
Server: envoy
User-Agent: Go-http-client/1.1,Go-http-client/1.1,curl/7.29.0
Version: 
X-B3-Parentspanid: 1ee43f7bf9ed5dcf
X-B3-Sampled: 1
X-B3-Spanid: d0c0b9d482cffb4e
X-B3-Traceid: 70c23cfc738b52bc8cb31591affd9f5e
X-Envoy-Attempt-Count: 1
X-Envoy-External-Address: 9.135.14.174
X-Envoy-Original-Path: /hello/service0
X-Envoy-Upstream-Service-Time: 1
X-Forwarded-Client-Cert: By=spiffe://cluster.local/ns/module12/sa/default;Hash=ca42ff7b95910e94cc0313d485e6689821b4e14e80e10aa36741fc3933e328eb;Subject="";URI=spiffe://cluster.local/ns/module12/sa/default
X-Forwarded-For: 9.135.14.174
X-Forwarded-Proto: https
X-Request-Id: 6f8ed0c3-5001-9bfe-a37f-beaebf955d8e

* Connection #0 to host httpsserver.cncamp.io left intact
Hello, World!
```
