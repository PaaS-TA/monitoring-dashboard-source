package models

type (
	UserInfo struct {
		Username      string `json:"username"`
		Password      string `json:"password"`
		Token         string `json:"token"`
		UserId        string `json:"userId"`
		UserPw        string `json:"userPw"`
		UserEmail     string `json:"userEmail"`
		UserNm        string `json:"userNm"`
		IaasUserId    string `json:"iaasUserId"`
		IaasUserPw    string `json:"iaasUserPw"`
		IaasToken     string `json:"iaasToken"`
		PaasUserId    string `json:"paasUserId"`
		PaasUserPw    string `json:"paasUserPw"`
		PaasToken     string `json:"paasToken"`
		PaasAdminYn   string `json:"paasAdminYn"`
		CaasUserId    string `json:"caasUserId"`
		CaasUserPw    string `json:"caasUserPw"`
		CaasToken     string `json:"caasToken"`
		UserAuth      string `json:"userAuth"`
		AuthI1        string `json:"authI1"`
		AuthI2        string `json:"authI2"`
		AuthP1        string `json:"authP1"`
		AuthP2        string `json:"authP2"`
		AuthC1        string `json:"authC1"`
		AuthC2        string `json:"authC2"`
		SysType       string `json:"sysType"`
		IaasUserUseYn string `json:"iaasUserUseYn"`
		PaasUserUseYn string `json:"paasUserUseYn"`
		CaasUserUseYn string `json:"caasUserUseYn"`
	}

	MemberInfo struct {
		UserId        string `json:"userId"`
		UserPw        string `json:"userPw"`
		UserEmail     string `json:"userEmail"`
		UserNm        string `json:"userNm"`
		IaasUserId    string `json:"iaasUserId"`
		IaasUserPw    string `json:"iaasUserPw"`
		PaasUserId    string `json:"paasUserId"`
		PaasUserPw    string `json:"paasUserPw"`
		CaasUserId    string `json:"caasUserId"`
		CaasUserPw    string `json:"caasUserPw"`
		IaasUserUseYn string `json:"iaasUserUseYn"`
		PaasUserUseYn string `json:"paasUserUseYn"`
		CaasUserUseYn string `json:"caasUserUseYn"`
	}
)

type ErrMessage map[string]interface{}
