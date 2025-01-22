package komoju

// PaymentStatus 決済ステータス
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"    // 保留
	PaymentStatusAuthorized PaymentStatus = "authorized" // 承認済み
	PaymentStatusCaptured   PaymentStatus = "captured"   // 入金済み
	PaymentStatusRefunded   PaymentStatus = "refunded"   // 返金済み
	PaymentStatusCancelled  PaymentStatus = "cancelled"  // キャンセル済み
	PaymentStatusExpired    PaymentStatus = "expired"    // 期限切れ
	PaymentStatusFailed     PaymentStatus = "failed"     // 失敗
)

// SessionStatus 決済トランザクションステータス
type SessionStatus string

const (
	SessionStatusPending   SessionStatus = "pending"   // 保留
	SessionStatusCompleted SessionStatus = "completed" // 完了
	SessionStatusCancelled SessionStatus = "cancelled" // キャンセル済み
)

// PaymentType 決済種別
type PaymentType string

const (
	PaymentTypeCreditCard   PaymentType = "credit_card"   // クレジットカード
	PaymentTypeBankTransfer PaymentType = "bank_transfer" // 銀行振込
	PaymentTypeKonbini      PaymentType = "konbini"       // コンビニ
	PaymentTypeLinePay      PaymentType = "linepay"       // LINE Pay
	PaymentTypeMerpay       PaymentType = "merpay"        // メルペイ
	PaymentTypePayPay       PaymentType = "paypay"        // PayPay
	PaymentTypeRakutenPay   PaymentType = "rakutenpay"    // 楽天ペイ
	PaymentTypeAUPay        PaymentType = "aupay"         // au PAY
	PaymentTypePaidy        PaymentType = "paidy"         // Paidy
	PaymentTypePayEasy      PaymentType = "pay_easy"      // Pay-easy
)

// KonbiniType コンビニ店舗種別
type KonbiniType string

const (
	KonbiniTypeDailyYamazaki KonbiniType = "daily-yamazaki" // デイリーヤマザキ
	KonbiniTypeFamilyMart    KonbiniType = "family-mart"    // ファミリーマート
	KonbiniTypeLawson        KonbiniType = "lawson"         // ローソン
	KonbiniTypeMinistop      KonbiniType = "ministop"       // ミニストップ
	KonbiniTypeSeicomart     KonbiniType = "seicomart"      // セイコーマート
	KonbiniTypeSevenEleven   KonbiniType = "seven-eleven"   // セブンイレブン
)

// BankAccountType 預金口座種別
type BankAccountType string

const (
	BankAccountTypeNormal   BankAccountType = "normal"   // 普通口座
	BankAccountTypeChecking BankAccountType = "checking" // 当座口座
	BankAccountTypeSaving   BankAccountType = "saving"   // 貯蓄口座
)

// CaptureMode 売上処理方法
type CaptureMode string

const (
	CaptureModeManual CaptureMode = "manual" // 仮売上・実売上
	CaptureModeAuto   CaptureMode = "auto"   // 即時売上
)
