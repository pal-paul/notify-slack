package main

import (
	"encoding/json"
	"fmt"
	"log"

	env "github.com/pal-paul/go-libraries/pkg/env"
	http "github.com/pal-paul/go-libraries/pkg/http-client"
	slack "github.com/pal-paul/go-libraries/pkg/slack"
)

type Environment struct {
	GitHub struct {
		Token    string `env:"INPUT_TOKEN,required=true"`
		Api      string `env:"GITHUB_API_URL,required=true"`
		Repo     string `env:"GITHUB_REPOSITORY,required=true"`
		Workflow string `env:"GITHUB_WORKFLOW,required=true"`
		Branch   string `env:"GITHUB_REF,required=true"`
		Commit   string `env:"GITHUB_SHA,required=true"`
		RunId    string `env:"GITHUB_RUN_ID,required=true"`
		JobName  string `env:"GITHUB_JOB,required=true"`
		Server   string `env:"GITHUB_SERVER_URL,required=true"`
	}
	Input struct {
		Status     string `env:"INPUT_STATUS,required=true"`
		NotifyWhen string `env:"INPUT_NOTIFY_WHEN,required=true"`
	}
	Slack struct {
		Token   string `env:"INPUT_SLACK_TOKEN,required=true"`
		Channel string `env:"INPUT_SLACK_CHANNEL,required=true"`
	}
}

var envVar Environment

var (
	slackClient slack.ISlack
	httpClient  http.IHttpClient
)

// Initializing environment variables
func init() {
	_, err := env.Unmarshal(&envVar)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	slackClient = slack.New(
		slack.WithToken(envVar.Slack.Token),
	)
	httpClient = http.New()
}

func main() {
	if envVar.Input.Status == envVar.Input.NotifyWhen || envVar.Input.NotifyWhen == "always" {
		message := SlackMessageBuilder()
		_, err := slackClient.AddFormattedMessage(envVar.Slack.Channel, message)
		if err != nil {
			log.Fatalf("error while sending message to slack: %v", err)
		}
	}
}

type GitHubWorkflowResponse struct {
	Workflows []struct {
		Id      int    `json:"id"`
		Name    string `json:"name"`
		State   string `json:"state"`
		HtmlUrl string `json:"html_url"`
	} `json:"workflows"`
}

func GithubWorkflowUrl() string {
	url := fmt.Sprintf("%s/repos/%s/actions/workflows", envVar.GitHub.Api, envVar.GitHub.Repo)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", envVar.GitHub.Token),
		"Accept":        "application/vnd.github.v3+json",
	}
	bytes, code, err := httpClient.Get(url, headers)
	if err != nil {
		log.Fatalf("Error while getting workflow url: %v", err)
	}
	if code != 200 {
		log.Fatalf("Error while getting workflow url: %s", string(bytes))
	}
	var response GitHubWorkflowResponse
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		log.Fatalf("Error while unmarshalling workflow url: %v", err)
	}
	for _, workflow := range response.Workflows {
		if workflow.Name == envVar.GitHub.Workflow {
			return workflow.HtmlUrl
		}
	}
	return ""
}

func SlackMessageBuilder() slack.Message {
	message := slack.Message{
		Channel: envVar.Slack.Channel,
	}
	var status string = ""
	var icon string = ""
	switch envVar.Input.Status {
	case "success":
		icon = " :done-check: "
		status = " succeeded "
	case "failure":
		icon = " :bangbang: "
		status = " failed "
	default:
		icon = " :heavy_exclamation_mark: "
	}
	message.Blocks = append(message.Blocks, slack.Block{
		Type: slack.HeaderBlock,
		Text: &slack.Text{
			Type: slack.PlainText,
			Text: fmt.Sprintf("%s%s", envVar.GitHub.Workflow, status),
		},
	})

	commitUrl := fmt.Sprintf("%s/%s/commit/%s", envVar.GitHub.Server, envVar.GitHub.Repo, envVar.GitHub.Commit)
	runUrl := fmt.Sprintf("%s/%s/actions/runs/%s", envVar.GitHub.Server, envVar.GitHub.Repo, envVar.GitHub.RunId)
	message.Blocks = append(message.Blocks, slack.Block{
		Type: slack.SectionBlock,
		Text: &slack.Text{
			Type: slack.Mrkdwn,
			Text: fmt.Sprintf("%s %s %s in %s on <%s|%s> ", icon, envVar.GitHub.Workflow, status, GithubWorkflowUrl(), commitUrl, envVar.GitHub.Commit),
		},
	})
	message.Blocks = append(message.Blocks, slack.Block{
		Type: slack.SectionBlock,
		Text: &slack.Text{
			Type: slack.Mrkdwn,
			Text: fmt.Sprintf("<%s|View Run>", runUrl),
		},
	})
	return message
}
