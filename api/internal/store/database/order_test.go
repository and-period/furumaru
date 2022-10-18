package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
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

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx,
		orderItemTable, orderActivityTable, orderPaymentTable, orderFulfillmentTable,
		orderTable, scheduleTable, shippingTable, productTable, productTypeTable, categoryTable,
	)
	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = m.db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 2)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	err = m.db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	products := make(entity.Products, 2)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "producer-id", now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "producer-id", now())
	err = m.db.DB.Create(&products).Error
	require.NoError(t, err)
	shipping := testShipping("shipping-id", now())
	err = m.db.DB.Create(&shipping).Error
	require.NoError(t, err)
	schedule := testSchedule("schedule-id", now())
	err = m.db.DB.Create(&schedule).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "schedule-id", "coordinator-id", now())
	orders[1] = testOrder("order-id02", "user-id", "schedule-id", "coordinator-id", now())
	err = m.db.DB.Create(&orders).Error
	require.NoError(t, err)
	payments := make(entity.OrderPayments, 2)
	payments[0] = testOrderPayment("payment-id01", "transaction-id01", "order-id01", "", "payment-id", now())
	orders[0].OrderPayment = *payments[0]
	payments[1] = testOrderPayment("payment-id02", "transaction-id02", "order-id02", "", "payment-id", now())
	orders[1].OrderPayment = *payments[1]
	err = m.db.DB.Create(&payments).Error
	require.NoError(t, err)
	fulfillments := make(entity.OrderFulfillments, 2)
	fulfillments[0] = testOrderFulfillment("fulfillment-id01", "order-id01", "shipping-id", now())
	orders[0].OrderFulfillment = *fulfillments[0]
	fulfillments[1] = testOrderFulfillment("fulfillment-id02", "order-id02", "shipping-id", now())
	orders[1].OrderFulfillment = *fulfillments[1]
	err = m.db.DB.Create(&fulfillments).Error
	require.NoError(t, err)
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("item-id01", "order-id01", "product-id01", now())
	orders[0].OrderItems = []*entity.OrderItem{items[0]}
	items[1] = testOrderItem("item-id02", "order-id02", "product-id02", now())
	orders[1].OrderItems = []*entity.OrderItem{items[1]}
	err = m.db.DB.Create(&items).Error
	require.NoError(t, err)
	activities := make(entity.OrderActivities, 2)
	activities[0] = testOrderActivity("activity-id01", "order-id01", "user-id", now())
	orders[0].OrderActivities = []*entity.OrderActivity{activities[0]}
	activities[1] = testOrderActivity("activity-id02", "order-id02", "user-id", now())
	orders[1].OrderActivities = []*entity.OrderActivity{activities[1]}
	err = m.db.DB.Create(&activities).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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

			tt.setup(ctx, t, m)

			db := &order{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			fillIgnoreOrdersField(actual, now())
			assert.Equal(t, tt.want.orders, actual)
		})
	}
}

