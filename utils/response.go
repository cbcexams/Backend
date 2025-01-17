package utils

import (
	"encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

// Response is a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendResponse sends a standardized JSON response
func SendResponse(c *beego.Controller, success bool, message string, data interface{}, err error) {
	fmt.Println("=== Sending Response ===")

	response := Response{
		Success: success,
		Message: message,
		Data:    data,
	}

	if err != nil {
		response.Error = err.Error()
	}

	// Debug print the response
	jsonData, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("Response JSON:\n%s\n", string(jsonData))

	c.Data["json"] = response
	c.ServeJSON()
}
