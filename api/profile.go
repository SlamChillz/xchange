package api

import (
	"bytes"
	"fmt"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
)

func (server *Server) CustomerProfilePicture(c *gin.Context) {
	server.router.MaxMultipartMemory = 8 << 20 // 8 MiB
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	buf := make([]byte, 0, file.Size)
	buffer := bytes.NewBuffer(buf)
	openFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer openFile.Close()
	buffer.ReadFrom(openFile)
	log.Printf("%+v\n", buffer.Len())
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
