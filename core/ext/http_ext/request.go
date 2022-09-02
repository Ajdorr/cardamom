package http_ext

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Post[Response any](url string, body any, headers map[string]string) (*Response, error) {

	client := &http.Client{}
	// rsp := new(Response)
	if data, err := json.Marshal(body); err != nil {
		return nil, err
	} else if req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data)); err != nil {
		return nil, err
	} else {
		for k, v := range headers {
			req.Header.Add(k, v)
		}

		if _, err := client.Do(req); err != nil {
			return nil, err
			// } else {
			// rsp.Body
		}
	}

	return nil, nil
}
