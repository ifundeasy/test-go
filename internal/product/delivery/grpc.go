package delivery

import (
	"context"

	"test-go/internal/product"
	"test-go/internal/product/pb"
)

// GRPCProductHandler is the gRPC handler for products
type GRPCProductHandler struct {
	pb.UnimplementedProductServiceServer
	usecase product.ProductUsecase
}

// NewGRPCProductHandler returns a new instance of GRPCProductHandler
func NewGRPCProductHandler(u product.ProductUsecase) *GRPCProductHandler {
	return &GRPCProductHandler{usecase: u}
}

// GetProduct handles the gRPC request to get a product by ID
func (h *GRPCProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	product, err := h.usecase.GetByID(req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.ProductResponse{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}

// CreateProduct handles the gRPC request to create a new product
func (h *GRPCProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product := &product.Product{
		Name:  req.GetName(),
		Price: req.GetPrice(),
	}

	err := h.usecase.Create(product)
	if err != nil {
		return nil, err
	}

	return &pb.ProductResponse{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}

// UpdateProduct handles the gRPC request to update an existing product
func (h *GRPCProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	product := &product.Product{
		ID:    req.GetId(),
		Name:  req.GetName(),
		Price: req.GetPrice(),
	}

	err := h.usecase.Update(product.ID, product)
	if err != nil {
		return nil, err
	}

	return &pb.ProductResponse{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}

// DeleteProduct handles the gRPC request to delete a product by ID
func (h *GRPCProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Empty, error) {
	err := h.usecase.Delete(req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
