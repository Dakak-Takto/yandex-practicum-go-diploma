package handlers

type UserRegisterDTO struct {
	Login    string
	Password string
}

type UserLoginDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type WithdrawDTO struct {
	Order string
	Sum   float64
}
