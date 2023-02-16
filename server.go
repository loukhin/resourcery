package main

import (
	"log"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func main() {
	var fileHash, filePath string
	resourcePath := os.Args[1]

	fileHash = generateResourcePack(resourcePath)
	filePath = path.Join("cache", fileHash+".zip")

	go watch(resourcePath, func() {
		fileHash = generateResourcePack(resourcePath)
		filePath = path.Join("cache", fileHash+".zip")
	})

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

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
