package mailer

// Personalization - 送信先情報
type Personalization struct {
	Name          string                 // 送信先名称
	Address       string                 // 送信先メールアドレス
	Type          AddressType            // 宛先タイプ
	Substitutions map[string]interface{} // 動的コンテンツ
}

// Content - メッセージ内容
type Content struct {
	ContentType string
	Value       string
}

// AddressType - メール宛先種別
type AddressType int32

const (
	AddressTypeTo  AddressType = iota // To
	AddressTypeCC                     // CC
	AddressTypeBCC                    // BCC
)

// Substitutions - メール動的コンテンツの生成
func NewSubstitutions(params map[string]string) map[string]interface{} {
	res := make(map[string]interface{}, len(params))
	for k, v := range params {
		res[k] = v
	}
	return res
}
