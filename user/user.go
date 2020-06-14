package user

type SessionUser struct {
	Uuid   string `json:"uuid,omitempty"`
	Status string `json:"status,omitempty"`
}
