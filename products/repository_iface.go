package products

import "context"

/*
mockery --dir=products --name=IProduct --filename=repository_mocks.go --output=products/mocks --outpkg=repomocks
*/
type IProduct interface {
	InsertProduct(ctx context.Context, product *ProductModel) (*ProductModel, error)
	FindAllProduct(ctx context.Context) []ProductModel
	FindProductsWithPagination(ctx context.Context, page int64, pageSize int64) ([]ProductModel, error)
	FindProductWithID(ctx context.Context, idProduct string) (*ProductModel, error)
	UpdateProduct(ctx context.Context, idProduct string, product *ProductModel) (*ProductModel, error)
	DeleteProduct(ctx context.Context, idProduct string) error
}
