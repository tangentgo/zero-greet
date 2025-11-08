# Jaeger UI 中查看 greet 服务追踪的完整指南

## 🎯 问题：在 Jaeger 中看不到 greet 服务

### ✅ 验证服务已正确配置

```powershell
# 1. 确认两个服务都在运行且有 Istio sidecar
kubectl get pods -n api-book
kubectl get pods -n api-greet
# 应该看到 READY 2/2

# 2. 确认 Telemetry 配置已应用
kubectl get telemetry -A
# 应该看到 mesh-default, book-tracing, greet-tracing

# 3. 确认服务收到请求
kubectl logs -n api-greet -l app=api-greet -c api-greet --tail=20 | Select-String "GET /you"
```

## 📊 在 Jaeger UI 中正确查看 greet 服务

### 方法 1：通过完整调用链查看（推荐）

1. **访问 Jaeger UI**
   ```powershell
   # 浏览器打开
   Start-Process "https://jaeger.concises.net"
   # 或
   kubectl port-forward -n istio-system svc/tracing 16686:80
   Start-Process "http://localhost:16686"
   ```

2. **选择入口服务**
   - Service 下拉框选择：`istio-ingressgateway.istio-system`
   - 或选择：`book-service.api-book`

3. **点击 "Find Traces" 按钮**

4. **点击任意一个 trace 查看详情**
   - 你会看到一个时间线视图
   - 展开 trace，应该看到多个 span：

   ```
   📊 Trace Timeline (点击展开查看):
   
   ┌─ istio-ingressgateway.istio-system [~50ms]
   │
   ├─ book-service.api-book [~45ms]
   │  └─ GET /info
   │     📍 Tags:
   │        - http.method: GET
   │        - http.status_code: 200
   │        - component: istio-proxy
   │
   │     └─ greet-service.api-greet [~20ms] 👈 这里就是 greet 服务！
   │        └─ GET /you
   │           📍 Tags:
   │              - http.method: GET
   │              - http.status_code: 200
   │              - component: istio-proxy
   └─
   ```

5. **点击 greet span 查看详细信息**
   - Duration: 服务响应时间
   - Tags: HTTP 方法、状态码、URL 等
   - Process: 服务的元数据（namespace, pod name 等）

### 方法 2：直接查询 greet 服务

1. **Service 下拉框尝试以下名称**：
   - `greet-service.api-greet`
   - `api-greet.api-greet`
   - `greet-api` （如果显示的话）

2. **如果服务名不在列表中**，说明：
   - greet 只作为下游服务被调用
   - 没有直接的外部流量到 greet
   - **这是正常的！** greet 的追踪数据在 book 的 trace 中

### 方法 3：使用高级搜索

1. 在 Jaeger UI 右上角点击 "🔍 Search"

2. 使用以下查询条件：
   ```
   service="greet-service.api-greet"
   ```

3. 或者搜索特定的 tag：
   ```
   http.url=/you
   ```

## 🧪 生成测试追踪数据

```powershell
# 生成完整的调用链（book → greet）
for ($i=1; $i -le 20; $i++) {
    curl -s "https://book.concises.net/info?title=TestBook$i" | Out-Null
    Write-Host "Request $i completed"
    Start-Sleep -Milliseconds 200
}

# 等待 5-10 秒让数据被 Jaeger 处理
Start-Sleep -Seconds 5

# 然后在 Jaeger UI 中刷新并查找
```

## 🔍 故障排查

### 问题 1: Service 列表中没有任何服务

**可能原因**：
- Jaeger 没有收到任何追踪数据
- 采样率设置为 0

**解决方案**：
```powershell
# 检查 Telemetry 配置
kubectl get telemetry -n istio-system mesh-default -o yaml

# 确认 randomSamplingPercentage 是 100.0
# 如果不是，更新配置：
kubectl apply -f istio/telemetry.yaml
```

### 问题 2: 只看到 book 服务，看不到 greet

**这是正常的！**

