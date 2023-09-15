1. 为 HTTPServer 添加 0-2 秒的随机延时；

   ![](https://s3plus.meituan.net/v1/mss_f32142e8d47149129e9550e929704625/yzz-test-image/18e4d572c9974545b42c11e5f43ac032)

2. 为 HTTPServer 项目添加延时 Metric；

   ![](https://s3plus.meituan.net/v1/mss_f32142e8d47149129e9550e929704625/yzz-test-image/96d317575f794d748f2f0950e3fab116)

   ![](https://s3plus.meituan.net/v1/mss_f32142e8d47149129e9550e929704625/yzz-test-image/055e0fb99fda40b08560fefb727ad6de)

3. 将 HTTPServer 部署至测试集群，并完成 Prometheus 配置；

   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: httpserver
   spec:
     replicas: 1
     selector:
       matchLabels:
         app: httpserver
     template:
       metadata:
         annotations:
           prometheus.io/scrape: "true"
           prometheus.io/port: "8080"
         labels:
           app: httpserver
       spec:
         containers:
           - name: httpserver
             image: isziyang/httpserver:v6.0
             ports:
               - containerPort: 80
   ```

4. 从 Promethus 界面中查询延时指标数据；

   ![](https://s3plus.meituan.net/v1/mss_f32142e8d47149129e9550e929704625/yzz-test-image/549a55d9bacc4e7084f36081ab749db0)

5. （可选）创建一个 Grafana Dashboard 展现延时分配情况。

   ![](https://s3plus.meituan.net/v1/mss_f32142e8d47149129e9550e929704625/yzz-test-image/22f35ea4e2d149bb9a8f1f147bf1c252)