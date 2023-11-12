package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrder(t *testing.T) {
	assert.NotNil(t, newOrder(nil))
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
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "coordinator-id", "producer-id", []string{}, 1, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err = db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "", "coordinator-id", now())
	orders[1] = testOrder("order-id02", "user-id", "", "coordinator-id", now())
	err = db.DB.Create(&orders).Error
	require.NoError(t, err)
	payments := make(entity.OrderPayments, 2)
	payments[0] = testOrderPayment("order-id01", 1, "transaction-id01", "payment-id", now())
	orders[0].OrderPayment = *payments[0]
	payments[1] = testOrderPayment("order-id02", 1, "transaction-id02", "payment-id", now())
	orders[1].OrderPayment = *payments[1]
	err = db.DB.Create(&payments).Error
	require.NoError(t, err)
	fulfillments := make(entity.OrderFulfillments, 2)
	fulfillments[0] = testOrderFulfillment("fulfillment-id01", "order-id01", 1, 1, now())
	orders[0].OrderFulfillments = entity.OrderFulfillments{fulfillments[0]}
	fulfillments[1] = testOrderFulfillment("fulfillment-id02", "order-id02", 1, 2, now())
	orders[1].OrderFulfillments = entity.OrderFulfillments{fulfillments[1]}
	err = db.DB.Create(&fulfillments).Error
	require.NoError(t, err)
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("fulfillment-id01", 1, "order-id01", now())
	orders[0].OrderItems = []*entity.OrderItem{items[0]}
	items[1] = testOrderItem("fulfillment-id02", 2, "order-id02", now())
	orders[1].OrderItems = []*entity.OrderItem{items[1]}
	err = db.DB.Create(&items).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListOrdersParams
	}
	type want struct {
		orders entity.Orders
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListOrdersParams{
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
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "coordinator-id", "producer-id", []string{}, 1, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err = db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "", "coordinator-id", now())
	orders[1] = testOrder("order-id02", "user-id", "", "coordinator-id", now())
	err = db.DB.Create(&orders).Error
	require.NoError(t, err)
	payments := make(entity.OrderPayments, 2)
	payments[0] = testOrderPayment("order-id01", 1, "transaction-id01", "payment-id", now())
	orders[0].OrderPayment = *payments[0]
	payments[1] = testOrderPayment("order-id02", 1, "transaction-id02", "payment-id", now())
	orders[1].OrderPayment = *payments[1]
	err = db.DB.Create(&payments).Error
	require.NoError(t, err)
	fulfillments := make(entity.OrderFulfillments, 2)
	fulfillments[0] = testOrderFulfillment("fulfillment-id01", "order-id01", 1, 1, now())
	orders[0].OrderFulfillments = entity.OrderFulfillments{fulfillments[0]}
	fulfillments[1] = testOrderFulfillment("fulfillment-id02", "order-id02", 1, 2, now())
	orders[1].OrderFulfillments = entity.OrderFulfillments{fulfillments[1]}
	err = db.DB.Create(&fulfillments).Error
	require.NoError(t, err)
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("fulfillment-id01", 1, "order-id01", now())
	orders[0].OrderItems = []*entity.OrderItem{items[0]}
	items[1] = testOrderItem("fulfillment-id02", 2, "order-id02", now())
	orders[1].OrderItems = []*entity.OrderItem{items[1]}
	err = db.DB.Create(&items).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListOrdersParams
	}
	type want struct {
		total  int64
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListOrdersParams{
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
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "coordinator-id", "producer-id", []string{}, 1, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err = db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	o := testOrder("order-id", "user-id", "", "coordinator-id", now())
	err = db.DB.Create(&o).Error
	require.NoError(t, err)
	payment := testOrderPayment("order-id", 1, "transaction-id", "payment-id", now())
	o.OrderPayment = *payment
	err = db.DB.Create(&payment).Error
	require.NoError(t, err)
	fulfillments := make(entity.OrderFulfillments, 1)
	fulfillments[0] = testOrderFulfillment("fulfillment-id", "order-id", 1, 1, now())
	o.OrderFulfillments = fulfillments
	err = db.DB.Create(&fulfillments).Error
	require.NoError(t, err)
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("fulfillment-id", 1, "order-id", now())
	items[1] = testOrderItem("fulfillment-id", 2, "order-id", now())
	o.OrderItems = items
	err = db.DB.Create(&items).Error
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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

func TestOrder_Create(t *testing.T) {
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
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "coordinator-id", "producer-id", []string{}, 1, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err = db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	fulfillments := make(entity.OrderFulfillments, 1)
	fulfillments[0] = testOrderFulfillment("fulfillment-id", "order-id", 1, 1, now())
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("fulfillment-id", 1, "order-id", now())
	items[1] = testOrderItem("fulfillment-id", 2, "order-id", now())

	o := testOrder("order-id", "user-id", "", "coordinator-id", now())
	o.OrderPayment = *testOrderPayment("order-id", 1, "transaction-id", "payment-id", now())
	o.OrderFulfillments = fulfillments
	o.OrderItems = items

	type args struct {
		order *entity.Order
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				order: o,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				err := db.DB.Create(&o).Error
				require.NoError(t, err)
			},
			args: args{
				order: o,
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, orderItemTable, orderFulfillmentTable, orderPaymentTable, orderTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			err = db.Create(ctx, tt.args.order)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
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
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "coordinator-id", "producer-id", []string{}, 1, now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Create(&products).Error
	require.NoError(t, err)
	for i := range products {
		err = db.DB.Create(&products[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "", "coordinator-id", now())
	orders[1] = testOrder("order-id02", "user-id", "", "coordinator-id", now())
	err = db.DB.Create(&orders).Error
	require.NoError(t, err)
	payments := make(entity.OrderPayments, 2)
	payments[0] = testOrderPayment("order-id01", 1, "transaction-id01", "payment-id", now())
	orders[0].OrderPayment = *payments[0]
	payments[1] = testOrderPayment("order-id02", 1, "transaction-id02", "payment-id", now())
	orders[1].OrderPayment = *payments[1]
	err = db.DB.Create(&payments).Error
	require.NoError(t, err)
	fulfillments := make(entity.OrderFulfillments, 2)
	fulfillments[0] = testOrderFulfillment("fulfillment-id01", "order-id01", 1, 1, now())
	orders[0].OrderFulfillments = entity.OrderFulfillments{fulfillments[0]}
	fulfillments[1] = testOrderFulfillment("fulfillment-id02", "order-id02", 1, 2, now())
	orders[1].OrderFulfillments = entity.OrderFulfillments{fulfillments[1]}
	err = db.DB.Create(&fulfillments).Error
	require.NoError(t, err)
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("fulfillment-id01", 1, "order-id01", now())
	orders[0].OrderItems = []*entity.OrderItem{items[0]}
	items[1] = testOrderItem("fulfillment-id02", 2, "order-id02", now())
	orders[1].OrderItems = []*entity.OrderItem{items[1]}
	err = db.DB.Create(&items).Error
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
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
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

func testOrder(id, userID, promotionID, coordinatorID string, now time.Time) *entity.Order {
	return &entity.Order{
		ID:            id,
		UserID:        userID,
		PromotionID:   promotionID,
		CoordinatorID: coordinatorID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func testOrderPayment(orderID string, addressID int64, transactionID, paymentMethodID string, now time.Time) *entity.OrderPayment {
	return &entity.OrderPayment{
		OrderID:           orderID,
		AddressRevisionID: addressID,
		Status:            entity.PaymentStatusCaptured,
		TransactionID:     transactionID,
		MethodType:        entity.PaymentMethodTypeCreditCard,
		Subtotal:          1800,
		Discount:          0,
		ShippingFee:       500,
		Tax:               230,
		Total:             2530,
		RefundTotal:       0,
		RefundReason:      "",
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func testOrderFulfillment(fulfillmentID, orderID string, addressID, boxNumber int64, now time.Time) *entity.OrderFulfillment {
	return &entity.OrderFulfillment{
		ID:                fulfillmentID,
		OrderID:           orderID,
		AddressRevisionID: addressID,
		Status:            entity.FulfillmentStatusFulfilled,
		TrackingNumber:    "",
		ShippingCarrier:   entity.ShippingCarrierUnknown,
		ShippingType:      entity.ShippingTypeNormal,
		BoxNumber:         boxNumber,
		BoxSize:           entity.ShippingSize60,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func testOrderItem(fulfillmentID string, productID int64, orderID string, now time.Time) *entity.OrderItem {
	return &entity.OrderItem{
		FulfillmentID:     fulfillmentID,
		ProductRevisionID: productID,
		OrderID:           orderID,
		Quantity:          1,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}
