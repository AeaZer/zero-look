package logic

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/zero-look/apps/user/rpc/internal/svc"
	"github.com/zero-look/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateTokenLogic) GenerateToken(in *user.GenerateTokenReq) (*user.GenerateTokenResp, error) {
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, accessExpire, in.UserID)
	if err != nil {
		return nil, err
	}
	return &user.GenerateTokenResp{
		AccessToken:  accessToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (l *GenerateTokenLogic) getJwtToken(secretKey string, iat, seconds, userID int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["user_id"] = userID
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
