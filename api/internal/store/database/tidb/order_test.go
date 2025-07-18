package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrder(t *testing.T) {
	assert.NotNil(t, NewOrder(nil))
}

func TestOrder_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now().Add(-time.Hour))
	orders[1] = testOrder("order-id02", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 2, now())
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
					ShopID: "shop-id",
					Limit:  2,
					Offset: 0,
				},
			},
			want: want{
				orders: orders,
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.orders, actual)
		})
	}
}

func TestOrder_ListUserIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}

	experienceTypes := make(entity.ExperienceTypes, 1)
	experienceTypes[0] = testExperienceType("experience-type-id", "体験", now())
	err = db.DB.Create(&experienceTypes).Error
	require.NoError(t, err)
	experiences := make(internalExperiences, 1)
	experiences[0] = testExperience("experience-id", "experience-type-id", "shop-id", "coordinator-id", "producer-id", 1, now())
	err = db.DB.Table(experienceTable).Create(&experiences).Error
	require.NoError(t, err)
	for i := range experiences {
		err = db.DB.Create(&experiences[i].ExperienceRevision).Error
		require.NoError(t, err)
	}

	orders := make(entity.Orders, 3)
	orders[0] = testOrder("order-id01", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now())
	orders[1] = testOrder("order-id02", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 2, now())
	orders[2] = testOrder("order-id03", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeExperience, 3, now())
	err = db.DB.Create(&orders).Error
	require.NoError(t, err)

	payments := make(entity.OrderPayments, 3)
	payments[0] = testOrderPayment("order-id01", 1, "transaction-id01", "payment-id", now())
	orders[0].OrderPayment = *payments[0]
	payments[1] = testOrderPayment("order-id02", 1, "transaction-id02", "payment-id", now())
	orders[1].OrderPayment = *payments[1]
	payments[2] = testOrderPayment("order-id03", 1, "transaction-id03", "payment-id", now())
	orders[2].OrderPayment = *payments[2]
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

	einternal := make(internalOrderExperiences, 1)
	einternal[0] = testOrderExperience("order-id03", 1, now())
	oexperiences, err := einternal.entities()
	require.NoError(t, err)
	orders[2].OrderExperience = *oexperiences[0]
	err = db.DB.Table(orderExperienceTable).Create(&einternal).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListOrdersParams
	}
	type want struct {
		userIDs []string
		total   int64
		hasErr  bool
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
					ShopID: "shop-id",
					Limit:  10,
					Offset: 0,
				},
			},
			want: want{
				userIDs: []string{"user-id"},
				total:   1,
				hasErr:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, total, err := db.ListUserIDs(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.userIDs, actual)
			assert.Equal(t, tt.want.total, total)
		})
	}
}

func TestOrder_Count(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now())
	orders[1] = testOrder("order-id02", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 2, now())
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestOrder_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	o := testOrder("order-id", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now())
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.orderID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.order, actual)
		})
	}
}

func TestOrder_GetByTransactionID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	o := testOrder("order-id", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now())
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
		userID        string
		transactionID string
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
				userID:        "user-id",
				transactionID: "transaction-id",
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
				userID:        "",
				transactionID: "",
			},
			want: want{
				order:  nil,
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.GetByTransactionID(ctx, tt.args.userID, tt.args.transactionID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.order, actual)
		})
	}
}

func TestOrder_GetByTransactionIDWithSessionID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	o := testOrder("order-id", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now())
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
		sessionID     string
		transactionID string
	}
	type want struct {
		order *entity.Order
		err   error
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
				sessionID:     "session-id",
				transactionID: "transaction-id",
			},
			want: want{
				order: o,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				sessionID:     "",
				transactionID: "",
			},
			want: want{
				order: nil,
				err:   database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.GetByTransactionIDWithSessionID(ctx, tt.args.sessionID, tt.args.transactionID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.order, actual)
		})
	}
}

