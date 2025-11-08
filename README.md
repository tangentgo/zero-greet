# Zero-Greet 微服务项目

基于 go-zero 框架构建的微服务示例项目，展示了 RESTful API 和 gRPC 服务的集成，并在 Kubernetes + Istio 环境中运行。

## 项目架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         外部访问                                  │
│  https://book.concises.net   https://greet.concises.net         │
└────────────────┬────────────────────────┬───────────────────────┘
                 │                        │
                 ▼                        ▼
        ┌─────────────────┐      ┌─────────────────┐
        │   Book API      │      │   Greet API     │
        │  (REST 8888)    │─────▶│  (REST 8888)    │
        │  api-book ns    │ HTTP │  api-greet ns   │
        └─────────────────┘      └────────┬────────┘
                                           │ gRPC
                                           ▼
                                  ┌─────────────────┐
                                  │   User RPC      │
                                  │  (gRPC 8080)    │
                                  │  grpc-user ns   │
                                  └─────────────────┘
```

### 服务说明

#### 1. Book API (REST)
- **命名空间**: `api-book`
- **端口**: 8888
- **域名**: `https://book.concises.net`
- **功能**: 
  - 提供图书信息查询 API
  - 调用 Greet API 获取问候语
  - 使用 faker 生成随机作者名
- **主要接口**:
  - `GET /info?title=<书名>` - 查询书籍信息

#### 2. Greet API (REST)
- **命名空间**: `api-greet`
- **端口**: 8888
- **域名**: `https://greet.concises.net`
- **功能**:
  - 提供问候语 API
  - 通过 gRPC 调用 User 服务获取用户信息
- **主要接口**:
  - `GET /:name` - 获取问候语（name 可选 you|me）
  - `GET /` - 基础健康检查

#### 3. User RPC (gRPC)
- **命名空间**: `grpc-user`
- **端口**: 8080
- **服务名**: `grpc-user-service`
- **功能**: 
  - 提供用户相关 gRPC 服务
  - 响应 Ping 请求
- **主要方法**:
  - `user.User/Ping` - Ping/Pong 测试

## 技术栈

- **框架**: go-zero (API + RPC)
- **容器编排**: Kubernetes
- **服务网格**: Istio
- **网关**: Istio Gateway (支持 HTTPS)
- **协议**: REST (HTTP/HTTPS) + gRPC
- **开发语言**: Go

## 项目结构

```
zero-greet/
├── book/                    # Book API 服务
│   ├── book.api            # API 定义
│   ├── book.go             # 主程序
│   ├── Dockerfile          # 容器镜像
│   ├── deploy.yaml         # K8s 部署配置
│   ├── Taskfile.yaml       # 构建任务
│   ├── etc/
│   │   └── book-api.yaml   # 服务配置（含 GreetAPI 配置）
│   └── internal/
│       ├── config/         # 配置结构
│       ├── handler/        # HTTP 处理器
│       ├── logic/          # 业务逻辑
│       ├── svc/            # 服务上下文（HTTP 客户端）
│       └── types/          # 类型定义
│
├── greet/                  # Greet API 服务
│   ├── greet.api           # API 定义
│   ├── greet.go            # 主程序
│   ├── Dockerfile
│   ├── deploy.yaml
│   ├── Taskfile.yaml
│   ├── etc/
│   │   └── greet-api.yaml  # 服务配置（含 RPCuser 配置）
│   ├── internal/
│   │   ├── config/         # 配置结构（含 RPC 客户端配置）
│   │   ├── handler/
│   │   ├── logic/          # 业务逻辑（调用 User RPC）
│   │   ├── svc/            # 服务上下文（gRPC 客户端）
│   │   └── types/
│   └── proto/user/         # User gRPC 客户端代码
│
├── user/                   # User RPC 服务
│   ├── user.proto          # gRPC 服务定义
│   ├── user.go             # 主程序
│   ├── Dockerfile
│   ├── deploy.yaml         # K8s 部署配置
│   ├── fleet.yaml
│   ├── Taskfile.yaml
│   ├── etc/
│   │   └── user.yaml       # 服务配置
│   ├── internal/
│   │   ├── config/
│   │   ├── logic/          # gRPC 业务逻辑
│   │   ├── server/         # gRPC 服务器
│   │   └── svc/            # 服务上下文
│   ├── user/               # 生成的 gRPC 代码
│   └── userclient/         # gRPC 客户端封装
│
├── istio/
│   └── enable-mtls.yaml    # Istio mTLS 配置
│
├── .vscode/
│   └── launch.json         # VS Code 调试配置
│
└── test-book-api.ps1       # API 测试脚本
```

