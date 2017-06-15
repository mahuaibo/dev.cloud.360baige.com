package response

type Response struct {
	Code    string        `json:"code"`
	Messgae string        `json:"messgae"`
	Data    interface{}   `json:"data,omitempty"`
}

