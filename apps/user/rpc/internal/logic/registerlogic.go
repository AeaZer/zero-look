package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-look/apps/user/rpc/internal/svc"
	zeroLookUser "github.com/zero-look/apps/user/rpc/models/user"
	"github.com/zero-look/apps/user/rpc/user"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	zUser := &zeroLookUser.User{
		UserID:    in.UserID,
		StaffName: in.StaffName,
		Email:     in.Email,
		Password:  in.Password,
	}
	err := l.svcCtx.DB.Save(zUser).Error
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	return &user.RegisterResp{}, nil
}
