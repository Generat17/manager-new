package domain

type Storage map[string]Element

type Element struct {
	Type        string `json:"type"`
	Password    string `json:"password"`
	Description string `json:"description"`
	Additional  string `json:"additional"`
	Favorite    bool   `json:"favorite"`
}