**原因**：
- greet 是 book 的下游服务
- 它不接收直接的外部流量（通过 Istio Gateway）
- 它的追踪数据嵌套在 book 的 trace 中

**验证方法**：
1. 在 Service 列表选择 `book-service.api-book`
2. 找到一个 trace
3. **点击并展开** trace
4. 你会看到 greet 作为 child span

### 问题 3: 追踪数据显示不完整

**可能原因**：
- 追踪头没有正确传播
- Handler 代码没有更新

**验证追踪头传播**：
```powershell
# 查看 book 服务日志，确认有 trace ID
kubectl logs -n api-book -l app=api-book -c api-book --tail=5

# 查看 greet 服务日志，确认有相同的 trace ID
kubectl logs -n api-greet -l app=api-greet -c api-greet --tail=5

# 两个服务的日志应该显示相同的 trace ID
```

### 问题 4: Jaeger UI 显示空白

**解决方案**：
```powershell
# 1. 确认 Jaeger 服务正在运行
kubectl get pods -n istio-system -l app=jaeger

# 2. 检查 Jaeger 日志
kubectl logs -n istio-system -l app=jaeger --tail=100

# 3. 确认可以访问 Jaeger UI
kubectl port-forward -n istio-system svc/tracing 16686:80
# 浏览器访问 http://localhost:16686
```

## 📝 关键理解

### greet 服务的追踪特点

1. **greet 不是入口服务**
   - 它不直接接收来自 Istio Gateway 的流量
   - 它被 book 服务调用（服务间通信）

2. **追踪数据的层级结构**
   ```
   istio-ingressgateway (根 span)
   └─ book-service (子 span)
       └─ greet-service (孙 span) 👈 嵌套在这里
   ```

3. **正常的显示方式**
   - 在 Jaeger UI 中，greet 不会作为独立的 trace 出现
   - 它作为 book trace 的一部分（child span）
   - 这是**分布式追踪的正常行为**

## ✅ 验证追踪正常工作的标志

1. **在 Service 列表中看到**：
   - `istio-ingressgateway.istio-system`
   - `book-service.api-book`

2. **点击 book 的 trace 后**：
   - 看到多个 span（至少 2 个）
   - 其中一个 span 的 service name 是 `greet-service.api-greet`
   - Span 之间有父子关系（缩进显示）

3. **Span 详情包含**：
   - Operation name: `GET /you`
   - Tags: `http.status_code=200`
   - Duration: 合理的响应时间（几毫秒到几十毫秒）

## 🎬 截图位置说明

在 Jaeger UI 中查看的正确位置：

```
1. Jaeger UI 首页
   ┌────────────────────────────────────┐
   │ [Service ▼] book-service.api-book  │
   │ [Find Traces]                      │
   └────────────────────────────────────┘
   
2. Traces 列表页面
   ┌────────────────────────────────────┐
   │ Trace ID: abc123...                │ ← 点击这里
   │ 2 Spans | 50ms                     │
   └────────────────────────────────────┘
   
3. Trace 详情页面（展开后）
   ┌────────────────────────────────────┐
   │ ▼ book-service.api-book [45ms]     │
   │   ├─ GET /info                     │
   │   │                                │
   │   └─ ▼ greet-service.api-greet     │ ← greet 在这里！
   │       └─ GET /you [20ms]           │
   └────────────────────────────────────┘
```

## 🚀 推荐的查看步骤

1. **生成追踪数据**（运行上面的测试脚本）
2. **打开 Jaeger UI**
3. **选择 `book-service.api-book`**
4. **点击 "Find Traces"**
5. **选择最近的一个 trace**
6. **点击展开 trace 的各个 span**
7. **找到 `greet-service.api-greet` 的 span**

这就是你应该看到 greet 服务追踪的地方！

---

**最后提示**：如果按照以上步骤仍然看不到 greet，可能需要：
1. 等待几分钟让 Jaeger 索引数据
2. 刷新 Jaeger UI 页面
3. 确认时间范围选择正确（默认是最近 1 小时）
