package request

import "time"

type PaymentRequest struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Resource  string          `json:"resource"`
	Payload   *PaymentPayload `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
	Reason    string          `json:"reason,omitempty"`
}

type PaymentPayload struct {
	ID                  string           `json:"id"`
	Resource            string           `json:"resource"`
	Status              string           `json:"status"`
	Amount              int64            `json:"amount"`
	Tax                 int64            `json:"tax"`
	Customer            string           `json:"customer,omitempty"`
	PaymentDeadline     time.Time        `json:"payment_deadline,omitempty"`
	PaymentDetails      *PaymentDetails  `json:"payment_details"`
	PaymentMethodFee    int64            `json:"payment_method_fee"`
	Total               int64            `json:"total"`
	Currency            string           `json:"currency"`
	Description         string           `json:"description,omitempty"`
	CapturedAt          time.Time        `json:"captured_at,omitempty"`
	ExternalOrderNumber string           `json:"external_order_num,omitempty"`
	Metadata            map[string]any   `json:"metadata"`
	CreatedAt           time.Time        `json:"created_at"`
	AmountRefunded      int64            `json:"amount_refunded"`
	Locale              string           `json:"locale"`
	Session             string           `json:"session,omitempty"`
	CustomerFamilyName  string           `json:"customer_family_name,omitempty"`
	CustomerGivenName   string           `json:"customer_given_name,omitempty"`
	MCC                 string           `json:"mcc,omitempty"`
	StatementDescriptor string           `json:"statement_descriptor,omitempty"`
	Refunds             []*Refund        `json:"refunds"`
	RefundRequests      []*RefundRequest `json:"refund_requests"`
}

type PaymentDetails struct {
	Type              string `json:"type"`
	Email             string `json:"email,omitempty"`
	OrderID           string `json:"order_id,omitempty"`
	Brand             string `json:"brand,omitempty"`
	LastFourDigits    string `json:"last_four_digits,omitempty"`
	Month             int64  `json:"month,omitempty"`
	Year              int64  `json:"year,omitempty"`
	BankName          string `json:"bank_name,omitempty"`
	AccountBranchName string `json:"account_branch_name,omitempty"`
	AccountNumber     string `json:"account_number,omitempty"`
	AccountType       string `json:"account_type,omitempty"`
	AccountName       string `json:"account_name,omitempty"`
	InstructionsURL   string `json:"instructions_url,omitempty"`
	RedirectURL       string `json:"redirect_url,omitempty"`
	ExternalPaymentID string `json:"external_payment_id,omitempty"`
}

type Refund struct {
	ID          string    `json:"id"`
	Resource    string    `json:"resource"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	Payment     string    `json:"payment"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Chargeback  bool      `json:"chargeback"`
}

type RefundRequest struct {
	ID            string `json:"id"`
	Payment       string `json:"payment"`
	CustomerName  string `json:"customer_name"`
	BankName      string `json:"bank_name"`
	BankCode      string `json:"bank_code"`
	BranchName    string `json:"branch_name"`
	BranchNumber  string `json:"branch_number"`
	AccountNumber string `json:"account_number"`
	Status        string `json:"status"`
	Description   string `json:"description"`
}
