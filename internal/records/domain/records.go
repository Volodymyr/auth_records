package domain

type Record struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewRecord(title string, content string) *Record {
	return &Record{
		Title:   title,
		Content: content,
	}
}
