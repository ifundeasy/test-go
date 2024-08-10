package product

// Product represents the product entity with its attributes.
type Product struct {
	ID    string  `json:"id" bson:"_id,omitempty" validate:"omitempty"`
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required,gt=0"`
}

// ProductRepository defines the interface for product repository operations.
type ProductRepository interface {
	Fetch() ([]Product, error)
	GetByID(id string) (*Product, error)
	Create(product *Product) error
	Update(id string, product *Product) error
	Delete(id string) error
}

// ProductUsecase defines the interface for product use case operations.
type ProductUsecase interface {
	Fetch() ([]Product, error)
	GetByID(id string) (*Product, error)
	Create(product *Product) error
	Update(id string, product *Product) error
	Delete(id string) error
}
