package models

import (
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FullName  string      `json:"full_name"`
	Email     string      `json:"email"`
	Phone     string      `json:"phone"`
	RoleCodes []string    `json:"role_codes"`
	StoreIDs  []uuid.UUID `json:"store_ids"`
}

type UpdateUserRequest struct {
	FullName *string `json:"full_name,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}

type AssignRolesRequest struct {
	RoleCodes []string `json:"role_codes"`
}

type AssignStoresRequest struct {
	StoreIDs []uuid.UUID `json:"store_ids"`
}

type ListUsersResponse struct {
	Users      []User `json:"users"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}
