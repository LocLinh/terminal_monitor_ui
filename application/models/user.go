package models

type MessageUser struct {
	Payload struct {
		Id     int64  `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		RoleId int64  `json:"role_id"`
		Status int64  `json:"status"`
		Config int64  `json:"config"`
	} `json:"payload"`
}
