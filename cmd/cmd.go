package main

import (
	"fmt"
	"log"

	env "github.com/pal-paul/go-libraries/pkg/env"
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

var slackClient slack.ISlack

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
	return fmt.Sprintf("%s/%s/actions",
		envVar.GitHub.Server,
		envVar.GitHub.Repo)
}

func SlackMessageBuilder() slack.Message {
	message := slack.Message{
		Channel: envVar.Slack.Channel,
	}
	var status string = ""
	var icon string = ""
	switch envVar.Input.Status {
	case "success":
		icon = " :white_check_mark: "
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

	// Construct URLs for GitHub web links (not API URLs)
	commitUrl := fmt.Sprintf("%s/%s/commit/%s", envVar.GitHub.Server, envVar.GitHub.Repo, envVar.GitHub.Commit)
	runUrl := fmt.Sprintf("%s/%s/actions/runs/%s", envVar.GitHub.Server, envVar.GitHub.Repo, envVar.GitHub.RunId)

	// Add workflow status section
	message.Blocks = append(message.Blocks, slack.Block{
		Type: slack.SectionBlock,
		Text: &slack.Text{
			Type: slack.Mrkdwn,
			Text: fmt.Sprintf("%s *%s*%s", icon, envVar.GitHub.Workflow, status),
		},
	})

	// Add commit and run links section
	message.Blocks = append(message.Blocks, slack.Block{
		Type: slack.SectionBlock,
		Text: &slack.Text{
			Type: slack.Mrkdwn,
			Text: fmt.Sprintf("• Run: <%s|%s>\n• Commit: <%s|%s>",
				runUrl,
				"View workflow run",
				commitUrl,
				envVar.GitHub.Commit[:7]),
		},
	})
	return message
}
