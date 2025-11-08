// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"net/http"
	"time"

	"book/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	GreetClient *http.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	timeout := time.Duration(c.GreetAPI.Timeout) * time.Millisecond
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	// 使用标准 HTTP 客户端
	// go-zero 会自动通过 httpx 传播追踪头
	return &ServiceContext{
		Config: c,
		GreetClient: &http.Client{
			Timeout: timeout,
		},
	}
}
