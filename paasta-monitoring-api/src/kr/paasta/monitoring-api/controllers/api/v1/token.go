package v1

import (
	"GoEchoProject/apiHelpers"
	"GoEchoProject/connections"
	"GoEchoProject/helpers"
	"GoEchoProject/models/api/v1"
	v1service "GoEchoProject/services/api/v1"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TokenController struct {
	DbInfo    *gorm.DB
	RedisInfo *redis.Client
}

func GetTokenController(conn connections.Connections) *TokenController {
	return &TokenController{
		DbInfo:    conn.DbInfo,
		RedisInfo: conn.RedisInfo,
	}
}

func (a *TokenController) CreateToken(c echo.Context) (err error) {
	var userRequest v1.UserInfo                             // 클라이언트의 리퀘스트 정보를 저장할 변수 선언
	err = helpers.BindRequestAndCheckValid(c, &userRequest) // 클라이언트의 리퀘스트 정보의 바인딩 & 유효성 결과를 반환
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	// Authentication의 CreateToken 발급을 호출한다.
	tokenDetails, err := v1service.GetTokenService(a.DbInfo, a.RedisInfo).CreateToken(userRequest, c)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Cannot create token", err.Error())
		return err
	}

	// TokenDetails로 토근 정보를 전달한다. (AccessToken, RefreshToken, AccessUuid, RefreshUuid, AtExpires, RtExpires)
	apiHelpers.Respond(c, http.StatusOK, "Success to create token", tokenDetails)
	return nil
}

func (a *TokenController) RefreshToken(c echo.Context) (err error) {
	var userRequest v1.TokenDetails                         // 클라이언트의 리퀘스트 정보를 저장할 변수 선언
	err = helpers.BindRequestAndCheckValid(c, &userRequest) // 클라이언트의 리퀘스트 정보의 바인딩 & 유효성 결과를 반환
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	// Authentication의 RefreshToken 발급을 호출한다.
	tokenDetails, err := v1service.GetTokenService(a.DbInfo, a.RedisInfo).RefreshToken(userRequest, c)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Cannot refresh token", err.Error())
		return err
	}

	// TokenDetails로 토근 정보를 전달한다. (AccessToken, RefreshToken, AccessUuid, RefreshUuid, AtExpires, RtExpires) .
	apiHelpers.Respond(c, http.StatusOK, "Success to refresh token", tokenDetails)
	return nil
}
