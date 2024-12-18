package model

type MailBody struct {
}

type RequestDetailBody struct {
	REQUEST_NO   string
	FACTORY_NAME string
	REV          int
	START        string
	END          string
	MINUTE_DIFF  float64
	NAME_STATUS  string
}
