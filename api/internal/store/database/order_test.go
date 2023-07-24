package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrder(t *testing.T) {
	assert.NotNil(t, NewOrder(nil))
}

func TestOrder_List(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 2)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	err = db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	products := make(entity.Products, 2)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "producer-id", []string{}, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "producer-id", []string{}, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	shipping := testShipping("shipping-id", "coordinator-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)
	schedule := testSchedule("schedule-id", "coordinator-id", "shipping-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)
	address := testAddress("address-id", "user-id", now())
	err = db.DB.Create(&address).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "schedule-id", "", "coordinator-id", now())
	orders[1] = testOrder("order-id02", "user-id", "schedule-id", "", "coordinator-id", now())
	err = db.DB.Create(&orders).Error
	require.NoError(t, err)
	payments := make(entity.Payments, 2)
	payments[0] = testPayment("order-id01", "address-id", "transaction-id01", "payment-id", now())
	orders[0].Payment = *payments[0]
	payments[1] = testPayment("order-id02", "address-id", "transaction-id02", "payment-id", now())
	orders[1].Payment = *payments[1]
	err = db.DB.Create(&payments).Error
	require.NoError(t, err)
	fulfillments := make(entity.Fulfillments, 2)
	fulfillments[0] = testFulfillment("order-id01", "address-id", now())
	orders[0].Fulfillment = *fulfillments[0]
	fulfillments[1] = testFulfillment("order-id02", "address-id", now())
	orders[1].Fulfillment = *fulfillments[1]
	err = db.DB.Create(&fulfillments).Error
	require.NoError(t, err)
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("order-id01", "product-id01", now())
	orders[0].OrderItems = []*entity.OrderItem{items[0]}
	items[1] = testOrderItem("order-id02", "product-id02", now())
	orders[1].OrderItems = []*entity.OrderItem{items[1]}
	err = db.DB.Create(&items).Error
	require.NoError(t, err)
	activities := make(entity.Activities, 2)
	activities[0] = testActivity("activity-id01", "order-id01", "user-id", now())
	orders[0].Activities = []*entity.Activity{activities[0]}
	activities[1] = testActivity("activity-id02", "order-id02", "user-id", now())
	orders[1].Activities = []*entity.Activity{activities[1]}
	err = db.DB.Create(&activities).Error
	require.NoError(t, err)

	type args struct {
		params *ListOrdersParams
	}
	type want struct {
		orders entity.Orders
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListOrdersParams{
					CoordinatorID: "coordinator-id",
					Limit:         2,
					Offset:        1,
				},
			},
			want: want{
				orders: orders[1:],
				hasErr: false,
			},
		},
		{
			name:  "success with sort",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListOrdersParams{
					Orders: []*ListOrdersOrder{
						{Key: entity.OrderOrderByCreatedAt, OrderByASC: true},
						{Key: entity.OrderOrderByUpdatedAt, OrderByASC: false},
					},
				},
			},
			want: want{
				orders: orders,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.orders, actual)
		})
	}
}

func TestOrder_Count(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 2)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	err = db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	products := make(entity.Products, 2)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "producer-id", []string{}, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "producer-id", []string{}, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	shipping := testShipping("shipping-id", "coordinator-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)
	schedule := testSchedule("schedule-id", "coordinator-id", "shipping-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)
	address := testAddress("address-id", "user-id", now())
	err = db.DB.Create(&address).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "schedule-id", "", "coordinator-id", now())
	orders[1] = testOrder("order-id02", "user-id", "schedule-id", "", "coordinator-id", now())
	err = db.DB.Create(&orders).Error
	require.NoError(t, err)
	payments := make(entity.Payments, 2)
	payments[0] = testPayment("order-id01", "address-id", "transaction-id01", "payment-id", now())
	orders[0].Payment = *payments[0]
	payments[1] = testPayment("order-id02", "address-id", "transaction-id02", "payment-id", now())
	orders[1].Payment = *payments[1]
	err = db.DB.Create(&payments).Error
	require.NoError(t, err)
	fulfillments := make(entity.Fulfillments, 2)
	fulfillments[0] = testFulfillment("order-id01", "address-id", now())
	orders[0].Fulfillment = *fulfillments[0]
	fulfillments[1] = testFulfillment("order-id02", "address-id", now())
	orders[1].Fulfillment = *fulfillments[1]
	err = db.DB.Create(&fulfillments).Error
	require.NoError(t, err)
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("order-id01", "product-id01", now())
	orders[0].OrderItems = []*entity.OrderItem{items[0]}
	items[1] = testOrderItem("order-id02", "product-id02", now())
	orders[1].OrderItems = []*entity.OrderItem{items[1]}
	err = db.DB.Create(&items).Error
	require.NoError(t, err)
	activities := make(entity.Activities, 2)
	activities[0] = testActivity("activity-id01", "order-id01", "user-id", now())
	orders[0].Activities = []*entity.Activity{activities[0]}
	activities[1] = testActivity("activity-id02", "order-id02", "user-id", now())
	orders[1].Activities = []*entity.Activity{activities[1]}
	err = db.DB.Create(&activities).Error
	require.NoError(t, err)

	type args struct {
		params *ListOrdersParams
	}
	type want struct {
		total  int64
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListOrdersParams{
					Limit:  2,
					Offset: 1,
				},
			},
			want: want{
				total:  2,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestOrder_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 2)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	err = db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	products := make(entity.Products, 2)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "producer-id", []string{}, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "producer-id", []string{}, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	shipping := testShipping("shipping-id", "coordinator-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)
	schedule := testSchedule("schedule-id", "coordinator-id", "shipping-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)
	address := testAddress("address-id", "user-id", now())
	err = db.DB.Create(&address).Error
	require.NoError(t, err)

	o := testOrder("order-id", "user-id", "schedule-id", "", "coordinator-id", now())
	err = db.DB.Create(&o).Error
	require.NoError(t, err)
	payment := testPayment("order-id", "address-id", "transaction-id", "payment-id", now())
	o.Payment = *payment
	err = db.DB.Create(&payment).Error
	require.NoError(t, err)
	fulfillment := testFulfillment("order-id", "address-id", now())
	o.Fulfillment = *fulfillment
	err = db.DB.Create(&fulfillment).Error
	require.NoError(t, err)
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("order-id", "product-id01", now())
	items[1] = testOrderItem("order-id", "product-id02", now())
	o.OrderItems = items
	err = db.DB.Create(&items).Error
	require.NoError(t, err)
	activities := make(entity.Activities, 2)
	activities[0] = testActivity("activity-id01", "order-id", "user-id", now())
	activities[1] = testActivity("activity-id02", "order-id", "user-id", now())
	o.Activities = activities
	err = db.DB.Create(&activities).Error
	require.NoError(t, err)

	type args struct {
		orderID string
	}
	type want struct {
		order  *entity.Order
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				orderID: "order-id",
			},
			want: want{
				order:  o,
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				orderID: "other-id",
			},
			want: want{
				order:  nil,
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.orderID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.order, actual)
		})
	}
}

