package model

type Vocab struct {
	Hanzi        string              `json:"hanzi"`
	ID           int                 `json:"id"`
	Level        int                 `json:"level"`
	Pinyin       string              `json:"pinyin"`
	Radicals     string              `json:"radicals"`
	Strokes      string              `json:"strokes"`
	Translations map[string][]string `json:"translations"`
}
