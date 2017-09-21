package util

import (
	"fmt"
	"math"
	"strconv"
	cb "kr/paasta/monitoring/monit-batch/models/base"
	"time"
)


func GetInsertCurrentTime() time.Time{
	now := time.Now()
	t := now.Format(cb.INSERT_DATE_FORMAT)
	currentTime, _ := time.Parse(time.RFC3339,t)
	return currentTime
}

func GetConnectionString(host, port, user, pass , dbname string) string {

	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
		user, pass, "tcp", host, port, dbname, "")

}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func RoundFloat(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num * output)) / output
}

func RoundFloatDigit2(num float64) float64 {
	return RoundFloat(num , 2)
}

func FloattostrDigit2(fv float64) string {
	return strconv.FormatFloat(RoundFloatDigit2(fv), 'f', 2, 64)
}


func Floattostr(fv float64) string {
	return strconv.FormatFloat(fv, 'f', 2, 64)
}

func floattostrwithprec(fv float64, prec int) string {
	return strconv.FormatFloat(fv, 'f', prec, 64)
}

func isExistArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}


/*

func GetAdminToken(uaaurl string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, _ []*http.Request) error {
			return errors.New("No redirects")
		},
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	type uaaErrorResponse struct {
		Code        string `json:"error"`
		Description string `json:"error_description"`
	}

	type AuthenticationResponse struct {
		AccessToken  string           `json:"access_token"`
		TokenType    string           `json:"token_type"`
		RefreshToken string           `json:"refresh_token"`
		Error        uaaErrorResponse `json:"error"`
	}

	path := uaaurl + "/oauth/token"
	authorization := "Basic "+base64.StdEncoding.EncodeToString([]byte("cf:"))

	data := Values{
		"grant_type": {"password"},
		"scope":      {""},
		"username":	[]string{"test"},
		"password":	[]string{"test"},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encode(data)))
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	req.Header.Set("Accept", "application/json;charset=utf-8")

	resp, err := client.Do(req)
	*/
/*if err != nil {
		l4g.Error("##### err at requesting admin token:", err.Error())
		return "", err
	}*//*

	//fmt.Println("respnonse body :", resp.Body)

	returndata := new(AuthenticationResponse)
	rawdata, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rawdata, returndata)
	//fmt.Println("oauth token response :", returndata)

	return returndata.TokenType + " " + returndata.AccessToken, err
}

func GetCFToken(uaaurl string) (string, error) {
	fmt.Println("### GetCFToken called ###")
	*/
/*client := &http.Client{
		CheckRedirect: func(req *http.Request, _ []*http.Request) error {
			return errors.New("No redirects")
		},
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}*//*

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS10,
			},
		},
	}


	type uaaErrorResponse struct {
		Code        string `json:"error"`
		Description string `json:"error_description"`
	}

	type AuthenticationResponse struct {
		AccessToken  string           `json:"access_token"`
		TokenType    string           `json:"token_type"`
		RefreshToken string           `json:"refresh_token"`
		Error        uaaErrorResponse `json:"error"`
	}

	path := uaaurl + "/oauth/token"
	authorization := "Basic "+base64.StdEncoding.EncodeToString([]byte("cf:"))

	*/
/*data := Values{
		"grant_type": {"client_credentials"},
		"scope":      {""},
		"username":	[]string{"admin"},
		"password":	[]string{"admin"},
	}*//*

	data := Values{
		"grant_type": {"client_credentials"},
		"response_type": {"token"},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encode(data)))
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	req.Header.Set("Accept", "application/json;charset=utf-8")

	resp, err := client.Do(req)
	*/
/*if err != nil {
		l4g.Error("##### err at requesting cf token:", err.Error())
		return "", err
	}*//*

	//fmt.Println("respnonse body :", resp.Body)

	returndata := new(AuthenticationResponse)
	rawdata, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rawdata, returndata)

	fmt.Println("oauth token response :", returndata)

	return returndata.TokenType + " " + returndata.AccessToken, err
}

func encode(v Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := queryEscape(k) + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(queryEscape(v))
		}
	}
	return buf.String()
}
func queryEscape(s string) string {
	return escape(s, encodeQueryComponent)
}

func escape(s string, mode encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == encodeQueryComponent:
			t[j] = '+'
			j++
		case shouldEscape(c, mode):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	if mode == encodeHost {
		// §3.2.2 Host allows
		//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		// as part of reg-name.
		// We add : because we include :port as part of host.
		// We add [ ] because we include [ipv6]:port as part of host
		switch c {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']':
			return false
		}
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case encodePath: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments. This package
			// only manipulates the path as a whole, so we allow those
			// last two as well. That leaves only ? to escape.
			return c == '?'

		case encodeUserPassword: // §3.2.1
			// The RFC allows ';', ':', '&', '=', '+', '$', and ',' in
			// userinfo, so we must escape only '@', '/', and '?'.
			// The parsing of userinfo treats ':' as special so we must escape
			// that too.
			return c == '@' || c == '/' || c == '?' || c == ':'

		case encodeQueryComponent: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case encodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	// Everything else must be escaped.
	return true
}
*/
