package dto

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	Limit        int `json:"limit"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}

type ResponseUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Photo string `json:"photo"`
	Token string `json:"token"`
}

type ResponValidatePhone struct {
	Status string `json:"status"`
	Phone  string `json:"phone"`
	Vendor string `json:"vendor"`
}