## 服务间调用关系

### Book → Greet (HTTP)
```yaml
# book/etc/book-api.yaml
GreetAPI:
  Endpoint: http://greet-service.api-greet:8888
  Timeout: 5000
```

```go
// book/internal/logic/bookinfologic.go
greetURL := fmt.Sprintf("%s/you", l.svcCtx.Config.GreetAPI.Endpoint)
httpResp, err := l.svcCtx.GreetClient.Do(httpReq)
```

### Greet → User (gRPC)
```yaml
# greet/etc/greet-api.yaml
RPCuser:
  Timeout: 300000
  Endpoints:
    - grpc-user-service.grpc-user:8080
```

```go
// greet/internal/logic/greetlogic.go
rpcResp, err := l.svcCtx.RPCuser.Ping(l.ctx, &user.Request{
    Ping: req.Name,
})
```

## Kubernetes 部署

### 命名空间配置

所有命名空间都启用了 Istio sidecar 自动注入：

```yaml
metadata:
  labels:
    istio-injection: enabled
```

### 服务配置

| 服务 | 命名空间 | Service 名称 | ClusterIP | 端口 |
|------|----------|--------------|-----------|------|
| Book API | api-book | book-service | 34.118.x.x | 8888 |
| Greet API | api-greet | greet-service | 34.118.238.93 | 8888 |
| User RPC | grpc-user | grpc-user-service | 34.118.231.86 | 8080 |

### Gateway 配置

通过 Istio Gateway 和 HTTPRoute 暴露服务：

- `book.concises.net` → book-service:8888
- `greet.concises.net` → greet-service:8888
- 自动 HTTP → HTTPS 重定向

## Istio 服务网格配置

### mTLS 配置

启用了命名空间级别的 mTLS 策略：

```yaml
# istio/enable-mtls.yaml
# PERMISSIVE 模式（允许明文和 mTLS 混合，用于本地调试）
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: grpc-user  # 同样适用于 api-greet, api-book
spec:
  mtls:
    mode: PERMISSIVE  # 或 STRICT（生产环境）
```

### DestinationRule

配置服务间通信使用 Istio mTLS：

```yaml
# api-greet → grpc-user
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: grpc-user-mtls
  namespace: api-greet
spec:
  host: grpc-user-service.grpc-user.svc.cluster.local
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL

# api-book → api-greet
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: greet-service-mtls
  namespace: api-book
spec:
  host: greet-service.api-greet.svc.cluster.local
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL
```

## 本地开发与调试

### 使用 Telepresence

#### 1. 连接到集群
```powershell
telepresence connect
```

作用：
- 本地可以解析集群内的 Service DNS
- 本地可以直接访问集群网络

#### 2. 拦截服务流量到本地
```powershell
# 拦截 user 服务
telepresence intercept rpc-user --namespace grpc-user --port 8080:8080

# 拦截 greet 服务
telepresence intercept greet-service --namespace api-greet --port 8888:8888
```

#### 3. 本地启动服务
```powershell
# 启动 user 服务
cd user
go run user.go -f etc/user.yaml

# 启动 greet 服务
cd greet
go run greet.go -f etc/greet-api.yaml
```

#### 4. 查看拦截状态
```powershell
telepresence list
telepresence status
```

#### 5. 测试本地服务
```powershell
# 测试 user RPC
grpcurl --plaintext grpc-user-service.grpc-user:8080 user.User/Ping

# 测试 greet API
curl http://greet-service.api-greet:8888/you

# 测试 book API（会调用本地 greet，greet 再调用本地 user）
curl http://book-service.api-book:8888/info?title=Go编程
```

#### 6. 重启集群服务使拦截生效
```powershell
# 重启 greet 服务（让它重新连接 user，流量会转到本地）
kubectl rollout restart deployment api-greet -n api-greet

# 重启 book 服务（让它重新连接 greet）
kubectl rollout restart deployment api-book -n api-book
```

#### 7. 清理拦截
```powershell
telepresence leave rpc-user
telepresence quit
```

### VS Code 调试配置

