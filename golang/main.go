package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

type Space struct {
	Id string `json:"id"`
}

type Watcher struct {
	Id             int    `json:"id"`
	Initials       string `json:"initials"`
	ProfilePicture any    `json:"profilePicture"`
	Username       string `json:"username"`
	Color          string `json:"color"`
	Email          string `json:"email"`
}

type Status struct {
	Orderindex int    `json:"orderindex"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	Color      string `json:"color"`
	Id         string `json:"id"`
}

type List struct {
	Access bool   `json:"access"`
	Id     string `json:"id"`
	Name   string `json:"name"`
}

type Task struct {
	Points          any           `json:"points"`
	DateDone        any           `json:"date_done"`
	GroupAssignees  []any         `json:"group_assignees"`
	Id              string        `json:"id"`
	Orderindex      string        `json:"orderindex"`
	Parent          any           `json:"parent"`
	Sharing         Sharing       `json:"sharing"`
	TextContent     string        `json:"text_content"`
	TopLevelParent  any           `json:"top_level_parent"`
	Assignees       []any         `json:"assignees"`
	Space           Space         `json:"space"`
	StartDate       any           `json:"start_date"`
	Tags            []any         `json:"tags"`
	TeamId          string        `json:"team_id"`
	Url             string        `json:"url"`
	Watchers        []Watcher     `json:"watchers"`
	CustomFields    []CustomField `json:"custom_fields"`
	CustomItemId    int           `json:"custom_item_id"`
	LinkedTasks     []any         `json:"linked_tasks"`
	Project         Project       `json:"project"`
	Status          Status        `json:"status"`
	Checklists      []any         `json:"checklists"`
	DateUpdated     string        `json:"date_updated"`
	Dependencies    []any         `json:"dependencies"`
	Description     string        `json:"description"`
	List            List          `json:"list"`
	PermissionLevel string        `json:"permission_level"`
	Priority        any           `json:"priority"`
	TimeSpent       int           `json:"time_spent"`
	Name            string        `json:"name"`
	TimeEstimate    any           `json:"time_estimate"`
	Attachments     []any         `json:"attachments"`
	Creator         Creator       `json:"creator"`
	DateCreated     string        `json:"date_created"`
	Locations       []any         `json:"locations"`
	Archived        bool          `json:"archived"`
	DateClosed      any           `json:"date_closed"`
	CustomId        any           `json:"custom_id"`
	DueDate         any           `json:"due_date"`
	Folder          Folder        `json:"folder"`
}

type Tasks struct {
	Task []Task `json:"tasks"`
}

type Sharing struct {
	Public               bool     `json:"public"`
	PublicFields         []string `json:"public_fields"`
	PublicShareExpiresOn any      `json:"public_share_expires_on"`
	SeoOptimized         bool     `json:"seo_optimized"`
	Token                any      `json:"token"`
}

type CustomField struct {
	Name           string     `json:"name"`
	Required       bool       `json:"required"`
	Type           string     `json:"type"`
	TypeConfig     TypeConfig `json:"type_config"`
	DateCreated    string     `json:"date_created"`
	HideFromGuests bool       `json:"hide_from_guests"`
	Id             string     `json:"id"`
}

type Project struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Access bool   `json:"access"`
	Hidden bool   `json:"hidden"`
}

type Creator struct {
	Color          string `json:"color"`
	Email          string `json:"email"`
	Id             int    `json:"id"`
	ProfilePicture any    `json:"profilePicture"`
	Username       string `json:"username"`
}

type Folder struct {
	Name   string `json:"name"`
	Access bool   `json:"access"`
	Hidden bool   `json:"hidden"`
	Id     string `json:"id"`
}

type TypeConfig struct{}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type Performance struct {
	Duration string
	RPM      string
	TPS      string
}

func calculatePerformance(totalTasks int, start time.Time) Performance {
	elapsed := time.Since(start)
	pages := float64(totalTasks) / 100.0
	if pages < 1 {
		pages = 1
	}

	rpm := (pages / elapsed.Minutes())
	tps := float64(totalTasks) / elapsed.Seconds()

	return Performance{
		Duration: elapsed.Round(time.Millisecond).String(),
		RPM:      fmt.Sprintf("%.2f", rpm),
		TPS:      fmt.Sprintf("%.2f", tps),
	}
}

func GetAllTasks(teamID string, token string) ([]Task, error) {
	var allTasks []Task
	var wg sync.WaitGroup

	limiter := rate.NewLimiter(rate.Every(time.Minute/1000), 1)

	taskChan := make(chan []Task, 100)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for tasks := range taskChan {
			allTasks = append(allTasks, tasks...)
		}
	}()

	fmt.Println("starting concurrent fetch...")
	start := time.Now()

	for page := 0; ; page++ {
		select {
		case <-ctx.Done():
			goto WaitAndFinish
		default:
		}

		if err := limiter.Wait(context.Background()); err != nil {
			break
		}

		wg.Add(1)
		go func(p int) {
			defer wg.Done()

			tasks, err := fetchPage(teamID, token, p)
			if err != nil || len(tasks) == 0 {
				cancel()
				return
			}

			taskChan <- tasks
		}(page)

	}

WaitAndFinish:
	wg.Wait()
	close(taskChan)

	time.Sleep(100 * time.Millisecond)

	performance := calculatePerformance(len(allTasks), start)
	fmt.Printf("fetched %d tasks in %s seconds. RPM: %s\n",
		len(allTasks), performance.Duration, performance.RPM)

	return allTasks, nil
}

func fetchPage(teamID string, token string, page int) ([]Task, error) {
	url := fmt.Sprintf("https://api.clickup.com/api/v2/team/%s/task?page=%d&include_closed=true&subtasks=true", teamID, page)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limit hit")
	}

	var tasksResponse TasksResponse
	if err := json.NewDecoder(resp.Body).Decode(&tasksResponse); err != nil {
		return nil, err
	}

	return tasksResponse.Tasks, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiToken := os.Getenv("TOKEN")
	GetAllTasks("36226098", apiToken)
}
