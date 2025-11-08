# 快速部署 Istio 链路追踪

# 1. 应用 Telemetry 配置
kubectl apply -f istio/telemetry.yaml

# 2. 重新部署 book 服务
kubectl apply -f book/deploy.yaml

# 3. 重启 Pod 以应用新配置
kubectl rollout restart deployment/api-book -n api-book

# 4. 等待 Pod 就绪
kubectl wait --for=condition=ready pod -l app=api-book -n api-book --timeout=60s

# 5. 查看 Pod 状态
kubectl get pods -n api-book

# 6. 发送测试请求
kubectl port-forward -n api-book svc/book-service 8888:8888

# 在另一个终端发送请求
curl http://localhost:8888/bookinfo -X POST -H "Content-Type: application/json" -d '{"title":"Go语言实战"}'

# 7. 访问 Jaeger UI
# 如果已配置域名：https://jaeger.concises.net
# 或使用 port-forward：
kubectl port-forward -n istio-system svc/tracing 16686:80
# 浏览器打开: http://localhost:16686
