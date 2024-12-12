package service

type Environment struct {
	ID            int64  `json:"id"`
	ApplicationID int64  `json:"ApplicationID"`
	Level         int    `json:"Level"`
	Name          string `json:"Name"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	DeletedAt     string `json:"deletedAt,omitempty"`
}
