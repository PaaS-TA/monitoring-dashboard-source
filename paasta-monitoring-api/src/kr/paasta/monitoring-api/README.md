### Environments
1. IDE : Goland
2. GOROOT : Go 1.17.8
3. GOPATH : 'Use GOPATH that's defined in system environment' 체크
4. GO Modules : 'Enable Go modules integration' 체크
5. `...\paasta-monitoring-api\src\kr\paasta\monitoring-api`에서 `go get -u` 실행
6. Swagger 실행 파일(`swag.exe`) 설치 : `go get github.com/swaggo/swag/cmd/swag` 실행

### Run/Debug Configurations
1. Name : paasta-monitoring-api
2. Runkind : File
3. Files : `...\paasta-monitoring-api\src\kr\paasta\monitoring-api\main.go`
4. Output directory : `...\paasta-monitoring-api\src\kr\paasta\monitoring-api`
5. Working directory : `...\paasta-monitoring-api\src\kr\paasta\monitoring-api`
6. Before launch :  
   1\. Add > Run External tool  
   2\. Add > Create Tool  
   . Name : swag init  
   . Description : swag init  
   . Program : `$GOPATH\bin\swag.exe`  
   . Arguments : init  
   . Working directory : `...\paasta-monitoring-api\src\kr\paasta\monitoring-api`

### Swagger Annotations Example Form
```go
// CreateToken
//  * Annotations for Swagger *
//  @Summary      [테스트] 토큰 생성하기
//  @Description  [테스트] 토큰 정보를 생성한다.
//  @Accept       json
//  @Produce      json
//  @Param        UserInfo  body      CreateToken  true  "Insert UserInfo"
//  @Success      200       {object}  apiHelpers.BasicResponseForm{responseInfo=TokenDetails}
//  @Router       /api/v1/token [post]
func (a *TokenController) CreateToken(c echo.Context) (err error) {
    ...
}
```
