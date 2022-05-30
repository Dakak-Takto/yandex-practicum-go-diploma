package handlers

type (
	userRegRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	userLoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	withdrawRequest struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}

	balanceResponse struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}
)
