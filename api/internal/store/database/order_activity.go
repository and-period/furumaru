package database

const orderActivityTable = "order_activities"

var orderActivityFields = []string{
	"id", "order_id", "user_id", "event_type", "detail",
	"metadata", "created_at", "updated_at",
}
