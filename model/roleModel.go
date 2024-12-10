package model

type Role struct {
	ID_ROLE    int         `json:"ID_ROLE"`
	NAME_ROLE  string      `json:"NAME_ROLE"`
	CREATED_AT interface{} `json:"CREATED_AT"`
	UPDATED_AT interface{} `json:"UPDATED_AT"`
}
