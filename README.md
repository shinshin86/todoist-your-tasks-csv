# todoist-your-tasks-csv

This is a CLI tool to write all active tasks on todoist to a CSV file.

## Install

```sh
go install github.com/shinshin86/todoist-your-tasks-csv@latest
```

## Usage

After execution, a file named `tasks.csv` will be generated in the current directory.

```sh
todoist_api_token=<Your Todoist API token> todoist-your-tasks-csv
```