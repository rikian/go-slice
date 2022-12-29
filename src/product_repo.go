package src

import (
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
)

var database = map[string]Product{}

type ProductRepo interface {
	GetProduct() *map[string]Product
	PostProduct(pName, imgName string) string
	DeleteProduct(pid, pName string) (string, error)
}

type productRepoImpl struct {
	product Product
}

func NewProductRepo() ProductRepo {
	return &productRepoImpl{
		product: Product{},
	}
}

func (p *productRepoImpl) GetProduct() *map[string]Product {
	return &database
}

func (p *productRepoImpl) PostProduct(pName, imgName string) string {
	id := uuid.New().String()

	p.product.Id = id
	p.product.Name = pName
	p.product.Image = imgName

	database[id] = p.product

	return id
}

func (p *productRepoImpl) DeleteProduct(pid, pName string) (string, error) {
	for key, product := range database {
		if product.Id == pid && product.Name == pName {
			err := os.Remove("./assets/images/" + product.Image)
			if os.IsNotExist(err) {
				log.Print("Image file not exist...")
			}

			delete(database, key)
			return "ok", nil
		}
	}

	return "not ok", errors.New("product not found")
}