func TestOrder_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	experienceType := testExperienceType("experience-type-id", "体験", now())
	err = db.DB.Create(&experienceType).Error
	require.NoError(t, err)
	experience := testExperience("experience-id", "experience-type-id", "shop-id", "coordinator-id", "producer-id", 1, now())
	err = db.DB.Table(experienceTable).Create(&experience).Error
	require.NoError(t, err)
	err = db.DB.Create(&experience.ExperienceRevision).Error
	require.NoError(t, err)

	fulfillments := make(entity.OrderFulfillments, 1)
	fulfillments[0] = testOrderFulfillment("fulfillment-id", "product-order-id", 1, 1, now())
	items := make(entity.OrderItems, 2)
	items[0] = testOrderItem("fulfillment-id", 1, "product-order-id", now())
	items[1] = testOrderItem("fulfillment-id", 2, "product-order-id", now())

	porder := testOrder("product-order-id", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now())
	porder.Type = entity.OrderTypeProduct
	porder.OrderPayment = *testOrderPayment("product-order-id", 1, "transaction-id", "payment-id", now())
	porder.OrderFulfillments = fulfillments
	porder.OrderItems = items

	eorder := testOrder("experience-order-id", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeExperience, 2, now())
	eorder.Type = entity.OrderTypeExperience
	eorder.OrderPayment = *testOrderPayment("experience-order-id", 1, "transaction-id", "payment-id", now())
	einternal := testOrderExperience("experience-order-id", 1, now())
	oexperience, err := einternal.entity()
	require.NoError(t, err)
	eorder.OrderExperience = *oexperience

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
			name:  "success product order",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				order: porder,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "success experience order",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				order: eorder,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				err := db.DB.Create(&porder).Error
				require.NoError(t, err)
			},
			args: args{
				order: porder,
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, orderItemTable, orderFulfillmentTable, orderPaymentTable, orderExperienceTable, orderTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			err = db.Create(ctx, tt.args.order)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestOrder_UpdateAuthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	create := func(t *testing.T, orderID string, status entity.OrderStatus, now time.Time) {
		order := testOrder(orderID, "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now)
		order.Status = status
		err := db.DB.Create(&order).Error
		require.NoError(t, err)

		payment := testOrderPayment(orderID, 1, "transaction-id", "payment-id", now)
		err = db.DB.Create(&payment).Error
		require.NoError(t, err)

		fulfillments := make(entity.OrderFulfillments, 1)
		fulfillments[0] = testOrderFulfillment("fulfillment-id", orderID, 1, 1, now)
		err = db.DB.Create(&fulfillments).Error
		require.NoError(t, err)

		items := make(entity.OrderItems, 2)
		items[0] = testOrderItem("fulfillment-id", 1, orderID, now)
		items[1] = testOrderItem("fulfillment-id", 2, orderID, now)
		err = db.DB.Create(&items).Error
		require.NoError(t, err)
	}

	type args struct {
		orderID string
		params  *database.UpdateOrderAuthorizedParams
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.OrderStatusUnpaid, now().AddDate(0, 0, -1))
			},
			args: args{
				orderID: "order-id",
				params: &database.UpdateOrderAuthorizedParams{
					PaymentID: "payment-id",
					IssuedAt:  now(),
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				orderID: "",
				params: &database.UpdateOrderAuthorizedParams{
					PaymentID: "payment-id",
					IssuedAt:  now(),
				},
			},
			want: want{
				err: database.ErrNotFound,
			},
		},
		{
			name: "not latest data",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.OrderStatusUnpaid, now().AddDate(0, 0, 1))
			},
			args: args{
				orderID: "order-id",
				params: &database.UpdateOrderAuthorizedParams{
					PaymentID: "",
					IssuedAt:  now(),
				},
			},
			want: want{
				err: database.ErrFailedPrecondition,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, orderItemTable, orderFulfillmentTable, orderPaymentTable, orderExperienceTable, orderTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			err = db.UpdateAuthorized(ctx, tt.args.orderID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestOrder_UpdateCaptured(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	create := func(t *testing.T, orderID string, status entity.OrderStatus, now time.Time) {
		order := testOrder(orderID, "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now)
		order.Status = status
		err := db.DB.Create(&order).Error
		require.NoError(t, err)

		payment := testOrderPayment(orderID, 1, "transaction-id", "payment-id", now)
		err = db.DB.Create(&payment).Error
		require.NoError(t, err)

		fulfillments := make(entity.OrderFulfillments, 1)
		fulfillments[0] = testOrderFulfillment("fulfillment-id", orderID, 1, 1, now)
		err = db.DB.Create(&fulfillments).Error
		require.NoError(t, err)

		items := make(entity.OrderItems, 2)
		items[0] = testOrderItem("fulfillment-id", 1, orderID, now)
		items[1] = testOrderItem("fulfillment-id", 2, orderID, now)
		err = db.DB.Create(&items).Error
		require.NoError(t, err)
	}

	type args struct {
		orderID string
		params  *database.UpdateOrderCapturedParams
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.OrderStatusUnpaid, now().AddDate(0, 0, -1))
			},
			args: args{
				orderID: "order-id",
				params: &database.UpdateOrderCapturedParams{
					PaymentID: "payment-id",
					IssuedAt:  now(),
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				orderID: "",
				params: &database.UpdateOrderCapturedParams{
					PaymentID: "payment-id",
					IssuedAt:  now(),
				},
			},
			want: want{
				err: database.ErrNotFound,
			},
		},
		{
			name: "not latest data",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.OrderStatusUnpaid, now().AddDate(0, 0, 1))
			},
			args: args{
				orderID: "order-id",
				params: &database.UpdateOrderCapturedParams{
					PaymentID: "",
					IssuedAt:  now(),
				},
			},
			want: want{
				err: database.ErrFailedPrecondition,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, orderItemTable, orderFulfillmentTable, orderPaymentTable, orderExperienceTable, orderTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			err = db.UpdateCaptured(ctx, tt.args.orderID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestOrder_UpdateFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	create := func(t *testing.T, orderID string, status entity.OrderStatus, now time.Time) {
		order := testOrder(orderID, "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now)
		order.Status = status
		err := db.DB.Create(&order).Error
		require.NoError(t, err)

		payment := testOrderPayment(orderID, 1, "transaction-id", "payment-id", now)
		err = db.DB.Create(&payment).Error
		require.NoError(t, err)

		fulfillments := make(entity.OrderFulfillments, 1)
		fulfillments[0] = testOrderFulfillment("fulfillment-id", orderID, 1, 1, now)
		err = db.DB.Create(&fulfillments).Error
		require.NoError(t, err)

		items := make(entity.OrderItems, 2)
		items[0] = testOrderItem("fulfillment-id", 1, orderID, now)
		items[1] = testOrderItem("fulfillment-id", 2, orderID, now)
		err = db.DB.Create(&items).Error
		require.NoError(t, err)
	}

	type args struct {
		orderID string
		params  *database.UpdateOrderFailedParams
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.OrderStatusUnpaid, now().AddDate(0, 0, -1))
			},
			args: args{
				orderID: "order-id",
				params: &database.UpdateOrderFailedParams{
					Status:    entity.PaymentStatusFailed,
					PaymentID: "",
					IssuedAt:  now(),
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				orderID: "",
				params: &database.UpdateOrderFailedParams{
					Status:    entity.PaymentStatusFailed,
					PaymentID: "payment-id",
					IssuedAt:  now(),
				},
			},
			want: want{
				err: database.ErrNotFound,
			},
		},
		{
			name: "not latest data",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.OrderStatusUnpaid, now().AddDate(0, 0, 1))
			},
			args: args{
				orderID: "order-id",
				params: &database.UpdateOrderFailedParams{
					Status:    entity.PaymentStatusFailed,
					PaymentID: "",
					IssuedAt:  now(),
				},
			},
			want: want{
				err: database.ErrFailedPrecondition,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, orderItemTable, orderFulfillmentTable, orderPaymentTable, orderExperienceTable, orderTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			err = db.UpdateFailed(ctx, tt.args.orderID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestOrder_UpdateRefunded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	create := func(t *testing.T, orderID string, status entity.PaymentStatus, now time.Time) {
		order := testOrder(orderID, "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now)
		err := db.DB.Create(&order).Error
		require.NoError(t, err)

		payment := testOrderPayment(orderID, 1, "transaction-id", "payment-id", now)
		payment.Status = status
		err = db.DB.Create(&payment).Error
		require.NoError(t, err)

		fulfillments := make(entity.OrderFulfillments, 1)
		fulfillments[0] = testOrderFulfillment("fulfillment-id", orderID, 1, 1, now)
		err = db.DB.Create(&fulfillments).Error
		require.NoError(t, err)

		items := make(entity.OrderItems, 2)
		items[0] = testOrderItem("fulfillment-id", 1, orderID, now)
		items[1] = testOrderItem("fulfillment-id", 2, orderID, now)
		err = db.DB.Create(&items).Error
		require.NoError(t, err)
	}

	type args struct {
		orderID string
		params  *database.UpdateOrderRefundedParams
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success canceled",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.PaymentStatusPending, now().AddDate(0, 0, -1))
			},
			args: args{
				orderID: "order-id",
				params: &database.UpdateOrderRefundedParams{
					Status:       entity.PaymentStatusCanceled,
					RefundType:   entity.RefundTypeCanceled,
					RefundTotal:  1980,
					RefundReason: "テストです。",
					IssuedAt:     now(),
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success refunded",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.PaymentStatusPending, now().AddDate(0, 0, -1))
			},
			args: args{
				orderID: "order-id",
				params: &database.UpdateOrderRefundedParams{
					Status:       entity.PaymentStatusRefunded,
					RefundType:   entity.RefundTypeRefunded,
					RefundTotal:  1980,
					RefundReason: "テストです。",
					IssuedAt:     now(),
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				orderID: "",
				params: &database.UpdateOrderRefundedParams{
					Status:   entity.PaymentStatusRefunded,
					IssuedAt: now(),
				},
			},
			want: want{
				err: database.ErrNotFound,
			},
		},
		{
			name: "not latest data",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.PaymentStatusAuthorized, now().AddDate(0, 0, 1))
			},
			args: args{
				orderID: "order-id",
				params: &database.UpdateOrderRefundedParams{
					Status:   entity.PaymentStatusCaptured,
					IssuedAt: now(),
				},
			},
			want: want{
				err: database.ErrFailedPrecondition,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, orderItemTable, orderFulfillmentTable, orderPaymentTable, orderExperienceTable, orderTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			err = db.UpdateRefunded(ctx, tt.args.orderID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestOrder_UpdateFulfillment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	create := func(t *testing.T, orderID string, status entity.OrderStatus, now time.Time) {
		order := testOrder(orderID, "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now)
		order.Status = status
		err := db.DB.Create(&order).Error
		require.NoError(t, err)

		payment := testOrderPayment(orderID, 1, "transaction-id", "payment-id", now)
		err = db.DB.Create(&payment).Error
		require.NoError(t, err)

		fulfillments := make(entity.OrderFulfillments, 1)
		fulfillments[0] = testOrderFulfillment("fulfillment-id", orderID, 1, 1, now)
		err = db.DB.Create(&fulfillments).Error
		require.NoError(t, err)

		items := make(entity.OrderItems, 2)
		items[0] = testOrderItem("fulfillment-id", 1, orderID, now)
		items[1] = testOrderItem("fulfillment-id", 2, orderID, now)
		err = db.DB.Create(&items).Error
		require.NoError(t, err)
	}

	type args struct {
		orderID       string
		fulfillmentID string
		params        *database.UpdateOrderFulfillmentParams
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success fulfilled",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.OrderStatusPreparing, now().AddDate(0, 0, -1))
			},
			args: args{
				orderID:       "order-id",
				fulfillmentID: "fulfillment-id",
				params: &database.UpdateOrderFulfillmentParams{
					Status:          entity.FulfillmentStatusFulfilled,
					ShippingCarrier: entity.ShippingCarrierYamato,
					TrackingNumber:  "tracking-number",
					ShippedAt:       now(),
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success uunfulfilled",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.OrderStatusPreparing, now().AddDate(0, 0, -1))
			},
			args: args{
				orderID:       "order-id",
				fulfillmentID: "fulfillment-id",
				params: &database.UpdateOrderFulfillmentParams{
					Status:          entity.FulfillmentStatusUnfulfilled,
					ShippingCarrier: entity.ShippingCarrierYamato,
					TrackingNumber:  "tracking-number",
					ShippedAt:       now(),
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				orderID:       "",
				fulfillmentID: "fulfillment-id",
				params: &database.UpdateOrderFulfillmentParams{
					Status:          entity.FulfillmentStatusFulfilled,
					ShippingCarrier: entity.ShippingCarrierYamato,
					TrackingNumber:  "tracking-number",
					ShippedAt:       now(),
				},
			},
			want: want{
				err: database.ErrNotFound,
			},
		},
		{
			name: "already completed",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.OrderStatusCompleted, now().AddDate(0, 0, -1))
			},
			args: args{
				orderID:       "order-id",
				fulfillmentID: "fulfillment-id",
				params: &database.UpdateOrderFulfillmentParams{
					Status:          entity.FulfillmentStatusFulfilled,
					ShippingCarrier: entity.ShippingCarrierYamato,
					TrackingNumber:  "tracking-number",
					ShippedAt:       now(),
				},
			},
			want: want{
				err: database.ErrFailedPrecondition,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, orderItemTable, orderFulfillmentTable, orderPaymentTable, orderExperienceTable, orderTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			err = db.UpdateFulfillment(ctx, tt.args.orderID, tt.args.fulfillmentID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestOrder_Draft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	create := func(t *testing.T, orderID string, status entity.PaymentStatus, now time.Time) {
		order := testOrder(orderID, "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now)
		err := db.DB.Create(&order).Error
		require.NoError(t, err)

		payment := testOrderPayment(orderID, 1, "transaction-id", "payment-id", now)
		payment.Status = status
		err = db.DB.Create(&payment).Error
		require.NoError(t, err)

		fulfillments := make(entity.OrderFulfillments, 1)
		fulfillments[0] = testOrderFulfillment("fulfillment-id", orderID, 1, 1, now)
		err = db.DB.Create(&fulfillments).Error
		require.NoError(t, err)

		items := make(entity.OrderItems, 2)
		items[0] = testOrderItem("fulfillment-id", 1, orderID, now)
		items[1] = testOrderItem("fulfillment-id", 2, orderID, now)
		err = db.DB.Create(&items).Error
		require.NoError(t, err)
	}

	type args struct {
		orderID string
		params  *database.DraftOrderParams
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.PaymentStatusPending, now().AddDate(0, 0, -1))
			},
			args: args{
				orderID: "order-id",
				params: &database.DraftOrderParams{
					ShippingMessage: "購入ありがとうございます。",
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, orderItemTable, orderFulfillmentTable, orderPaymentTable, orderExperienceTable, orderTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			err = db.Draft(ctx, tt.args.orderID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestOrder_Complete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	create := func(t *testing.T, orderID string, status entity.PaymentStatus, now time.Time) {
		order := testOrder(orderID, "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now)
		err := db.DB.Create(&order).Error
		require.NoError(t, err)

		payment := testOrderPayment(orderID, 1, "transaction-id", "payment-id", now)
		payment.Status = status
		err = db.DB.Create(&payment).Error
		require.NoError(t, err)

		fulfillments := make(entity.OrderFulfillments, 1)
		fulfillments[0] = testOrderFulfillment("fulfillment-id", orderID, 1, 1, now)
		err = db.DB.Create(&fulfillments).Error
		require.NoError(t, err)

		items := make(entity.OrderItems, 2)
		items[0] = testOrderItem("fulfillment-id", 1, orderID, now)
		items[1] = testOrderItem("fulfillment-id", 2, orderID, now)
		err = db.DB.Create(&items).Error
		require.NoError(t, err)
	}

	type args struct {
		orderID string
		params  *database.CompleteOrderParams
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				create(t, "order-id", entity.PaymentStatusPending, now().AddDate(0, 0, -1))
			},
			args: args{
				orderID: "order-id",
				params: &database.CompleteOrderParams{
					ShippingMessage: "購入ありがとうございます。",
					CompletedAt:     now(),
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, orderItemTable, orderFulfillmentTable, orderPaymentTable, orderExperienceTable, orderTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			err = db.Complete(ctx, tt.args.orderID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestOrder_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now())
	orders[1] = testOrder("order-id02", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 2, now())
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
		params *database.AggregateOrdersParams
	}
	type want struct {
		order *entity.AggregatedOrder
		err   error
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
				params: &database.AggregateOrdersParams{
					ShopID:       "shop-id",
					CreatedAtGte: now().AddDate(0, 0, -1),
					CreatedAtLt:  now().AddDate(0, 0, 1),
				},
			},
			want: want{
				order: &entity.AggregatedOrder{
					OrderCount:    2,
					UserCount:     1,
					SalesTotal:    3600,
					DiscountTotal: 0,
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.Aggregate(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.order, actual)
		})
	}
}

func TestOrder_AggregateByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now())
	orders[1] = testOrder("order-id02", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 2, now())
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
		params *database.AggregateOrdersByUserParams
	}
	type want struct {
		orders entity.AggregatedUserOrders
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
				params: &database.AggregateOrdersByUserParams{
					ShopID:  "shop-id",
					UserIDs: []string{"user-id", "other-id"},
				},
			},
			want: want{
				orders: entity.AggregatedUserOrders{
					{
						UserID:     "user-id",
						OrderCount: 2,
						Subtotal:   3600,
						Discount:   0,
						Total:      4600,
					},
				},
				hasErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.AggregateByUser(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.orders, actual)
		})
	}
}

func TestOrder_AggregateByPaymentMethodType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now())
	orders[1] = testOrder("order-id02", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 2, now())
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
		params *database.AggregateOrdersByPaymentMethodTypeParams
	}
	type want struct {
		orders entity.AggregatedOrderPayments
		err    error
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
				params: &database.AggregateOrdersByPaymentMethodTypeParams{
					ShopID:             "shop-id",
					PaymentMethodTypes: entity.AllPaymentMethodTypes,
				},
			},
			want: want{
				orders: entity.AggregatedOrderPayments{
					{
						PaymentMethodType: entity.PaymentMethodTypeCreditCard,
						OrderCount:        2,
						UserCount:         1,
						SalesTotal:        3600,
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.AggregateByPaymentMethodType(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.orders, actual)
		})
	}
}

func TestOrder_AggregateByPromotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)
	promotion := testPromotion("promotion-id", "12345678", "", now())
	err = db.DB.Create(&promotion).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "promotion-id", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now())
	orders[1] = testOrder("order-id02", "user-id", "promotion-id", "shop-id", "coordinator-id", entity.OrderTypeProduct, 2, now())
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
		params *database.AggregateOrdersByPromotionParams
	}
	type want struct {
		orders entity.AggregatedOrderPromotions
		err    error
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
				params: &database.AggregateOrdersByPromotionParams{
					ShopID:       "shop-id",
					PromotionIDs: []string{"promotion-id"},
				},
			},
			want: want{
				orders: entity.AggregatedOrderPromotions{
					{
						PromotionID:   "promotion-id",
						OrderCount:    2,
						DiscountTotal: 0,
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.AggregateByPromotion(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.orders, actual)
		})
	}
}

func TestOrder_AggregateByPeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
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
	pinternal := make(internalProducts, 2)
	pinternal[0] = testProduct("product-id01", "type-id01", "shop-id", "coordinator-id", "producer-id", []string{}, 1, now())
	pinternal[1] = testProduct("product-id02", "type-id02", "shop-id", "coordinator-id", "producer-id", []string{}, 2, now())
	err = db.DB.Table(productTable).Create(&pinternal).Error
	require.NoError(t, err)
	for i := range pinternal {
		err = db.DB.Create(&pinternal[i].ProductRevision).Error
		require.NoError(t, err)
	}
	schedule := testSchedule("schedule-id", "shop-id", "coordinator-id", now())
	err = db.DB.Create(&schedule).Error
	require.NoError(t, err)

	orders := make(entity.Orders, 2)
	orders[0] = testOrder("order-id01", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 1, now())
	orders[1] = testOrder("order-id02", "user-id", "", "shop-id", "coordinator-id", entity.OrderTypeProduct, 2, now())
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
		params *database.AggregateOrdersByPeriodParams
	}
	type want struct {
		orders entity.AggregatedPeriodOrders
		err    error
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
				params: &database.AggregateOrdersByPeriodParams{
					ShopID:       "shop-id",
					PeriodType:   entity.AggregateOrderPeriodTypeDay,
					CreatedAtGte: now().AddDate(0, 0, -1),
					CreatedAtLt:  now().AddDate(0, 0, 1),
				},
			},
			want: want{
				orders: entity.AggregatedPeriodOrders{
					{
						Period:        jst.Date(now().Year(), now().Month(), now().Day(), 0, 0, 0, 0),
						OrderCount:    2,
						UserCount:     1,
						SalesTotal:    3600,
						DiscountTotal: 0,
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &order{db: db, now: now}
			actual, err := db.AggregateByPeriod(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.orders, actual)
		})
	}
}

