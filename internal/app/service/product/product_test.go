package sproduct

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/kenyako/catalog-service/internal/app/entity"
	"github.com/kenyako/catalog-service/internal/app/repository/mocks"
	"github.com/kenyako/catalog-service/internal/pkg/testutil"
)

func TestNewService(t *testing.T) {
	productRepo := mocks.NewMockProduct(t)
	categoryRepo := mocks.NewMockCategory(t)

	result := NewService(productRepo, categoryRepo)

	require.NotNil(t, result)
}

type createProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *createProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestCreateProductSuite(t *testing.T) {
	suite.Run(t, new(createProductSuite))
}

func (s *createProductSuite) TestCreate() {
	type args struct {
		req entity.RequestProductCreate
	}

	type want struct {
		err error
	}

	categoryGUID := uuid.Must(uuid.NewV4())

	listErr := errors.New("list error")
	categoryErr := errors.New("category error")
	createErr := errors.New("create error")

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Description:  testutil.PtrString("A test product"),
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: nil},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{{GUID: categoryGUID}}, nil).
					Once()

				s.productRepo.EXPECT().
					Create(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return p.Name == args.req.Name &&
							p.Description == args.req.Description &&
							p.Price == args.req.Price &&
							p.CategoryGUID == args.req.CategoryGUID
					})).
					Return(nil).
					Once()
			},
		},
		{
			name: "already exists",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Existing Product",
					Price:        500,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: entity.ErrAlreadyExists},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{{Name: "Existing Product"}}, nil).
					Once()
			},
		},
		{
			name: "category not found",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "New Product",
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{}, nil).
					Once()
			},
		},
		{
			name: "list error",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: listErr},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return(nil, listErr).
					Once()
			},
		},
		{
			name: "category get error",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: categoryErr},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return(nil, categoryErr).
					Once()
			},
		},
		{
			name: "create error",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Description:  testutil.PtrString("A test product"),
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: createErr},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{{GUID: categoryGUID}}, nil).
					Once()

				s.productRepo.EXPECT().
					Create(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return p.Name == args.req.Name &&
							p.Description == args.req.Description &&
							p.Price == args.req.Price &&
							p.CategoryGUID == args.req.CategoryGUID
					})).
					Return(createErr).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.Create(s.ctx, tc.args.req)

			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
			} else {
				s.NoError(err)
				s.NotEmpty(result.GUID)
				s.Equal(tc.args.req.Name, result.Name)
				s.Equal(tc.args.req.Description, result.Description)
				s.Equal(tc.args.req.Price, result.Price)
				s.Equal(tc.args.req.CategoryGUID, result.CategoryGUID)
			}
		})
	}
}

type getByGUIDsProductSuite struct {
	suite.Suite
	srv         *srv
	productRepo *mocks.MockProduct
	ctx         context.Context
}

func (s *getByGUIDsProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.srv = &srv{
		repoProduct: s.productRepo,
	}
}

func TestGetByGUIDsProductSuite(t *testing.T) {
	suite.Run(t, new(getByGUIDsProductSuite))
}

func (s *getByGUIDsProductSuite) TestGetByGUIDs() {
	type args struct {
		guids []uuid.UUID
	}

	type want struct {
		products []entity.Product
		err      error
	}

	productGUID := uuid.Must(uuid.NewV4())
	categoryGUID := uuid.Must(uuid.NewV4())

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "single product found",
			args: args{
				guids: []uuid.UUID{productGUID},
			},
			want: want{
				products: []entity.Product{
					{
						GUID:         productGUID,
						Name:         "Test Product",
						Price:        1000,
						CategoryGUID: categoryGUID,
					},
				},
				err: nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, args.guids).
					Return([]entity.Product{
						{
							GUID:         productGUID,
							Name:         "Test Product",
							Price:        1000,
							CategoryGUID: categoryGUID,
						},
					}, nil).
					Once()
			},
		},
		{
			name: "not found returns empty slice",
			args: args{
				guids: []uuid.UUID{productGUID},
			},
			want: want{
				products: []entity.Product{},
				err:      nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, args.guids).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.GetByGUIDs(s.ctx, tc.args.guids)

			s.NoError(err)
			s.Equal(tc.want.products, result)
		})
	}
}

type deleteProductSuite struct {
	suite.Suite
	srv         *srv
	productRepo *mocks.MockProduct
	ctx         context.Context
}

func (s *deleteProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.srv = &srv{
		repoProduct: s.productRepo,
	}
}

func TestDeleteProductSuite(t *testing.T) {
	suite.Run(t, new(deleteProductSuite))
}

