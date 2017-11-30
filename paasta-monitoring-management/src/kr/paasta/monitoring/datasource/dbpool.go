package datasource

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	_ "github.com/go-sql-driver/mysql"
)

type DBConfigCf struct {
	URI     string
	MAXCONN int
	DBMAP   *sql.DB
}
type DBConfig struct {
	ID      string
	PWD     string
	URL     string
	MAXCONN int
	DBMAP   *sql.DB
}

func NewDBConfig(id, pwd, url string, maxConn int) *DBConfig {
	dbmap := InitDb(id, pwd, url, maxConn)

	return &DBConfig{
		ID:      id,
		PWD:     pwd,
		URL:     url,
		MAXCONN: maxConn,
		DBMAP:   dbmap,
	}
}

// 데이터베이스 연동
func InitDb(id, pwd, urlStr string, maxConn int) *sql.DB {
	databaseUrl := "mysql2://" + id + ":" + pwd + "@" + urlStr + "?reconnect=true"
	fmt.Println("Mysql databaseurl :", databaseUrl)

	url, err := url.Parse(databaseUrl)
	if err != nil {
		log.Fatalln("Error parsing DATABASE_URL", err)
	}

	targerUrl := formattedUrl(url)
	db, err := sql.Open("mysql", targerUrl)

	if err != nil {
		return nil
	}
	db.SetMaxIdleConns(maxConn)
	return db
}

//데이터 베이스 URL 포맷 생성
func formattedUrl(url *url.URL) string {
	return fmt.Sprintf(
		"%v@tcp(%v)%v?parseTime=true",
		url.User,
		url.Host,
		url.Path,
	)
}

//데이터베이스 Connection
func (config *DBConfig) getConnection() *sql.DB {
	return config.DBMAP
}

//MySql 데이터베이스 연결 종료
func (config *DBConfig) CloseDb() {
	config.DBMAP.Close()
}
