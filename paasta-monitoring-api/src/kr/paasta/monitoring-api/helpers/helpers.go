package helpers

import (
	"fmt"
	"strconv"
)

//Int64ToString function convert a float number to a string
func Int64ToString(inputNum int64) string {
	return strconv.FormatInt(inputNum, 10)
}

func GetDBConnectionString(dbtype, user, password, protocol, host, port, dbname, charset, parseTime string) (string, string) {
	return dbtype, fmt.Sprintf("%s:%s@%s([%s]:%s)/%s?charset=%s&parseTime=%s",
		user, password, protocol, host, port, dbname, charset, parseTime)
}
