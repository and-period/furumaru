package database

const customerTable = "customers"

var customerFields = []string{
	"user_id", "lastname", "firstname", "lastname_kana", "firstname_kana",
	"postal_code", "prefecture", "city", "address_line1", "address_line2",
	"created_at", "updated_at",
}
