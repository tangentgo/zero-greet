# 循环调用 book API 脚本
param(
    [int]$Count = 10,           # 调用次数，默认 10 次
    [int]$DelaySeconds = 1      # 每次调用间隔秒数，默认 1 秒
)

$url = "https://book.concises.net/info?title=nihao"

Write-Host "开始循环调用 Book API..." -ForegroundColor Green
Write-Host "URL: $url" -ForegroundColor Cyan
Write-Host "调用次数: $Count, 间隔: $DelaySeconds 秒`n" -ForegroundColor Cyan

for ($i = 1; $i -le $Count; $i++) {
    Write-Host "[$i/$Count] 调用中..." -ForegroundColor Yellow
    
    try {
        $response = Invoke-WebRequest -Uri $url -Method Get -UseBasicParsing
        $statusCode = $response.StatusCode
        $content = $response.Content
        
        Write-Host "  状态码: $statusCode" -ForegroundColor Green
        Write-Host "  响应: $content`n" -ForegroundColor White
    }
    catch {
        Write-Host "  错误: $($_.Exception.Message)`n" -ForegroundColor Red
    }
    
    # 如果不是最后一次调用，则等待
    if ($i -lt $Count) {
        Start-Sleep -Seconds $DelaySeconds
    }
}

Write-Host "调用完成！" -ForegroundColor Green
