package job

import "encoding/base64"

type client struct {
	BaseURL      string
	APIKeyBase64 string
}

func NewClient(apikey string) *client {
	return &client{
		BaseURL:      "https://taucoder.com/api/v1",
		APIKeyBase64: base64.StdEncoding.EncodeToString([]byte(apikey)),
	}
}
