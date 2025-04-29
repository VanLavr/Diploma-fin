package dto

type CreateGroupDTO struct {
	Name string `json:"name"`
}

type UpdateGroupDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
