package main

import (
	"flag"
	"fmt"

	"github.com/zero-look/apps/user/api/internal/config"
	"github.com/zero-look/apps/user/api/internal/handler"
	"github.com/zero-look/apps/user/api/internal/svc"
	"github.com/zero-look/pkg/common/api"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

// 业务系统编号
const bizNumber = 1000

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	api.SetGlobalErrorHandler(bizNumber)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
