package handlers

type UserRegisterDTO struct {
	Login    string
	Password string
}

type UserLoginDTO struct {
	Login    string
	Password string
}

type WithdrawDTO struct {
	Order int
	Sum   int
}
