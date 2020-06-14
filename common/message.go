package common

type Message struct {
	Datetime string   `json:"Datetime,omitempty"`
	Tags     []string `json:"Tags,omitempty"`
	TimeFrom string   `json:"TimeFrom,omitempty"`
	TimeTo   string   `json:"TimeTo,omitempty"`
	Message  string   `json:"Message,omitempty"`
	Token    string   `json:"Token,omitempty"`
}

type Response struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type UserLogin struct {
	Username string `json:"Username,omitempty"`
	Password string `json:"Password,omitempty"`
	Token    string `json:"Token,omitempty"`
}
