package src

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductService interface {
	GetService() gin.HandlerFunc
	PostService() gin.HandlerFunc
	DeleteService() gin.HandlerFunc
	UpdateService() gin.HandlerFunc
}

type productServiceImpl struct {
	s ProductRepo
}

func NewProductService() ProductService {
	return &productServiceImpl{
		s: NewProductRepo(),
	}
}

func (p *productServiceImpl) GetService() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		database := p.s.GetProduct()

		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"ProductList": database,
		})
	}
}

func (p *productServiceImpl) PostService() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 1000000)
		file, err := ctx.FormFile("product_image")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
			})

			return
		}

		pName := ctx.PostForm("product_name")

		if pName == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
			})

			return
		}

		imgName := uuid.New().String()

		switch strings.Split(file.Header.Get("content-type"), "/")[1] {
		case "jpeg":
			imgName += ".jpeg"
		case "jpg":
			imgName += ".jpg"
		case "png":
			imgName += ".png"
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
			})

			return
		}

		postProduck := p.s.PostProduct(pName, imgName)

		err = ctx.SaveUploadedFile(file, "./assets/images/"+imgName)

		if err != nil {
			log.Print(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
			})

			p.s.DeleteProduct(postProduck, pName)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

func (p *productServiceImpl) DeleteService() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productId := ctx.PostForm("product_id")
		productName := ctx.PostForm("product_name")

		deleteProduct, err := p.s.DeleteProduct(productId, productName)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
			})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": deleteProduct,
		})
	}
}

func (p *productServiceImpl) UpdateService() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
		})
	}
}