func TestOrder_Count(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx,
		orderItemTable, orderActivityTable, orderPaymentTable, orderFulfillmentTable,
		orderTable, shippingTable, productTable, productTypeTable, categoryTable,
	)
	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = m.db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 2)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	err = m.db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	products := make(entity.Products, 2)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "producer-id", now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "producer-id", now())
	err = m.db.DB.Create(&products).Error
	require.NoError(t, err)
	shipping := testShipping("shipping-id", now())
	err = m.db.DB.Create(&shipping).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "schedule-id", "coordinator-id", now())
	orders[1] = testOrder("order-id02", "user-id", "schedule-id", "coordinator-id", now())
	err = m.db.DB.Create(&orders).Error
	require.NoError(t, err)
	payments := make(entity.OrderPayments, 2)
	payments[0] = testOrderPayment("payment-id01", "transaction-id01", "order-id01", "", "payment-id", now())
	orders[0].OrderPayment = *payments[0]
	payments[1] = testOrderPayment("payment-id02", "transaction-id02", "order-id02", "", "payment-id", now())
	orders[1].OrderPayment = *payments[1]
	err = m.db.DB.Create(&payments).Error
	require.NoError(t, err)
	fulfillments := make(entity.OrderFulfillments, 2)
	fulfillments[0] = testOrderFulfillment("fulfillment-id01", "order-id01", "shipping-id", now())
	orders[0].OrderFulfillment = *fulfillments[0]
	fulfillments[1] = testOrderFulfillment("fulfillment-id02", "order-id02", "shipping-id", now())
	orders[1].OrderFulfillment = *fulfillments[1]
	err = m.db.DB.Create(&fulfillments).Error
	require.NoError(t, err)
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("item-id01", "order-id01", "product-id01", now())
	orders[0].OrderItems = []*entity.OrderItem{items[0]}
	items[1] = testOrderItem("item-id02", "order-id02", "product-id02", now())
	orders[1].OrderItems = []*entity.OrderItem{items[1]}
	err = m.db.DB.Create(&items).Error
	require.NoError(t, err)
	activities := make(entity.OrderActivities, 2)
	activities[0] = testOrderActivity("activity-id01", "order-id01", "user-id", now())
	orders[0].OrderActivities = []*entity.OrderActivity{activities[0]}
	activities[1] = testOrderActivity("activity-id02", "order-id02", "user-id", now())
	orders[1].OrderActivities = []*entity.OrderActivity{activities[1]}
	err = m.db.DB.Create(&activities).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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

			tt.setup(ctx, t, m)

			db := &order{db: m.db, now: now}
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

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx,
		orderItemTable, orderActivityTable, orderPaymentTable, orderFulfillmentTable,
		orderTable, shippingTable, productTable, productTypeTable, categoryTable,
	)
	categories := make(entity.Categories, 2)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "果物", now())
	err = m.db.DB.Create(&categories).Error
	require.NoError(t, err)
	productTypes := make(entity.ProductTypes, 2)
	productTypes[0] = testProductType("type-id01", "category-id01", "野菜", now())
	productTypes[1] = testProductType("type-id02", "category-id02", "果物", now())
	err = m.db.DB.Create(&productTypes).Error
	require.NoError(t, err)
	products := make(entity.Products, 2)
	products[0] = testProduct("product-id01", "type-id01", "category-id01", "producer-id", now())
	products[1] = testProduct("product-id02", "type-id02", "category-id02", "producer-id", now())
	err = m.db.DB.Create(&products).Error
	require.NoError(t, err)
	shipping := testShipping("shipping-id", now())
	err = m.db.DB.Create(&shipping).Error
	require.NoError(t, err)

	o := testOrder("order-id", "user-id", "schedule-id", "coordinator-id", now())
	err = m.db.DB.Create(&o).Error
	require.NoError(t, err)
	payment := testOrderPayment("payment-id", "transaction-id", "order-id", "", "payment-id", now())
	o.OrderPayment = *payment
	err = m.db.DB.Create(&payment).Error
	require.NoError(t, err)
	fulfillment := testOrderFulfillment("fulfillment-id", "order-id", "shipping-id", now())
	o.OrderFulfillment = *fulfillment
	err = m.db.DB.Create(&fulfillment).Error
	require.NoError(t, err)
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("item-id01", "order-id", "product-id01", now())
	items[1] = testOrderItem("item-id02", "order-id", "product-id02", now())
	o.OrderItems = items
	err = m.db.DB.Create(&items).Error
	require.NoError(t, err)
	activities := make(entity.OrderActivities, 2)
	activities[0] = testOrderActivity("activity-id01", "order-id", "user-id", now())
	activities[1] = testOrderActivity("activity-id02", "order-id", "user-id", now())
	o.OrderActivities = activities
	err = m.db.DB.Create(&activities).Error
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
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
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

			tt.setup(ctx, t, m)

			db := &order{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.orderID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			fillIgnoreOrderField(actual, now())
			assert.Equal(t, tt.want.order, actual)
		})
	}
}

func testOrder(id, userID, scheduleID, coordinatorID string, now time.Time) *entity.Order {
	return &entity.Order{
		ID:                id,
		UserID:            userID,
		ScheduleID:        scheduleID,
		CoordinatorID:     coordinatorID,
		PaymentStatus:     entity.PaymentStatusCaptured,
		FulfillmentStatus: entity.FulfillmentStatusFulfilled,
		CancelType:        entity.CancelTypeUnknown,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func fillIgnoreOrderField(o *entity.Order, now time.Time) {
	if o == nil {
		return
	}
	o.CreatedAt = now
	o.UpdatedAt = now
	o.OrderPayment.CreatedAt = now
	o.OrderPayment.UpdatedAt = now
	o.OrderFulfillment.CreatedAt = now
	o.OrderFulfillment.UpdatedAt = now
	for i := range o.OrderActivities {
		o.OrderActivities[i].CreatedAt = now
		o.OrderActivities[i].UpdatedAt = now
	}
	for i := range o.OrderItems {
		o.OrderItems[i].CreatedAt = now
		o.OrderItems[i].UpdatedAt = now
	}
}

func fillIgnoreOrdersField(os entity.Orders, now time.Time) {
	for i := range os {
		fillIgnoreOrderField(os[i], now)
	}
}