func (s *deleteProductSuite) TestDelete() {
	type args struct {
		guid uuid.UUID
	}

	type want struct {
		err error
	}

	productGUID := uuid.Must(uuid.NewV4())
	deleteErr := errors.New("delete product error")
	getErr := errors.New("get product error")

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				guid: productGUID,
			},
			want: want{
				err: nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{
						{
							GUID: args.guid,
						},
					}, nil).
					Once()

				s.productRepo.EXPECT().
					Delete(s.ctx, args.guid).
					Return(nil).
					Once()
			},
		},
		{
			name: "not found",
			args: args{
				guid: productGUID,
			},
			want: want{
				err: entity.ErrNotFound,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
		{
			name: "delete error",
			args: args{
				guid: productGUID,
			},
			want: want{
				err: deleteErr,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{
						{GUID: args.guid},
					}, nil).
					Once()

				s.productRepo.EXPECT().
					Delete(s.ctx, args.guid).
					Return(deleteErr).
					Once()
			},
		},
		{
			name: "get product error",
			args: args{
				guid: productGUID,
			},
			want: want{
				err: getErr,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(
						s.ctx,
						mock.AnythingOfType("func(context.Context) error"),
					).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.guid},
					).
					Return(nil, getErr).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			err := s.srv.Delete(s.ctx, tc.args.guid)

			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

type listProductSuite struct {
	suite.Suite
	srv         *srv
	productRepo *mocks.MockProduct
	ctx         context.Context
}

func (s *listProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.srv = &srv{
		repoProduct: s.productRepo,
	}
}

func TestListProductSuite(t *testing.T) {
	suite.Run(t, new(listProductSuite))
}

func (s *listProductSuite) TestList() {
	type args struct {
		req entity.RequestProductList
	}

	type want struct {
		products []entity.Product
		err      error
	}

	categoryGUID := uuid.Must(uuid.NewV4())

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				req: entity.RequestProductList{
					CategoryGUID: &categoryGUID,
					MinPrice:     testutil.PtrInt64(500),
					MaxPrice:     testutil.PtrInt64(1500),
				},
			},
			want: want{
				err: nil,
				products: []entity.Product{
					{
						Name:  "Product 1",
						Price: 1000,
					},
					{
						Name:  "Product 2",
						Price: 1200,
					},
				},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					List(
						s.ctx,
						(*string)(nil),
						args.req.CategoryGUID,
						args.req.MinPrice,
						args.req.MaxPrice,
					).
					Return([]entity.Product{
						{
							Name:  "Product 1",
							Price: 1000,
						},
						{
							Name:  "Product 2",
							Price: 1200,
						},
					}, nil).
					Once()
			},
		},
		{
			name: "empty result",
			args: args{
				req: entity.RequestProductList{
					CategoryGUID: &categoryGUID,
					MinPrice:     testutil.PtrInt64(500),
					MaxPrice:     testutil.PtrInt64(1500),
				},
			},
			want: want{
				err:      nil,
				products: []entity.Product{},
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					List(
						s.ctx,
						(*string)(nil),
						args.req.CategoryGUID,
						args.req.MinPrice,
						args.req.MaxPrice,
					).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.List(s.ctx, tc.args.req)

			s.NoError(err)
			s.Len(result, len(tc.want.products))

			for i, product := range tc.want.products {
				s.Equal(product.Name, result[i].Name)
				s.Equal(product.Price, result[i].Price)
			}
		})
	}
}

type updateProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *updateProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestUpdateProductSuite(t *testing.T) {
	suite.Run(t, new(updateProductSuite))
}

