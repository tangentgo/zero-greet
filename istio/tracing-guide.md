# Istio 链路追踪配置指南

本文档说明如何在 book 服务中启用和使用 Istio 分布式链路追踪。

## 📋 前置条件

- ✅ Kubernetes 集群已部署
- ✅ Istio 已安装
- ✅ Jaeger 已部署在 `istio-system` namespace
- ✅ 命名空间已启用 Istio 注入 (`istio-injection: enabled`)

## 🚀 部署步骤

### 1. 应用 Istio Telemetry 配置

启用全局和命名空间级别的追踪配置：

```powershell
# 应用 Telemetry 配置
kubectl apply -f .\istio\telemetry.yaml

# 验证配置
kubectl get telemetry -n istio-system
kubectl get telemetry -n api-book
```

**说明：**
- 全局配置 (`mesh-default`): 为整个网格设置 100% 采样率
- 命名空间配置 (`book-tracing`): 为 api-book 添加自定义标签

### 2. 重新构建和部署服务

```powershell
# 进入 book 目录
cd book

# 构建 Docker 镜像
docker build -t fengyuxiu/zero-greet-book:latest .

# 推送镜像
docker push fengyuxiu/zero-greet-book:latest

# 部署到 Kubernetes
kubectl apply -f deploy.yaml

# 查看部署状态
kubectl get pods -n api-book
kubectl describe pod <pod-name> -n api-book
```

### 3. 验证 Istio Sidecar 注入

```powershell
# 检查 Pod 是否有 2 个容器（应用 + Envoy）
kubectl get pods -n api-book

# 应该看到 READY 显示 2/2
# NAME                        READY   STATUS    RESTARTS   AGE
# api-book-xxxxxxxxxx-xxxxx   2/2     Running   0          1m

# 查看 Envoy 代理日志
kubectl logs -n api-book <pod-name> -c istio-proxy --tail=50
```

## 🧪 测试链路追踪

### 1. 发送测试请求

```powershell
# 方式1：通过域名访问（如果已配置）
curl https://book.concises.net/bookinfo -d '{"title":"Go语言实战"}'

# 方式2：通过 port-forward
kubectl port-forward -n api-book svc/book-service 8888:8888
curl http://localhost:8888/bookinfo -X POST -H "Content-Type: application/json" -d '{"title":"Go语言实战"}'

# 发送多个请求以生成追踪数据
for ($i=1; $i -le 10; $i++) {
    curl http://localhost:8888/bookinfo -X POST -H "Content-Type: application/json" -d "{`"title`":`"Book $i`"}"
    Start-Sleep -Milliseconds 500
}
```

### 2. 访问 Jaeger UI

```powershell
# 方式1：通过域名访问（如果已配置）
# 浏览器打开: https://jaeger.concises.net

# 方式2：通过 port-forward
kubectl port-forward -n istio-system svc/tracing 16686:80

# 浏览器打开: http://localhost:16686
```

### 3. 在 Jaeger UI 中查看追踪

1. **选择服务**: 在左侧 "Service" 下拉框选择 `book-service.api-book`
2. **点击 "Find Traces"**: 查看最近的追踪记录
3. **点击某个追踪**: 查看详细的调用链路

**预期看到的追踪信息：**
```
istio-ingressgateway.istio-system
  └─ book-service.api-book (GET /bookinfo)
      └─ greet-service.api-greet (GET /you)
```

## 📊 追踪信息说明

### Span 标签

每个 span 会包含以下信息：

- **Service Name**: 服务名称 (book-api, greet-api)
- **Operation**: HTTP 方法和路径 (GET /bookinfo)
- **Duration**: 请求耗时
- **Tags**: 
  - `http.method`: HTTP 方法
  - `http.url`: 请求 URL
  - `http.status_code`: 响应状态码
  - `component`: istio-proxy
  - `node_id`: Pod 标识

### 追踪头

Istio 使用以下 HTTP 头进行追踪：

| 头名称 | 说明 | 示例 |
|--------|------|------|
| `x-request-id` | 请求唯一标识 | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |
| `x-b3-traceid` | 追踪 ID | `80f198ee56343ba864fe8b2a57d3eff7` |
| `x-b3-spanid` | Span ID | `05e3ac9a4f6e3b90` |
| `x-b3-parentspanid` | 父 Span ID | `e457b5a2e4d86bd1` |
| `x-b3-sampled` | 是否采样 | `1` |
| `b3` | 紧凑格式 | `80f198ee56343ba864fe8b2a57d3eff7-05e3ac9a4f6e3b90-1` |

## 🔍 调试技巧

### 1. 检查追踪配置

```powershell
# 查看 Telemetry 配置
kubectl get telemetry -n istio-system -o yaml
kubectl get telemetry -n api-book -o yaml

