package utils

import (
	"fmt"
	"encoding/json"
)

/*************
	 Utils
**************/
func PrintJson(p interface{}) {
	doc, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\n\n============  RESULT ============\n")
	fmt.Println(string(doc))
}