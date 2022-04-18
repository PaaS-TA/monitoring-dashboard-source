package v1

import (
	"GoEchoProject/connections"
	"GoEchoProject/models"
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

	/* Request Body Data를 매핑한다. */
	var apiRequest models.UserInfo // -> &추가
	if err = c.Bind(&apiRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid json provided")
	}

	// Authentication의 CreateToken 발급을 호출한다.
	tokenDetails, err := v1service.GetTokenService(a.DbInfo, a.RedisInfo).CreateToken(apiRequest, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TokenDetails로 토근 정보를 전달한다. (AccessToken, RefreshToken, AccessUuid, RefreshUuid, AtExpires, RtExpires)
	// JSON 같은 경우 Unhandled error에 대한 처리를 필요로 한다.
	err = c.JSON(http.StatusOK, tokenDetails)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (a *TokenController) RefreshToken(c echo.Context) (err error) {

	/* Request Body Data를 매핑한다.  */
	var apiRequest models.TokenDetails // -> &추가
	if err = c.Bind(&apiRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid json provided")
	}

	// Authentication의 CreateToken 발급을 호출한다.
	tokenDetails, err := v1service.GetTokenService(a.DbInfo, a.RedisInfo).RefreshToken(apiRequest, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TokenDetails로 토근 정보를 전달한다. (AccessToken, RefreshToken, AccessUuid, RefreshUuid, AtExpires, RtExpires)
	// JSON 같은 경우 Unhandled error에 대한 처리를 필요로 한다.
	err = c.JSON(http.StatusOK, tokenDetails)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
