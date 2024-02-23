package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type SentryConf struct {
	Dsn              string
	Debug            bool
	AttachStacktrace bool
}

type Config struct {
	rest.RestConf
	SentryConf SentryConf
	JwtAuth    struct {
		AccessSecret string
	}

	UserRPCConf zrpc.RpcClientConf
}
