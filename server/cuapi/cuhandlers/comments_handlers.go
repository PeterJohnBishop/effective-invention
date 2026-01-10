package cuhandlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetTaskCommentsReq struct {
	TaskId  string
	Start   *int
	StartId string
}

type GetTaskCommentsResp struct {
	Comments []Comment
}

type Comment struct {
	Id         string
	Comment    []map[string]interface{}
	User       User
	Resolved   bool
	Assignee   User
	AssignedBy User
	Reactions  []string
	Date       int
	ReplyCount int
}

type User struct {
	Id             string
	Username       string
	Initials       string
	Email          string
	Color          string
	ProfilePicture string
}

func GetTaskComments(request *GetTaskCommentsReq) ([]byte, error) {

	baseUrl := fmt.Sprintf("https://api.clickup.com/api/v2/task/%s/comment", request.TaskId)
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	if request.Start != nil {
		q.Set("start", strconv.Itoa(*request.Start))
	}
	if request.StartId != "" {
		q.Set("start_id", request.StartId)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	return body, nil
}

func HandleTaskComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		//
	}
}
