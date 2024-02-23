package main

import (
	"flag"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap/zapcore"

	"github.com/zero-look/apps/user/api/internal/config"
	"github.com/zero-look/apps/user/api/internal/handler"
	"github.com/zero-look/apps/user/api/internal/svc"
	"github.com/zero-look/pkg/common/api"
	"github.com/zero-look/pkg/zapx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

// 业务系统编号
const (
	bizNumber = 1000
	bizName   = "user-api"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	err := setLoggerWriter(&c)
	if err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)
	api.SetGlobalErrorHandler(bizNumber)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	logx.Errorf("this is a test")
	logx.Debug("this is a test")
	server.Start()
}

func setLoggerWriter(c *config.Config) error {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              c.SentryConf.Dsn,
		Debug:            c.SentryConf.Debug,
		AttachStacktrace: c.SentryConf.AttachStacktrace,
		ServerName:       bizName,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		return err
	}
	zapCore := zapx.SentryCoreConfig{
		Tags:              map[string]string{"env": "local"},
		DisableStacktrace: true,
		Level:             -1,
		FlushTimeout:      0,
	}
	sentryCore := zapx.NewSentryCore(zapCore, client)
	optConf := zapx.OptConf{
		ZapCores: []zapcore.Core{sentryCore},
	}
	writer, err := zapx.NewZapWriter(optConf)
	logx.Must(err)
	logx.SetWriter(writer)

	return nil
}
