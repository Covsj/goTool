package localEmail

type EmailResult struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type wxEmailSetEmailRep struct {
	EmPrefix string `json:"em_prefix"`
}

type wxEmailSetResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		DomainName string `json:"domain_name"`
		EmPrefix   string `json:"em_prefix"`
		Fulldomain string `json:"fulldomain"`
	} `json:"data"`
}

type wxEmailGetMsgResp struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    []WxEmailRespData `json:"data"`
}

type WxEmailRespData struct {
	Id              int    `json:"id"`
	From            string `json:"from"`
	Fromemail       string `json:"fromemail"`
	To              string `json:"to"`
	Subject         string `json:"subject"`
	ToFormat        string `json:"toFormat"`
	MailContentType string `json:"mail_content_type"`
	MailContent     string `json:"mail_content"`
	CreateTime      string `json:"create_time"`
	Hosturl         string `json:"hosturl"`
	UpdateTime      string `json:"update_time"`
	Uid             int    `json:"uid"`
	Paid            int    `json:"paid"`
}

type ImapEmail struct {
	TimeStamp int64  `json:"time_stamp"`
	From      string `json:"from"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}