func TestOrder_Aggregate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 2)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	err = db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	products := make(entity.Products, 2)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "producer-id", []string{}, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "producer-id", []string{}, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	shipping := testShipping("shipping-id", "coordinator-id", now())
	err = db.DB.Create(&shipping).Error
	require.NoError(t, err)
	schedule := testSchedule("schedule-id", "coordinator-id", "shipping-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)
	address := testAddress("address-id", "user-id", now())
	err = db.DB.Create(&address).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "schedule-id", "", "coordinator-id", now())
	orders[1] = testOrder("order-id02", "user-id", "schedule-id", "", "coordinator-id", now())
	err = db.DB.Create(&orders).Error
	require.NoError(t, err)
	payments := make(entity.Payments, 2)
	payments[0] = testPayment("order-id01", "address-id", "transaction-id01", "payment-id", now())
	orders[0].Payment = *payments[0]
	payments[1] = testPayment("order-id02", "address-id", "transaction-id02", "payment-id", now())
	orders[1].Payment = *payments[1]
	err = db.DB.Create(&payments).Error
	require.NoError(t, err)
	fulfillments := make(entity.Fulfillments, 2)
	fulfillments[0] = testFulfillment("order-id01", "address-id", now())
	orders[0].Fulfillment = *fulfillments[0]
	fulfillments[1] = testFulfillment("order-id02", "address-id", now())
	orders[1].Fulfillment = *fulfillments[1]
	err = db.DB.Create(&fulfillments).Error
	require.NoError(t, err)
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("order-id01", "product-id01", now())
	orders[0].OrderItems = []*entity.OrderItem{items[0]}
	items[1] = testOrderItem("order-id02", "product-id02", now())
	orders[1].OrderItems = []*entity.OrderItem{items[1]}
	err = db.DB.Create(&items).Error
	require.NoError(t, err)
	activities := make(entity.Activities, 2)
	activities[0] = testActivity("activity-id01", "order-id01", "user-id", now())
	orders[0].Activities = []*entity.Activity{activities[0]}
	activities[1] = testActivity("activity-id02", "order-id02", "user-id", now())
	orders[1].Activities = []*entity.Activity{activities[1]}
	err = db.DB.Create(&activities).Error
	require.NoError(t, err)

	type args struct {
		userIDs []string
	}
	type want struct {
		orders entity.AggregatedOrders
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				userIDs: []string{"user-id", "other-id"},
			},
			want: want{
				orders: entity.AggregatedOrders{
					{
						UserID:     "user-id",
						OrderCount: 2,
						Subtotal:   3600,
						Discount:   0,
					},
				},
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.Aggregate(ctx, tt.args.userIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.orders, actual)
		})
	}
}

func testOrder(id, userID, scheduleID, promotionID, coordinatorID string, now time.Time) *entity.Order {
	return &entity.Order{
		ID:                id,
		UserID:            userID,
		ScheduleID:        scheduleID,
		PromotionID:       promotionID,
		CoordinatorID:     coordinatorID,
		PaymentStatus:     entity.PaymentStatusCaptured,
		FulfillmentStatus: entity.FulfillmentStatusFulfilled,
		CancelType:        entity.CancelTypeUnknown,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
