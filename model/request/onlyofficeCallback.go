package request

type Action struct {
	Type   int    `json:"type"`
	UserID string `json:"userid"`
}

type HistoryChange struct {
	DocumentSha256 string `json:"documentSha256"`
	Created        string `json:"created"`
	User           struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
}

type CallbackRequest struct {
	Key        string   `json:"key"`
	Status     int      `json:"status"`
	Users      []string `json:"users"`
	Actions    []Action `json:"actions"`
	URL        string   `json:"url"`
	ChangesURL string   `json:"changesurl,omitempty"`
	History    struct {
		ServerVersion string          `json:"serverVersion"`
		Changes       []HistoryChange `json:"changes"`
	} `json:"history,omitempty"`
	LastSave      string `json:"lastsave,omitempty"`
	ForceSaveType int    `json:"forcesavetype,omitempty"`
	FileType      string `json:"filetype,omitempty"`
}

type CallbackResponse struct {
	Error int `json:"error"`
}
