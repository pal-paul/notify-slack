package slack

import (
	"encoding/json"
	"fmt"
	"io"
)

func (m *slack) AddFormattedMessage(
	channel string,
	message Message,
) (messageRef MessageRef, err error) {
	message.Channel = channel
	var response SlackResponse

	apiEndpoint := fmt.Sprintf("%s/chat.postMessage", baseUrl)
	header := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
	reqBody, err := json.Marshal(message)
	if err != nil {
		return messageRef, err
	}
	resp, err := m.postRequest(apiEndpoint, header, reqBody)
	if err != nil {
		return messageRef, fmt.Errorf("error post to slack: %v", err)
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return messageRef, err
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return messageRef, err
		}
		if !response.Ok {
			return messageRef, fmt.Errorf("error slack response")
		}
		messageRef.Channel = response.Channel
		messageRef.Timestamp = response.Ts
	}
	return messageRef, nil
}
