package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/exp/slices"
)

var token string

type Due struct {
	Date      string `json:"date"`
	String    string `json:"string"`
	Lang      string `json:"lang"`
	Recurring bool   `json:"recurring"`
}

type Task struct {
	Id           int `json:"id"`
	Assigner     int `json:"assigner"`
	ProjectId    int `json:"project_id"`
	Project      Project
	SectionId    int    `json:"section_id"`
	Order        int    `json:"order"`
	Content      string `json:"content"`
	Description  string `json:"description"`
	Completed    bool   `json:"completed"`
	LabelIds     []int  `json:"label_ids"`
	Labels       []Label
	Priority     int    `json:"priority"`
	CommentCount int    `json:"comment_count"`
	Creater      int    `json:"creator"`
	Created      string `json:"created"`
	Due          Due    `json:"due"`
	Url          string `json:"url"`
}

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Label struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Favorite bool   `json:"favorite"`
}

func allProjects() []Project {
	projectUrl := "https://api.todoist.com/rest/v1/projects"

	req, _ := http.NewRequest("GET", projectUrl, nil)
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)

	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("ERROR: http request error.")
		os.Exit(1)
	}

	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: fetch data read error.")
		os.Exit(1)
	}

	var projects []Project

	json.Unmarshal(byteArray, &projects)

	return projects
}

func allLabels() []Label {
	labelsUrl := "https://api.todoist.com/rest/v1/labels"

	req, _ := http.NewRequest("GET", labelsUrl, nil)
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)

	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("ERROR: http request error.")
		os.Exit(1)
	}

	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: fetch data read error.")
		os.Exit(1)
	}

	var lebels []Label

	json.Unmarshal(byteArray, &lebels)

	return lebels
}

func main() {
	if token = os.Getenv("todoist_api_token"); token == "" {
		fmt.Println("ERROR: API token must be set in an environment variable (key: todoist_api_token)")
		os.Exit(1)
	}

	projects := allProjects()
	labels := allLabels()

	tasksUrl := "https://api.todoist.com/rest/v1/tasks"
	req, _ := http.NewRequest("GET", tasksUrl, nil)
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)

	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("ERROR: http request error.")
		os.Exit(1)
	}

	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: fetch data read error.")
		os.Exit(1)
	}

	var tasks []Task

	json.Unmarshal(byteArray, &tasks)

	var formattedTasks []Task
	for _, task := range tasks {
		projectIndex := slices.IndexFunc(projects, func(project Project) bool { return project.Id == task.ProjectId })
		task.Project = projects[projectIndex]

		var taskLabels []Label
		for _, id := range task.LabelIds {
			labelIndex := slices.IndexFunc(labels, func(label Label) bool { return label.Id == id })
			taskLabels = append(taskLabels, labels[labelIndex])
		}

		task.Labels = taskLabels

		formattedTasks = append(formattedTasks, task)
	}

	file, err := os.Create("tasks.csv")
	if err != nil {
		fmt.Println("ERROR: create file error")
		os.Exit(1)
	}
	defer file.Close()

	w := csv.NewWriter(file)

	w.Write([]string{"Id", "Content", "Description", "Due date", "Completed", "Priority", "Project name", "Labels", "Created"})

	for _, task := range formattedTasks {
		labelNames := ""
		for i, label := range task.Labels {
			if i == 0 {
				labelNames = label.Name
			} else {
				labelNames = labelNames + ", " + label.Name
			}
		}

		w.Write([]string{
			strconv.Itoa(task.Id),
			task.Content,
			task.Description,
			task.Due.Date,
			strconv.FormatBool(task.Completed),
			strconv.Itoa(task.Priority),
			task.Project.Name,
			labelNames,
			task.Created,
		})
	}

	if err := w.Error(); err != nil {
		fmt.Println("ERROR: write csv error")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("SUCESS: All your active tasks have been written to a csv file")
}
