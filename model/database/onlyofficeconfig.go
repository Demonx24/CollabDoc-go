package database

type OnlyOfficeConfig struct {
	Document     Document     `json:"document"`
	DocumentType string       `json:"documentType"`
	EditorConfig EditorConfig `json:"editorConfig"`
}

type Document struct {
	FileType string `json:"fileType"`
	Key      string `json:"key"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

type EditorConfig struct {
	AutoSave      bool           `json:"autosave"`
	Mode          string         `json:"mode"`
	CallbackUrl   string         `json:"callbackUrl"`
	User          Useronlyoffice `json:"user"`
	Permissions   Permissions    `json:"permissions"`
	Lang          string         `json:"lang"`
	Customization Customization  `json:"customization"`
	CoEditing     CoEditing      `json:"coEditing"`
}

type Useronlyoffice struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Permissions struct {
	Edit     bool `json:"edit"`
	Download bool `json:"download"`
	Print    bool `json:"print"`
	Review   bool `json:"review"`
	Comment  bool `json:"comment"`
}

type Customization struct {
	ForceSave         bool `json:"forcesave"`
	Chat              bool `json:"chat"`
	Comments          bool `json:"comments"`
	CompactHeader     bool `json:"compactHeader"`
	Feedback          bool `json:"feedback"`
	Help              bool `json:"help"`
	ToolbarNoTabs     bool `json:"toolbarNoTabs"`
	HideRightMenu     bool `json:"hideRightMenu"`
	HideRuler         bool `json:"hideRuler"`
	HideToolbar       bool `json:"hideToolbar"`
	HideFileMenu      bool `json:"hideFileMenu"`
	HideReviewTab     bool `json:"hideReviewTab"`
	ShowReviewChanges bool `json:"showReviewChanges"`
	HideInsertTab     bool `json:"hideInsertTab"`
	HideHomeTab       bool `json:"hideHomeTab"`
	HideViewTab       bool `json:"hideViewTab"`
}

type CoEditing struct {
	Mode   string `json:"mode"`   // 推荐 "fast"
	Change bool   `json:"change"` // 是否开启实时变更同步
}
