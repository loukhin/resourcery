package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func regenerate(resourcePath string) (string, string) {
	fmt.Println("Regenerating...")
	fileHash := generateResourcePack(resourcePath)
	filePath := path.Join("cache", fileHash+".zip")
	fmt.Println("Regenerated!")
	return fileHash, filePath
}

func main() {
	var fileHash, filePath string
	resourcePath := os.Args[1]

	fileHash, filePath = regenerate(resourcePath)

	go watch(resourcePath, func() {
		fmt.Println("Changes detected!")
		fileHash, filePath = regenerate(resourcePath)
	})

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/resourcepack/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/resourcepack?h="+fileHash)
	})

	r.GET("/resourcepack", func(c *gin.Context) {
		c.FileAttachment(filePath, path.Base(filePath))
	})

	r.GET("/resourcepack.sha1", func(c *gin.Context) {
		c.String(200, fileHash)
	})

	err := r.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}
