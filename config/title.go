package config

type Question struct {
	Title    string   `json:"title"`
	Options  []string `json:"options"`
	IsLiKeTe bool     `json:"isLiKeTe"`
	IsMulti  bool     `json:"isMulti"`
	IsFill   bool     `json:"isFill"`
}

type Questions struct {
	Questions []*Question `json:"questions"`
}
