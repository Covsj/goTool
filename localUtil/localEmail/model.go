package localEmail

type LinShiEmailResp struct {
	List []struct {
		NameTo      string `json:"name_to"`
		NameFrom    string `json:"name_from"`
		Eid         string `json:"eid"`
		ESubject    string `json:"e_subject"`
		EDate       int64  `json:"e_date"`
		AddressFrom string `json:"address_from"`
	} `json:"list"`
	Status string `json:"status"`
}

type LinShiEmail struct {
	Data struct {
		To      string `json:"to"`
		Seqno   int    `json:"seqno"`
		Subject string `json:"subject"`
		From    struct {
			Address string `json:"address"`
			Name    string `json:"name"`
		} `json:"from"`
		Date      int64       `json:"date"`
		Html      interface{} `json:"html"`
		MessageId string      `json:"messageId"`
		Name      string      `json:"name"`
		Eid       string      `json:"eid"`
	} `json:"data"`
	Status string `json:"status"`
}

type RespWx struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		DomainName string `json:"domain_name"`
		EmPrefix   string `json:"em_prefix"`
		Fulldomain string `json:"fulldomain"`
	} `json:"data"`
}

type WxShiYinNuoCheEmail struct {
	Message string `json:"message"`
	Data    []struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
	} `json:"data"`
}
