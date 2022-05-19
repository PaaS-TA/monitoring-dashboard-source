package v1service

import (
	dao "GoEchoProject/dao/api/v1"
	"GoEchoProject/models/api/v1"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/twinj/uuid"
	"net/http"
	"os"
	"strings"
	"time"
)

//Gorm Object Struct
type TokenService struct {
	DbInfo    *gorm.DB
	RedisInfo *redis.Client
}

func GetTokenService(DbInfo *gorm.DB, RedisInfo *redis.Client) *TokenService {
	return &TokenService{
		DbInfo:    DbInfo,
		RedisInfo: RedisInfo,
	}
}

// 1. 토큰 추출
func ExtractToken(c echo.Context) (string, error) {
	var err error
	req := c.Request()
	bearToken := req.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		// Bearer 제외
		return strArr[1], err
	} else if len(bearToken) == 0 || len(strArr) == 0 {
		// Token 정보가 없으면
		return "", fmt.Errorf("Please enter token")
	}
	return bearToken, err
}

// 2. 토큰 검증 (signing method 검증, 서명 검증)
func VerifyToken(bearToken string, secretType string, c echo.Context) (*jwt.Token, error) {
	token, err := jwt.Parse(bearToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(secretType)), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// 3. 토큰 만료 검증
func TokenValid(token *jwt.Token) error {
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return fmt.Errorf("Invaild token")
	}
	return nil
}

// 4. 메타 데이터 추출 (메타데이터 이용한 Redis 확인)
func ExtractTokenMetadata(token *jwt.Token, tokenType string) (map[string]interface{}, error) {
	// metadata 초기화 선언
	metadata := make(map[string]interface{})
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		switch tokenType {
		case "ACCESS":
			accessUuid, ok := claims["access_uuid"].(string)
			if !ok {
				return metadata, fmt.Errorf("No access_uuid in metadata")
			}
			userId, ok := claims["user_id"].(string)
			if !ok {
				return metadata, fmt.Errorf("No user_id in metadata")
			}
			metadata["access_uuid"] = accessUuid
			metadata["user_id"] = userId
			return metadata, nil
		case "REFRESH":
			refreshUuid, ok := claims["refresh_uuid"].(string)
			if !ok {
				return metadata, fmt.Errorf("No refresh_uuid in metadata")
			}
			userId, ok := claims["user_id"].(string)
			if !ok {
				return metadata, fmt.Errorf("No user_id in metadata")
			}
			metadata["refresh_uuid"] = refreshUuid
			metadata["user_id"] = userId
			return metadata, nil
		}
	}
	return metadata, fmt.Errorf("Token is invalid")
}

// 5. Redis 추출 (UUID를 이용한 userId 추출)
func FetchAuth(metadata map[string]interface{}, RedisInfo *redis.Client) (string, error) {
	userid, err := RedisInfo.Get(metadata["access_uuid"].(string)).Result()
	if err != nil {
		return "", fmt.Errorf("User is not exist")
	}
	return userid, nil
}

// redis Token 데이터 삭제
func DeleteAuth(givenUuid string, RedisInfo *redis.Client) (int64, error) {
	deleted, err := RedisInfo.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// Token 데이터 생성
func CreateToken(td v1.TokenDetails, user_id string) (v1.TokenDetails, error) {
	var err error

	// Access Token 만료 시간 (현재시간 + 15분)
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()
	// Refresh Token 만료 시간 (현재시간 + 7일)
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = user_id
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = user_id
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))

	if err != nil {
		return td, err
	}
	return td, nil
}

// Token 데이터 Redis 저장
func CreateAuth(td v1.TokenDetails, user_id string, RedisInfo *redis.Client) (v1.TokenDetails, error) {
	// at는 AccessToken의 접근 유효 시간
	// rt는 RefreshToken의 만료 시간
	redis_at := time.Unix(td.AtExpires, 0) //converting Unix to UTC
	redis_rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	// JWT을 Redis에 저장한다.
	errAccess := RedisInfo.Set(td.AccessUuid, user_id, redis_at.Sub(now)).Err()
	if errAccess != nil {
		return td, errAccess
	}
	errRefresh := RedisInfo.Set(td.RefreshUuid, user_id, redis_rt.Sub(now)).Err()
	if errRefresh != nil {
		return td, errRefresh
	}
	return td, nil
}

// Refresh 토큰 정리 후 적용
func (h *TokenService) CreateToken(apiRequest v1.CreateToken, c echo.Context) (v1.TokenDetails, error) {
	// 아이디 및 비밀번호 확인 시 JWT 토큰 발급 및 Redis 저장

	// 1. Token 모델을 선언한다.
	td := v1.TokenDetails{}

	userInfo := v1.UserInfo{
		Username: apiRequest.Username,
		Password: apiRequest.Password,
	}

	// 2. 전달 받은 계정 정보로 데이터베이스에 계정이 존재하는지 확인한다.
	results, err := dao.GetUserDao(h.DbInfo).GetUser(userInfo, c)
	if err != nil {
		return td, err
	}
	if len(results) == 0 {
		return td, fmt.Errorf("reason: cannot found username")
	}
	// 계정에 대한 비밀번호를 확인한다.
	if results[0].Password != apiRequest.Password {
		return td, fmt.Errorf("reason: password is incorrect")
	}

	// 3. Token 생성
	td, err = CreateToken(td, results[0].Username)

	// 4. Token 저장 (Redis)
	td, err = CreateAuth(td, results[0].Username, h.RedisInfo)
	return td, nil
}

func (h *TokenService) RefreshToken(apiRequest v1.RefreshToken, c echo.Context) (v1.TokenDetails, error) {
	// RefreshToken 확인 시 기존 JWT 토큰 정보 삭제 및 생성 후 Redis 저장

	// 1. Token 모델을 선언한다.
	td := v1.TokenDetails{}

	// 2. 토큰 검증 (signing method 검증, 서명 검증)
	token, err := VerifyToken(apiRequest.RefreshToken, "REFRESH_SECRET", c)
	if err != nil {
		return td, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// 3. 토큰 만료 검증 (동작 방식 확인 필요)
	err = TokenValid(token)
	if err != nil {
		return td, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// 4. 메타 데이터 추출 (메타데이터 이용한 Redis 확인)
	metadata, err := ExtractTokenMetadata(token, "REFRESH")
	if err != nil {
		return td, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(metadata)

	// 5. Redis Token 삭제
	//Delete the previous Refresh Token
	deleted, delErr := DeleteAuth(metadata["refresh_uuid"].(string), h.RedisInfo)
	if delErr != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return td, err
	}

	// 6. Token 생성
	td, err = CreateToken(td, metadata["user_id"].(string))
	if err != nil {
		return td, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 7. Token 저장 (Redis)
	td, err = CreateAuth(td, metadata["user_id"].(string), h.RedisInfo)
	if err != nil {
		return td, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return td, nil
}
