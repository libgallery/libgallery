package rule34

type post struct {
	Score     string `xml:"score,attr"`
	FileURL   string `xml:"file_url,attr"`
	Tags      string `xml:"tags,attr"`
	ID        string `xml:"id,attr"`
	CreatedAt string `xml:"created_at,attr"`
	Source    string `xml:"source,attr"`
}

type searchResponse struct {
	Success *bool  `xml:"success,attr"`
	Error   error  `xml:"reason,attr"`
	Posts   []post `xml:"post"`
}
