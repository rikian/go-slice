package src

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

type Listener interface {
	Run(addres, port string)
}

type listenerImpl struct {
	s ProductService
}

func NewListener() Listener {
	return &listenerImpl{
		s: NewProductService(),
	}
}

func clearDirImages() {
	folderPath := "./assets/images"

	// check folder images exist or not
	_, err := os.Stat(folderPath)

	exists := os.IsExist(err)

	if !exists {
		// Delete the folder if exist and all its contents
		err := os.RemoveAll(folderPath)
		if err != nil {
			log.Printf("Error deleting folder: %v\n", err)
			return
		}
	}

	// create dir images
	err = os.Mkdir(folderPath, 0755)
	if err != nil {
		log.Printf("Error creating folder: %v\n", err)
		return
	}
}

func (l *listenerImpl) Run(addres, port string) {
	clearDirImages()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(gin.Recovery())
	router.LoadHTMLGlob("./views/*.html")
	router.Static("/assets", "./assets")
	router.GET("/", l.s.GetService())
	router.POST("/products/post", l.s.PostService())
	router.POST("/products/delete", l.s.DeleteService())
	router.POST("/products/update", l.s.UpdateService())

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := router.Run(addres + ":" + port)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	log.Print("server running at " + addres + ":" + port)

	<-done

	log.Print("server is going to stop")
	clearDirImages()
}
