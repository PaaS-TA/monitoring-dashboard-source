package model

type (

	UserInfo struct {
		Username 		string   `json:"username"`
		Password 		string   `json:"password"`
		Token    		string   `json:"token"`
		UserId         	string   `json:"userId"`
		UserPw         	string   `json:"userPw"`
		UserEmail      	string   `json:"userEmail"`
		UserNm         	string   `json:"userNm"`
		IaasUserId     	string   `json:"iaasUserId"`
		IaasUserPw     	string   `json:"iaasUserPw"`
		IaasToken    	string   `json:"iaasToken"`
		PaasUserId      string   `json:"paasUserId"`
		PaasUserPw      string   `json:"paasUserPw"`
		PaasToken    	string   `json:"paasToken"`
		PaasAdminYn    	string   `json:"paasAdminYn"`
		UserAuth    	string   `json:"userAuth"`
		AuthI1    	    string   `json:"authI1"`
		AuthI2    	    string   `json:"authI2"`
		AuthP1    	    string   `json:"authP1"`
		AuthP2    	    string   `json:"authP2"`
		SysType    	    string   `json:"sysType"`
		IaasUserUseYn   string   `json:"iaasUserUseYn"`
		PaasUserUseYn   string   `json:"paasUserUseYn"`
	}

	MemberInfo struct{
		UserId         	string   `json:"userId"`
		UserPw         	string   `json:"userPw"`
		UserEmail      	string   `json:"userEmail"`
		UserNm     	    string   `json:"userNm"`
		IaasUserId     	string   `json:"iaasUserId"`
		IaasUserPw     	string   `json:"iaasUserPw"`
		PaasUserId      string   `json:"paasUserId"`
		PaasUserPw      string   `json:"paasUserPw"`
		IaasUserUseYn   string   `json:"iaasUserUseYn"`
		PaasUserUseYn   string   `json:"paasUserUseYn"`
	}

)

type ErrMessage map[string]interface{}