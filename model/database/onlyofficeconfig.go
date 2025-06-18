package database

type OnlyOfficeConfig struct {
	Document struct {
		FileType    string `json:"fileType"`
		Key         string `json:"key"`
		Title       string `json:"title"`
		URL         string `json:"url"`
		Permissions struct {
			Edit bool `json:"edit"`
		} `json:"permissions"`
	} `json:"document"`
	EditorConfig struct {
		CallbackURL string `json:"callbackUrl"`
		User        struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
	} `json:"editorConfig"`
}
