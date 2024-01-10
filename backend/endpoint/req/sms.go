package req

type SmsVerify struct {
	Type string `json:"type"`
	Code string `json:"code"`
}
