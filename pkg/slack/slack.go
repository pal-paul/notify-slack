package slack

//go:generate mockgen -source=slack.go -destination=mocks/mock-slack.go -package=mocks
import (
	"context"
	"net/http"
)

const (
	baseUrl = "https://slack.com/api"
)

type (
	SlackInterface interface {
		UploadFileWithContent(
			fileType string,
			fileName string,
			title string,
			content string,
			messageRef MessageRef,
		) (err error)

		AddFormattedMessage(
			channel string,
			message Message,
		) (messageRef MessageRef, err error)

		AddReaction(name string, item MessageRef) (err error)
		RemoveReaction(name string, item MessageRef) error
	}
	slack struct {
		httpClient *http.Client
		token      string
		ctx        context.Context
	}
)

func New(token string) SlackInterface {
	return &slack{
		httpClient: &http.Client{},
		token:      token,
		ctx:        context.Background(),
	}
}
