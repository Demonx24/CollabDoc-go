package request

type Websocketlogin struct {
	DocId    string `json:"docId"form:"docId"`
	Useruuid string `json:"useruuid"form:"useruuid"`
}
