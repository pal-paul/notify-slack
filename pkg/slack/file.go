package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
)

func (api *slack) UploadFileWithContent(
	fileType string,
	fileName string,
	title string,
	content string,
	messageRef MessageRef,
) (err error) {
	values := url.Values{}
	if fileType != "" {
		values.Add("filetype", fileType)
	}
	if fileName != "" {
		values.Add("filename", fileName)
	}
	if title != "" {
		values.Add("title", title)
	}
	if messageRef.Timestamp != "" {
		values.Add("thread_ts", messageRef.Timestamp)
	}
	if messageRef.Channel != "" {
		values.Add("channels", strings.Join([]string{messageRef.Channel}, ","))
	}
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	var response SlackResponse
	if content != "" {
		values.Add("content", content)
		if resp, err := api.postForm(fmt.Sprintf("%s/reactions.remove", baseUrl), headers, values); err != nil {
			return err
		} else {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(body, &response); err != nil {
				return err
			}
			if !response.Ok {
				return fmt.Errorf("error slack response")
			}
		}
	}
	return nil
}
