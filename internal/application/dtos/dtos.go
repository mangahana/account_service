package dtos

type RegisterInput struct {
	Phone string `json:"phone"`
	IP    string // client ip address
}

type ConfirmCodeInput struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type CompleteRegisterInput struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginInput struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type AuthOutput struct {
	AccessToken string `json:"access_token"`
}

type RecoveryInput struct {
	Phone string `json:"phone"`
	IP    string
}

type CompleteRecovery struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Password string `json:"password"`
}
