package database

const orderFulfillmentTable = "order_fulfillments"

var orderFulfillmentFields = []string{
	"id", "order_id", "shipping_id", "tracking_number",
	"shipping_carrier", "shipping_method", "box_size", "box_count", "weight_total",
	"lastname", "firstname", "postal_code", "prefecture", "city",
	"address_line1", "address_line2", "phone_number", "created_at", "updated_at",
}
