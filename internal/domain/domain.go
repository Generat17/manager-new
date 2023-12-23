package domain

type Storage map[string]Service

type Service struct {
	Type     string             `json:"type"`
	Favorite bool               `json:"favorite"`
	Elements map[string]Element `json:"elements"`
}

type Element struct {
	Password    string `json:"password"`
	Description string `json:"description"`
	Additional  string `json:"additional"`
}

type LoginBody struct {
	Login   string  `json:"login"`
	Element Element `json:"element"`
}

type ServiceBody struct {
}
