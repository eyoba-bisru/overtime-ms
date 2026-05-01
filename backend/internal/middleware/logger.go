package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		// -------------------------
		// Read Request Body
		// -------------------------
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// -------------------------
		// Capture Response Body
		// -------------------------
		bw := &bodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = bw

		// process request
		c.Next()

		// -------------------------
		// Prepare Log Data
		// -------------------------
		latency := time.Since(start)

		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		reqJSON := formatJSON(requestBody)
		resJSON := formatJSON(bw.body.Bytes())

		log.Println("========== API LOG ==========")
		log.Printf("Method   : %s\n", method)
		log.Printf("Path     : %s\n", path)
		log.Printf("Status   : %d\n", status)
		log.Printf("Latency  : %s\n", latency)
		log.Printf("Request  : %s\n", reqJSON)
		log.Printf("Response : %s\n", resJSON)
		log.Println("=============================")
	}
}

// sensitiveFields are field names whose values should be masked in logs.
var sensitiveFields = []string{"password", "new_password", "password_hash"}

func maskSensitiveFields(data []byte) []byte {
	if len(data) == 0 {
		return data
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return data
	}

	for key := range parsed {
		for _, sensitive := range sensitiveFields {
			if strings.EqualFold(key, sensitive) {
				parsed[key] = "[REDACTED]"
			}
		}
	}

	masked, err := json.Marshal(parsed)
	if err != nil {
		return data
	}
	return masked
}

func formatJSON(data []byte) string {
	if len(data) == 0 {
		return "{}"
	}

	masked := maskSensitiveFields(data)

	var pretty bytes.Buffer
	err := json.Indent(&pretty, masked, "", "  ")
	if err != nil {
		return string(masked)
	}

	return pretty.String()
}
