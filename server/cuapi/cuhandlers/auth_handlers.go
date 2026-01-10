package cuhandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
)

type ClickUpTokenResponse struct {
	AccessToken string `json:"access_token"`
}

var redirect = "http://localhost:8080/clickup/success"
var client_id = "4SIB37ZKY9WMNTU3M2ADMVX238K0I2IS"
var client_secret = "1NV4A2B2V3TJKWZ36AXNYGNAAR6VBN4PRGQ6LP6CHC1HGEZXL5X9MS0EJ14GN4XP"
var oauth_url = fmt.Sprintf("https://app.clickup.com/api?client_id=%s&redirect_uri=%s", client_id, redirect)
var oAuthCode string
var authToken string

func AuthClickUp() error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", oauth_url}
	case "darwin": // macOS
		cmd = "open"
		args = []string{oauth_url}
	case "linux":
		cmd = "xdg-open"
		args = []string{oauth_url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}

func HandleOAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No authorization code provided",
			})
			return
		}

		tokenURL := "https://api.clickup.com/api/v2/oauth/token"
		params := url.Values{}
		params.Add("client_id", client_id)
		params.Add("client_secret", client_secret)
		params.Add("code", code)

		resp, err := http.PostForm(tokenURL, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to auth ClickUp"})
			return
		}
		defer resp.Body.Close()

		var tokenData ClickUpTokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Authorization successful, you can close this window.",
			"token":   tokenData.AccessToken,
		})
	}
}
