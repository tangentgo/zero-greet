// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	GreetAPI GreetAPIConf
}

type GreetAPIConf struct {
	Endpoint string
	Timeout  int64
}
