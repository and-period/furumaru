package database

const orderItemTable = "order_items"

var orderItemFields = []string{
	"id", "order_id", "product_id",
	"price", "quantity", "weight", "weight_unit",
	"created_at", "updated_at",
}
