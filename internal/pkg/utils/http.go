package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetClientRequestBody(tID string, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Print(err.Error(),
			"traceID", tID,
		)
		return []byte{}, err
	}

	if len(body) == 0 {
		error_message := "not enough parameters\n"
		log.Print(error_message,
			"traceID", tID,
		)
		return []byte{}, fmt.Errorf("%s", error_message)
	}

	return body, nil
}