func testOrder(id, userID, promotionID, shopID, coordinatorID string, typ entity.OrderType, mgmtID int64, now time.Time) *entity.Order {
	return &entity.Order{
		ID:            id,
		SessionID:     "session-id",
		UserID:        userID,
		PromotionID:   promotionID,
		ShopID:        shopID,
		CoordinatorID: coordinatorID,
		ManagementID:  mgmtID,
		Type:          typ,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func testOrderPayment(orderID string, addressID int64, transactionID, paymentID string, now time.Time) *entity.OrderPayment {
	return &entity.OrderPayment{
		OrderID:           orderID,
		AddressRevisionID: addressID,
		Status:            entity.PaymentStatusCaptured,
		TransactionID:     transactionID,
		PaymentID:         paymentID,
		MethodType:        entity.PaymentMethodTypeCreditCard,
		Subtotal:          1800,
		Discount:          0,
		ShippingFee:       500,
		Tax:               209,
		Total:             2300,
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

func testOrderExperience(orderID string, experienceID int64, now time.Time) *internalOrderExperience {
	experience := &entity.OrderExperience{
		OrderID:               orderID,
		ExperienceRevisionID:  experienceID,
		AdultCount:            1,
		JuniorHighSchoolCount: 1,
		ElementarySchoolCount: 2,
		PreschoolCount:        0,
		SeniorCount:           0,
		Remarks: entity.OrderExperienceRemarks{
			Transportation: "電車",
			RequestedDate:  jst.Date(2024, 1, 2, 0, 0, 0, 0),
			RequestedTime:  jst.Date(0, 1, 1, 18, 30, 0, 0),
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	internal, _ := newInternalOrderExperience(experience)
	return internal
}
