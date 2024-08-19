package grpc

import (
	"context"

	"test-go/internal/adapters/primary/grpc/proto"
	"test-go/internal/application"
)

// ProductHandler implements the gRPC server interface for managing products
type ProductHandler struct {
	proto.UnimplementedProductServiceServer
	service *application.ProductService
}

// NewProductHandler creates a new instance of ProductHandler
func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProduct handles the creation of a new product via gRPC
func (h *ProductHandler) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	id, err := h.service.CreateProduct(ctx, req.Name, req.Price)
	if err != nil {
		return nil, err
	}

	return &proto.CreateProductResponse{Id: id}, nil
}

// GetProductByID retrieves a product by its ID via gRPC
func (h *ProductHandler) GetProductByID(ctx context.Context, req *proto.GetProductByIDRequest) (*proto.GetProductByIDResponse, error) {
	product, err := h.service.GetProductByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.GetProductByIDResponse{
		Product: &proto.Product{
			Id:    product.ID.Hex(),
			Name:  product.Name,
			Price: product.Price,
		},
	}, nil
}

// UpdateProduct handles updating an existing product via gRPC
func (h *ProductHandler) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	err := h.service.UpdateProduct(ctx, req.Product.Id, req.Product.Name, req.Product.Price)
	if err != nil {
		return nil, err
	}

	return &proto.UpdateProductResponse{Success: true}, nil
}

// DeleteProduct handles deleting a product by its ID via gRPC
func (h *ProductHandler) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error) {
	err := h.service.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteProductResponse{Success: true}, nil
}

// ListProducts retrieves all products via gRPC
func (h *ProductHandler) ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error) {
	products, err := h.service.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	// Convert the list of products to the protobuf format
	var protoProducts []*proto.Product
	for _, product := range products {
		protoProducts = append(protoProducts, &proto.Product{
			Id:    product.ID.Hex(),
			Name:  product.Name,
			Price: product.Price,
		})
	}

	return &proto.ListProductsResponse{Products: protoProducts}, nil
}
