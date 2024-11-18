package dtos

type Account struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	IsActive bool   `json:"isActive"`
}

type UpdateAccountRequest struct {
	ID       string  `json:"id"`
	Login    *string `json:"login,omitempty"`
	Password *string `json:"password,omitempty"`
	IsActive *bool   `json:"isActive,omitempty"`
}
