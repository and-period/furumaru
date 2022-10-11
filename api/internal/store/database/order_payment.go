package database

const orderPaymentTable = "order_payments"

var orderPaymentFields = []string{
	"id", "transaction_id", "order_id", "promotion_id", "payment_id", "payment_type",
	"subtotal", "discount", "shipping_charge", "tax", "total",
	"lastname", "firstname", "postal_code", "prefecture", "city",
	"address_line1", "address_line2", "phone_number", "created_at", "updated_at",
}
