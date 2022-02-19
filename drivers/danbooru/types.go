package danbooru

type post struct {
	ID           uint   `json:"id"`
	CreatedAt    string `json:"created_at"`
	UploaderID   uint   `json:"uploader_id"`
	Score        int    `json:"score"`
	Source       string `json:"source"`
	Rating       string `json:"rating"`
	Tags         string `json:"tag_string"`
	LargeFileURL string `json:"large_file_url,omitempty"`
}

type comments []struct {
	ID        uint   `json:"id"`
	PostID    uint   `json:"post_id"`
	CreatorID uint   `json:"creator_id"`
	Body      string `json:"body"`
	Score     int64  `json:"score"`
	CreatedAt string `json:"created_at"`
	IsDeleted bool   `json:"is_deleted"`
	IsSticky  bool   `json:"is_sticky"`
}
