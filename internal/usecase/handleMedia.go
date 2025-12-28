package usecase

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

func (u *UseCase) HandleMedia(file io.ReadCloser, mimeType string) (string, string, error) {
	const op = "usecase.HandleMedia"

	data, err := io.ReadAll(file)
	if err != nil {
		u.log.Error("Invalid data", "op", op, "error", err)
		return "", "", err
	}

	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}

	base64Data := base64.StdEncoding.EncodeToString(data)
	dataurl := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)

	return dataurl, mimeType, nil
}
