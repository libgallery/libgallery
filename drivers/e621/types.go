package e621

type file struct {
	URL string `json:"url"`
}

type post struct {
	ID          uint64   `json:"id"`
	Tags        tags     `json:"tags"`
	CreatedAt   string   `json:"created_at"`
	File        file     `json:"file"`
	Description string   `json:"description"`
	UploaderID  uint64   `json:"uploader"`
	Rating      string   `json:"rating"`
	Source      []string `json:"sources"`
	Score       struct {
		Total int64 `json:"total"`
	} `json:"score"`
}

type tags struct {
	General   []string `json:"general"`
	Species   []string `json:"species"`
	Character []string `json:"character"`
	Copyright []string `json:"copyright"`
	Artist    []string `json:"artist"`
	Lore      []string `json:"lore"`
	Meta      []string `json:"meta"`
}
