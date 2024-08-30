package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createFormFile(w *multipart.Writer, fieldname, filename string, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}

func (c *client) NewJobCreateRequest(quality int, filenames []string, contentTypes []string) (*http.Request, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for i := range filenames {
		imgBytes, err := os.ReadFile(filenames[i])
		if err != nil {
			return nil, err
		}

		fw, err := createFormFile(w, "image", filepath.Base(filenames[i]), contentTypes[i])
		if err != nil {
			return nil, err
		}

		_, err = fw.Write(imgBytes)
		if err != nil {
			return nil, err
		}
	}

	options := struct {
		Quality        int    `json:"quality"`
		EncoderVersion string `json:"encoder_version"`
	}{
		Quality:        quality,
		EncoderVersion: "latest",
	}

	optionsBytes, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	err = w.WriteField("options", string(optionsBytes))
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/job-create", &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Basic "+c.APIKeyBase64)

	return req, nil
}
