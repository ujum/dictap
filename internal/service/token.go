package service

import (
	ctx "context"
	"encoding/json"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/internal/dto"
	"github.com/ujum/dictap/pkg/logger"
	"net/http"
	"time"
)

type TokenService interface {
	Generate(requestCtx ctx.Context, credentials *dto.UserCredentials) (*dto.TokenDTO, error)
	Refresh(requestCtx ctx.Context, refreshToken json.RawMessage) (*dto.TokenDTO, error)
	GetClaimsTypeVerifyFunc() func() interface{}
}

type JwtTokenService struct {
	log         logger.Logger
	cfg         *config.Config
	signer      *jwt.Signer
	verifier    *jwt.Verifier
	userService UserService
}

type userClaims struct {
	Uid string `json:"uid,required"`
	App string `json:"app"`
}

func NewJwtSigner(cfg *config.Config) (*jwt.Signer, error) {
	privateKeyRSA, err := jwt.LoadPrivateKeyRSA(cfg.ConfigDir + "/rsa_private_key.pem")
	if err != nil {
		return nil, err
	}
	min := cfg.Server.Security.ApiKeyAuth.AccessTokenMaxAgeMin
	signer := jwt.NewSigner(jwt.RS256, privateKeyRSA, time.Duration(min)*time.Minute)
	return signer, err
}

func NewJwtVerifier(cfg *config.Config) (*jwt.Verifier, error) {
	publicKeyRSA, err := jwt.LoadPublicKeyRSA(cfg.ConfigDir + "/rsa_public_key.pem")
	if err != nil {
		return nil, err
	}
	verifier := jwt.NewVerifier(jwt.RS256, publicKeyRSA)
	verifier.WithDefaultBlocklist()
	verifier.ErrorHandler = func(ctx *context.Context, err error) {
		ctx.StopWithJSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
	}
	return verifier, nil
}

func NewJwtTokenService(cfg *config.Config, appLogger logger.Logger,
	verifier *jwt.Verifier, signer *jwt.Signer, userService UserService) *JwtTokenService {
	tokenService := &JwtTokenService{
		log:         appLogger,
		cfg:         cfg,
		userService: userService,
		signer:      signer,
		verifier:    verifier,
	}
	return tokenService
}

func (tokenSrv *JwtTokenService) Generate(requestCtx ctx.Context, credentials *dto.UserCredentials) (*dto.TokenDTO, error) {
	user, err := tokenSrv.userService.GetByCredentials(requestCtx, credentials)
	if err != nil {
		return nil, err
	}
	return tokenSrv.generate(user)
}

func (tokenSrv *JwtTokenService) Refresh(requestCtx ctx.Context, refreshToken json.RawMessage) (*dto.TokenDTO, error) {
	verifiedToken, err := tokenSrv.verifier.VerifyToken(refreshToken)
	if err != nil {
		tokenSrv.log.Errorf("verify refresh token: %v", err)
		return nil, err
	}
	uid, err := tokenSrv.userService.GetByUid(requestCtx, verifiedToken.StandardClaims.Subject)
	if err != nil {
		return nil, err
	}

	return tokenSrv.generate(uid)
}

func (tokenSrv *JwtTokenService) generate(user *domain.User) (*dto.TokenDTO, error) {
	refreshClaims := jwt.Claims{Subject: user.Uid}
	accessClaims := userClaims{
		Uid: user.Uid,
		App: "Dictup",
	}
	refreshMin := tokenSrv.cfg.Server.Security.ApiKeyAuth.RefreshTokenMaxAgeMin
	tokenPair, err := tokenSrv.signer.NewTokenPair(accessClaims, refreshClaims, time.Duration(refreshMin)*time.Minute)
	if err != nil {
		tokenSrv.log.Errorf("token pair generation error: %v", err)
		return nil, err
	}
	return &dto.TokenDTO{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (tokenSrv *JwtTokenService) GetClaimsTypeVerifyFunc() func() interface{} {
	return func() interface{} {
		return new(userClaims)
	}
}
