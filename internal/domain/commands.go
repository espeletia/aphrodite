package domain

type Message struct {
	Type string `json:"type"`
	Body Body   `json:"body"`
}

type Body struct {
	Type  string `json:"type"`
	Mode  string `json:"mode"`
	Pin   string `json:"pin"`
	Value string `json:"value"`
}
