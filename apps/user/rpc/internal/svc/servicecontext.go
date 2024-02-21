package svc

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/zero-look/apps/user/rpc/internal/config"
	servicedb "github.com/zero-look/pkg/service/db"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := servicedb.New(&c.DB)
	if err != nil {
		panic(fmt.Sprintf("connect mysql db errror: %v", err))
	}
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
