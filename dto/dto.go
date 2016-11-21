package dto

type (
	APIResult struct {
		Result  interface{} `json:"result"`
		Success bool        `json:"success"`
		Error   APIError    `json:"error"`
	}

	APIError struct {
		Code    int         `json:"code"`
		Details interface{} `json:"details"`
		Message string      `json:"message"`
	}
)
