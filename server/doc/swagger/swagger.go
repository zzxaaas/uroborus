package swagger

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
)

// Serve serves the documentation.
func Serve(c *gin.Context) {
	c.Writer.WriteString(get(c.Request.Host))
}

func get(host string) string {
	log.Println(host)
	box := packr.NewBox("./files") // v1
	// box := packr.New("swagger", "./server/swagger/files")  // v2
	json, err := box.FindString("swagger.json")
	if err != nil {
		log.Fatal(err)
		return "error finding swagger.json"
	}
	return strings.Replace(json, "localhost", host, 1)
}
