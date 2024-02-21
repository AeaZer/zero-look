package config

import (
	"github.com/zeromicro/go-zero/zrpc"

	servicedb "github.com/zero-look/pkg/service/db"
)

type Config struct {
	zrpc.RpcServerConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	DB servicedb.Config
}
