# 测试完整的追踪链路

Write-Host "=== 测试 Istio 链路追踪 ===" -ForegroundColor Cyan
Write-Host ""

# 1. 检查服务状态
Write-Host "1. 检查服务状态:" -ForegroundColor Yellow
kubectl get pods -n api-book -l app=api-book
kubectl get pods -n api-greet -l app=api-greet
Write-Host ""

# 2. 检查 Telemetry 配置
Write-Host "2. 检查 Telemetry 配置:" -ForegroundColor Yellow
kubectl get telemetry -A
Write-Host ""

# 3. 直接测试 greet 服务
Write-Host "3. 直接访问 greet 服务（生成追踪）:" -ForegroundColor Yellow
$response1 = curl -s https://greet.concises.net/you -X GET
Write-Host "Response: $response1"
Write-Host ""

# 4. 通过 book 服务测试（生成完整链路）
Write-Host "4. 通过 book 服务访问（完整调用链）:" -ForegroundColor Yellow
$body = @{
    title = "Istio 追踪测试 $(Get-Date -Format 'HH:mm:ss')"
} | ConvertTo-Json

$response2 = curl -s https://book.concises.net/bookinfo `
    -X POST `
    -H "Content-Type: application/json" `
    -d $body

Write-Host "Response: $response2"
Write-Host ""

# 5. 查看最近的日志
Write-Host "5. 查看 greet 服务最近的请求日志:" -ForegroundColor Yellow
kubectl logs -n api-greet -l app=api-greet -c api-greet --tail=3 | Select-String "GET /you"
Write-Host ""

Write-Host "6. 查看 book 服务最近的请求日志:" -ForegroundColor Yellow
kubectl logs -n api-book -l app=api-book -c api-book --tail=3 | Select-String "POST /bookinfo"
Write-Host ""

# 7. 提示如何查看 Jaeger
Write-Host "=== 在 Jaeger UI 中查看追踪 ===" -ForegroundColor Green
Write-Host ""
Write-Host "方式1: 直接访问域名" -ForegroundColor Cyan
Write-Host "  https://jaeger.concises.net" -ForegroundColor White
Write-Host ""
Write-Host "方式2: 使用 port-forward" -ForegroundColor Cyan
Write-Host "  kubectl port-forward -n istio-system svc/tracing 16686:80" -ForegroundColor White
Write-Host "  然后访问: http://localhost:16686" -ForegroundColor White
Write-Host ""
Write-Host "在 Jaeger UI 中查找服务:" -ForegroundColor Yellow
Write-Host "  1. Service 下拉框可能的名称:" -ForegroundColor White
Write-Host "     - istio-ingressgateway.istio-system" -ForegroundColor Gray
Write-Host "     - book-service.api-book" -ForegroundColor Gray
Write-Host "     - greet-service.api-greet" -ForegroundColor Gray
Write-Host "     - api-book.api-book" -ForegroundColor Gray
Write-Host "     - api-greet.api-greet" -ForegroundColor Gray
Write-Host ""
Write-Host "  2. 点击 'Find Traces' 按钮" -ForegroundColor White
Write-Host "  3. 查看最近的追踪记录" -ForegroundColor White
Write-Host ""
Write-Host "提示: greet 服务可能显示为 'greet-service.api-greet' 或 'api-greet.api-greet'" -ForegroundColor Magenta
