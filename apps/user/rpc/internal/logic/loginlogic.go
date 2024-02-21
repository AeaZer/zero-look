package logic

import (
	"context"
	"errors"

	"github.com/zero-look/apps/user/rpc/internal/svc"
	perfUser "github.com/zero-look/apps/user/rpc/models/user"
	"github.com/zero-look/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	u, err := perfUser.FindOneForLogin(l.svcCtx.DB, in.GetStaffName(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("rpc not exists")
	}
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	token, err := generateTokenLogic.GenerateToken(&user.GenerateTokenReq{
		UserID: u.UserID,
	})
	if err != nil {
		return nil, err
	}
	return &user.LoginResp{
		UserID:       u.UserID,
		StaffName:    u.StaffName,
		AccessToken:  token.AccessToken,
		AccessExpire: token.AccessExpire,
		RefreshAfter: token.RefreshAfter,
	}, nil
}
