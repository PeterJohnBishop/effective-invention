package server

import (
	"bytes"
	"effective-invention/server/builder"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func logAsPrettyJSON(label string, data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("%s: error marshaling JSON: %v", label, err)
		return
	}
	log.Printf("%s:\n%s", label, string(jsonBytes))
}

// prints raw JSON, essentially a webhook endpoint to view payloads
func HandleRawData() gin.HandlerFunc {
	return func(c *gin.Context) {
		logAsPrettyJSON("Headers", c.Request.Header)
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			log.Printf("JSON bind error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Invalid JSON payload",
				"detail": err.Error(),
			})
			return
		}
		logAsPrettyJSON("Parsed Payload", data)
		c.JSON(http.StatusCreated, gin.H{
			"message": "Raw data received successfully",
			"data":    data,
		})
	}
}

func HandleTypeMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON payload",
				"data":  err.Error(),
			})
			return
		}
		name := c.Query("name")
		if name == "" {
			name = "GeneratedType"
		}
		output, err := builder.GenerateGoTypeMap(name, payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate type map",
				"data":  err.Error(),
			})
			return
		}

		clean := strings.ReplaceAll(output, `\n`, "\n")
		clean = strings.ReplaceAll(clean, `\"`, `"`)
		clean = strings.TrimPrefix(clean, "```go")
		clean = strings.TrimPrefix(clean, "```")
		clean = strings.TrimSuffix(clean, "```")
		clean = strings.TrimSpace(clean)

		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(clean))
	}
}
