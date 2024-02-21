package logic

import (
	"context"

	"github.com/zero-look/apps/user/api/internal/svc"
	"github.com/zero-look/apps/user/api/internal/types"
	"github.com/zero-look/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) error {
	_, err := l.svcCtx.UserRpc.Register(l.ctx, &userclient.RegisterReq{
		UserID:    req.UserID,
		StaffName: req.StaffName,
		Email:     req.Email,
		Password:  req.Password,
	})
	return err
}
