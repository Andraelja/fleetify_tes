package responses

type RegisterResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginResponse struct {
	Success bool        `json:"success"`
	Token   string      `json:"token"`
	Data    interface{} `json:"data"`
}