func (s *updateProductSuite) TestUpdate() {
	type args struct {
		guid uuid.UUID
		req  entity.RequestProductUpdate
	}

	type want struct {
		err error
	}

	productGUID := uuid.Must(uuid.NewV4())
	oldCategoryGUID := uuid.Must(uuid.NewV4())
	newCategoryGUID := uuid.Must(uuid.NewV4())
	anotherProductGUID := uuid.Must(uuid.NewV4())

	oldProduct := entity.Product{
		GUID:         productGUID,
		Name:         "Old Name",
		Description:  testutil.PtrString("Old description"),
		Price:        1000,
		CategoryGUID: oldCategoryGUID,
	}

	getProductErr := errors.New("get product error")
	listErr := errors.New("list error")
	categoryErr := errors.New("category error")
	updateErr := errors.New("update error")

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "full update",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name:         "New Name",
					Description:  testutil.PtrString("New description"),
					Price:        testutil.PtrInt64(2000),
					CategoryGUID: newCategoryGUID,
				},
			},
			want: want{
				err: nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(
						s.ctx,
						mock.AnythingOfType("func(context.Context) error"),
					).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.guid},
					).
					Return([]entity.Product{oldProduct}, nil).
					Once()

				s.productRepo.EXPECT().
					List(
						s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil),
					).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.req.CategoryGUID},
					).
					Return([]entity.Category{
						{
							GUID: args.req.CategoryGUID,
						},
					}, nil).
					Once()

				s.productRepo.EXPECT().
					Update(
						s.ctx,
						mock.MatchedBy(func(p entity.Product) bool {
							return p.GUID == args.guid &&
								p.Name == args.req.Name &&
								p.Description == args.req.Description &&
								p.Price == *args.req.Price &&
								p.CategoryGUID == args.req.CategoryGUID
						}),
					).
					Return(nil).
					Once()
			},
		},
		{
			name: "partial update - name only",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name: "New Name",
				},
			},
			want: want{
				err: nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(
						s.ctx,
						mock.AnythingOfType("func(context.Context) error"),
					).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.guid},
					).
					Return([]entity.Product{oldProduct}, nil).
					Once()

				s.productRepo.EXPECT().
					List(
						s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil),
					).
					Return([]entity.Product{}, nil).
					Once()

				s.productRepo.EXPECT().
					Update(
						s.ctx,
						mock.MatchedBy(func(p entity.Product) bool {
							return p.GUID == args.guid &&
								p.Name == args.req.Name &&
								p.Description == oldProduct.Description &&
								p.Price == oldProduct.Price &&
								p.CategoryGUID == oldProduct.CategoryGUID
						}),
					).
					Return(nil).
					Once()
			},
		},
		{
			name: "not found",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name: "New Name",
				},
			},
			want: want{
				err: entity.ErrNotFound,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(
						s.ctx,
						mock.AnythingOfType("func(context.Context) error"),
					).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.guid},
					).
					Return([]entity.Product{}, nil).
					Once()
			},
		},
		{
			name: "duplicate name",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name: "Existing Name",
				},
			},
			want: want{
				err: entity.ErrAlreadyExists,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(
						s.ctx,
						mock.AnythingOfType("func(context.Context) error"),
					).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.guid},
					).
					Return([]entity.Product{oldProduct}, nil).
					Once()

				s.productRepo.EXPECT().
					List(
						s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil),
					).
					Return([]entity.Product{
						{
							GUID: anotherProductGUID,
							Name: args.req.Name,
						},
					}, nil).
					Once()
			},
		},
		{
			name: "category not found",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name:         "New Name",
					CategoryGUID: newCategoryGUID,
				},
			},
			want: want{
				err: entity.ErrNotFound,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(
						s.ctx,
						mock.AnythingOfType("func(context.Context) error"),
					).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.guid},
					).
					Return([]entity.Product{oldProduct}, nil).
					Once()

				s.productRepo.EXPECT().
					List(
						s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil),
					).
					Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.req.CategoryGUID},
					).
					Return([]entity.Category{}, nil).
					Once()
			},
		},
		{
			name: "get product error",
			args: args{
				guid: productGUID,
				req:  entity.RequestProductUpdate{},
			},
			want: want{
				err: getProductErr,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(
						s.ctx,
						mock.AnythingOfType("func(context.Context) error"),
					).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.guid},
					).
					Return(nil, getProductErr).
					Once()
			},
		},
		{
			name: "list error",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Name: "New Name",
				},
			},
			want: want{
				err: listErr,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(
						s.ctx,
						mock.AnythingOfType("func(context.Context) error"),
					).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.guid},
					).
					Return([]entity.Product{oldProduct}, nil).
					Once()

				s.productRepo.EXPECT().
					List(
						s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil),
					).
					Return(nil, listErr).
					Once()
			},
		},
		{
			name: "category get error",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					CategoryGUID: newCategoryGUID,
				},
			},
			want: want{
				err: categoryErr,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(
						s.ctx,
						mock.AnythingOfType("func(context.Context) error"),
					).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.guid},
					).
					Return([]entity.Product{oldProduct}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.req.CategoryGUID},
					).
					Return(nil, categoryErr).
					Once()
			},
		},
		{
			name: "update error",
			args: args{
				guid: productGUID,
				req: entity.RequestProductUpdate{
					Price: testutil.PtrInt64(2000),
				},
			},
			want: want{
				err: updateErr,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					InsideTx(
						s.ctx,
						mock.AnythingOfType("func(context.Context) error"),
					).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(
						s.ctx,
						[]uuid.UUID{args.guid},
					).
					Return([]entity.Product{oldProduct}, nil).
					Once()

				s.productRepo.EXPECT().
					Update(
						s.ctx,
						mock.MatchedBy(func(p entity.Product) bool {
							return p.GUID == oldProduct.GUID &&
								p.Name == oldProduct.Name &&
								p.Description == oldProduct.Description &&
								p.Price == *args.req.Price &&
								p.CategoryGUID == oldProduct.CategoryGUID
						}),
					).
					Return(updateErr).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.Update(s.ctx, tc.args.guid, tc.args.req)

			if tc.want.err != nil {
				s.ErrorIs(err, tc.want.err)
			} else {
				s.NoError(err)
				s.Equal(tc.args.guid, result.GUID)
			}
		})
	}
}
