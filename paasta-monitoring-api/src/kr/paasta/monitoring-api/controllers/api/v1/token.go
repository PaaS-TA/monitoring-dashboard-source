package v1

import (
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	"paasta-monitoring-api/helpers"
	"paasta-monitoring-api/models/api/v1"
	v1service "paasta-monitoring-api/services/api/v1"
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

// CreateToken
//  * Annotations for Swagger *
//  @Summary      토큰 생성하기
//  @Description  토큰 정보를 생성한다.
//  @Accept       json
//  @Produce      json
//  @Param        UserInfo  body      CreateToken  true  "토큰을 생성하기 위해 필요한 유저 정보를 제공한다."
//  @Success      200       {object}  apiHelpers.BasicResponseForm{responseInfo=TokenDetails}
//  @Router       /api/v1/token [post]
func (a *TokenController) CreateToken(c echo.Context) (err error) {
	var userRequest v1.CreateToken                          // 클라이언트의 리퀘스트 정보를 저장할 변수 선언
	err = helpers.BindRequestAndCheckValid(c, &userRequest) // 클라이언트의 리퀘스트 정보의 바인딩 & 유효성 결과를 반환
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	// Authentication의 CreateToken 발급을 호출한다.
	tokenDetails, err := v1service.GetTokenService(a.DbInfo, a.RedisInfo).CreateToken(userRequest, c)
	if err != nil {
		apiHelpers.Respond(c, http.StatusInternalServerError, "Cannot create token", err.Error())
		return err
	}

	// TokenDetails로 토근 정보를 전달한다. (AccessToken, RefreshToken, AccessUuid, RefreshUuid, AtExpires, RtExpires)
	apiHelpers.Respond(c, http.StatusOK, "Success to create token", tokenDetails)
	return nil
}


func (a *TokenController) CreateAccessToken(ctx echo.Context) error {
	var params v1.TokenParam
	err := helpers.BindRequestAndCheckValid(ctx, &params)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	tokenMap, err := v1service.GetTokenService(a.DbInfo, a.RedisInfo).CreateAccessToken(params, ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusInternalServerError, "Cannot create token", err.Error())
		return err
	}

	// TokenDetails로 토근 정보를 전달한다. (AccessToken, RefreshToken, AccessUuid, RefreshUuid, AtExpires, RtExpires)
	apiHelpers.Respond(ctx, http.StatusOK, "Success to create token", tokenMap)
	return nil
}


// RefreshToken
//  * Annotations for Swagger *
//  @Summary      토큰 리프레시하기
//  @Description  토큰 정보를 리프레시한다.
//  @Accept       json
//  @Produce      json
//  @Param        TokenDetails  body      RefreshToken  true  "토큰을 리프레시하기 위한 토큰 정보를 제공한다."
//  @Success      200           {object}  apiHelpers.BasicResponseForm{responseInfo=TokenDetails}
//  @Router       /api/v1/token [put]
func (a *TokenController) RefreshToken(c echo.Context) (err error) {
	var userRequest v1.RefreshToken                         // 클라이언트의 리퀘스트 정보를 저장할 변수 선언
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

	// TokenDetails로 토근 정보를 전달한다. (AccessToken, RefreshToken, AccessUuid, RefreshUuid, AtExpires, RtExpires)
	apiHelpers.Respond(c, http.StatusOK, "Success to refresh token", tokenDetails)
	return nil
}