# 查看 Envoy 配置中的追踪设置
kubectl exec -n api-book <pod-name> -c istio-proxy -- curl localhost:15000/config_dump | grep -A 20 tracing
```

### 2. 查看 Envoy 追踪日志

```powershell
# 启用 Envoy 调试日志
kubectl exec -n api-book <pod-name> -c istio-proxy -- curl -X POST localhost:15000/logging?level=debug

# 查看日志
kubectl logs -n api-book <pod-name> -c istio-proxy -f | grep -i trace
```

### 3. 检查 Jaeger 连接

```powershell
# 检查 Jaeger Collector 服务
kubectl get svc -n istio-system jaeger-collector

# 测试连接
kubectl run -it --rm debug --image=curlimages/curl --restart=Never -- \
  curl -v http://jaeger-collector.istio-system:14268/api/traces
```

### 4. 验证追踪头传播

```powershell
# 查看应用日志
kubectl logs -n api-book <pod-name> -c api-book

# 添加调试代码打印追踪头
# 在 bookinfologic.go 中添加：
# logx.Infof("Trace ID: %v", l.ctx.Value("x-b3-traceid"))
```

## 🎯 常见问题

### Q1: 在 Jaeger 中看不到追踪数据

**解决方案：**

1. 检查采样率是否 > 0
```powershell
kubectl get telemetry -n istio-system -o yaml | grep samplingPercentage
```

2. 检查 Jaeger 服务是否运行
```powershell
kubectl get pods -n istio-system -l app=jaeger
```

3. 检查 Envoy 是否正确注入
```powershell
kubectl get pods -n api-book -o jsonpath='{.items[*].spec.containers[*].name}'
# 应该看到: api-book istio-proxy
```

### Q2: 跨服务调用没有关联

**解决方案：**

确保所有服务都正确传播追踪头：
- book → greet 调用时需要传递所有 B3 头
- 检查代码中是否调用了 `propagateTracingHeaders`

### Q3: 追踪数据不完整

**解决方案：**

1. 增加采样率到 100%（开发环境）
```yaml
spec:
  tracing:
  - randomSamplingPercentage: 100.0
```

2. 检查网络策略是否阻止了 Jaeger 流量

## 📈 生产环境建议

### 1. 调整采样率

生产环境建议使用较低的采样率以减少性能影响：

```yaml
spec:
  tracing:
  - randomSamplingPercentage: 1.0  # 1% 采样率
```

### 2. 配置 Jaeger 持久化

默认 Jaeger 使用内存存储，建议配置持久化后端：
- Elasticsearch
- Cassandra
- Badger（文件系统）

### 3. 启用追踪标签

添加更多自定义标签以便分析：

```yaml
spec:
  tracing:
  - customTags:
      environment:
        literal:
          value: "production"
      version:
        literal:
          value: "v1.0.0"
      region:
        literal:
          value: "us-west-2"
```

### 4. 设置追踪保留策略

在 Jaeger 中配置数据保留时间，避免存储过多历史数据。

## 🔗 相关资源

- [Istio Distributed Tracing](https://istio.io/latest/docs/tasks/observability/distributed-tracing/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [B3 Propagation](https://github.com/openzipkin/b3-propagation)

## 📝 代码改动总结

### 修改的文件

1. **istio/telemetry.yaml** (新建)
   - 配置全局和命名空间级别的追踪

2. **book/etc/book-api.yaml**
   - 添加 Telemetry 配置（可选）

3. **book/internal/handler/bookinfohandler.go**
   - 添加 `propagateTracingHeaders` 函数从请求中提取追踪头

4. **book/internal/logic/bookinfologic.go**
   - 添加 `propagateTracingHeaders` 方法将追踪头传播到下游服务

5. **book/deploy.yaml**
   - 添加 `version: v1` 标签
   - 添加环境变量（POD_NAME, POD_NAMESPACE, POD_IP）
   - 添加 Istio 注解

### 关键代码片段

**传播追踪头（Handler 层）：**
```go
func propagateTracingHeaders(r *http.Request) context.Context {
    ctx := r.Context()
    tracingHeaders := []string{
        "x-request-id", "x-b3-traceid", "x-b3-spanid",
        "x-b3-parentspanid", "x-b3-sampled", "x-b3-flags", "b3",
    }
    for _, header := range tracingHeaders {
        if value := r.Header.Get(header); value != "" {
            ctx = context.WithValue(ctx, header, value)
        }
    }
    return ctx
}
```

**传播追踪头（Logic 层）：**
```go
func (l *BookInfoLogic) propagateTracingHeaders(req *http.Request) {
    tracingHeaders := []string{...}
    for _, header := range tracingHeaders {
        if value := l.ctx.Value(header); value != nil {
            if strValue, ok := value.(string); ok && strValue != "" {
                req.Header.Set(header, strValue)
            }
        }
    }
}
```

---

**文档版本**: v1.0  
**最后更新**: 2025-11-09  
**适用环境**: Istio + Jaeger + go-zero
