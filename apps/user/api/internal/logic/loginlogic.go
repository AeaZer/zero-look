package logic

import (
	"context"

	"github.com/zero-look/apps/user/api/internal/svc"
	"github.com/zero-look/apps/user/api/internal/types"
	"github.com/zero-look/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	r, err := l.svcCtx.UserRpc.Login(l.ctx, &userclient.LoginReq{
		StaffName: req.StaffName,
		Password:  req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &types.LoginResp{
		UserID:       r.UserID,
		Email:        r.Email,
		AccessToken:  r.AccessToken,
		AccessExpire: r.AccessExpire,
		RefreshAfter: r.RefreshAfter,
	}, nil
}
