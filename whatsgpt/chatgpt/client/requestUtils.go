package client

import (
	"net/http"
)

func AddRequiredHeaders(httpRequest *http.Request) {
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer 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")
}
