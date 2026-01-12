package github

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var auth string

func GetGit() {
	token := os.Getenv("GH_TOKEN")
	auth = fmt.Sprintf("Bearer %s", token)
}

func GetUserRepos(username string) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Github-Api-Version", "2022-11-28")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		log.Println("Rate Limit Reached")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%s", string(body))
}
