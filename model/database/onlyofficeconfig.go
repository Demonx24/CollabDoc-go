package database

type Document struct {
	FileType string `json:"fileType"`
	Key      string `json:"key"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

type Useronly struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EditorConfig struct {
	CallbackURL string   `json:"callbackUrl"`
	User        Useronly `json:"user"`
}

type OnlyOfficeConfig struct {
	Document     Document     `json:"document"`
	EditorConfig EditorConfig `json:"editorConfig"`
}