`.vscode/launch.json` 已配置，支持直接在 VS Code 中调试：

- **Launch greet service** - 调试 greet 服务
- **Launch user service** - 调试 user 服务

注意：`program` 字段需指向模块目录（非单个 .go 文件）。

### mTLS 调试注意事项

如果启用了 STRICT mTLS，本地直接访问集群服务会失败。解决方案：

1. **临时改为 PERMISSIVE 模式**（推荐用于调试）：
```powershell
kubectl patch peerauthentication default -n grpc-user --type merge -p "{`"spec`":{`"mtls`":{`"mode`":`"PERMISSIVE`"}}}"
kubectl patch peerauthentication default -n api-greet --type merge -p "{`"spec`":{`"mtls`":{`"mode`":`"PERMISSIVE`"}}}"
```

2. **使用 telepresence intercept**（生产环境推荐）：
   - 流量在集群内走 mTLS，到达本地后已解密
   - 无需修改 mTLS 策略

3. **调试完成后恢复 STRICT**：
```powershell
kubectl apply -f istio/enable-mtls.yaml
```

## 部署步骤

### 1. 部署所有服务
```powershell
# 部署 user RPC
kubectl apply -f user/deploy.yaml

# 部署 greet API
kubectl apply -f greet/deploy.yaml

# 部署 book API
kubectl apply -f book/deploy.yaml
```

### 2. 应用 Istio mTLS 配置
```powershell
kubectl apply -f istio/enable-mtls.yaml
```

### 3. 重启服务使配置生效
```powershell
kubectl rollout restart deployment rpc-user -n grpc-user
kubectl rollout restart deployment api-greet -n api-greet
kubectl rollout restart deployment api-book -n api-book
```

### 4. 验证部署
```powershell
# 检查 Pod 状态（应该有 2 个容器：app + istio-proxy）
kubectl get pods -n grpc-user
kubectl get pods -n api-greet
kubectl get pods -n api-book

# 检查 Service
kubectl get svc -n grpc-user
kubectl get svc -n api-greet
kubectl get svc -n api-book

# 检查 Istio 配置
kubectl get peerauthentication -A
kubectl get destinationrule -A
```

## 测试

### 使用测试脚本
```powershell
# 默认调用 10 次，间隔 1 秒
.\test-book-api.ps1

# 自定义次数和间隔
.\test-book-api.ps1 -Count 20 -DelaySeconds 2

# 无限循环测试
while ($true) { .\test-book-api.ps1 -Count 1; Start-Sleep -Seconds 1 }
```

### 手动测试
```powershell
# 测试 Book API
curl https://book.concises.net/info?title=Go编程

# 测试 Greet API
curl https://greet.concises.net/you
curl https://greet.concises.net/me

# 在集群内测试 User RPC
kubectl exec -it <pod-name> -n api-greet -- grpcurl --plaintext grpc-user-service.grpc-user:8080 user.User/Ping
```

## 常见问题

### 1. curl 返回 "Empty reply from server"
**原因**: Istio mTLS STRICT 模式阻止了明文连接。

**解决**: 
- 临时改为 PERMISSIVE 模式
- 或使用 telepresence intercept

### 2. telepresence connect 后仍无法访问服务
**原因**: 
- mTLS STRICT 模式
- 网络路由问题

**排查**:
```powershell
# 检查 telepresence 状态
telepresence status

# 检查路由
Test-NetConnection -ComputerName <service-ip> -Port <port>

# 检查 mTLS 模式
kubectl get peerauthentication -A
```

### 3. 本地调试时断点不触发
**原因**: 集群服务使用了旧连接，未经过 intercept。

**解决**: 重启对应的客户端服务：
```powershell
kubectl rollout restart deployment <deployment-name> -n <namespace>
```

### 4. gRPC 调用失败
**检查**:
- User RPC 服务是否正常运行
- Greet 配置的 RPC 端点是否正确
- mTLS 配置是否正确

## 镜像构建

每个服务都包含 Dockerfile 和 Taskfile.yaml：

```powershell
# 使用 Task 构建镜像
cd book
task build

cd greet
task build

cd user
task build
```

## 参考资料

- [go-zero 文档](https://go-zero.dev/)
- [Istio 文档](https://istio.io/)
- [Telepresence 文档](https://www.telepresence.io/)
- [Kubernetes 文档](https://kubernetes.io/)

## License

MIT
