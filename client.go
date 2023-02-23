package dockerhub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	baseURL    string
	httpClient *http.Client

	token string
}

func NewClient(baseurl string) *Client {
	baseurl = strings.TrimSuffix(baseurl, "/")
	c := &Client{
		baseURL:    baseurl,
		httpClient: http.DefaultClient,
	}
	return c
}

type APIRequest interface {
	PostUsersLoginRequest | GetImageSummaryRequest | GetImageDetailsRequest | GetImageTagsRequest | ListRepositoryRequest
}
type APIResponse interface {
	PostUsersLoginResponse | GetImageSummaryResponse | GetImageDetailsResponse | GetImageTagsResponse | ListRepositoryResponse
}

type APIErrorResponse interface {
	PostUsersLoginErrorResponse | GetImageSummaryErrorResponse | GetImageDetailsErrorResponse | GetImageTagsErrorResponse | ListRepositoryErrorResponse
}

func callAPI[REQ APIRequest, RESP APIResponse, ERRRESP APIErrorResponse](c *http.Client, ctx context.Context, token *string, url string, reqData *REQ) (*RESP, *ERRRESP, error) {
	var reqBytes io.Reader
	if reqData != nil {
		b, _ := json.Marshal(reqData)
		reqBytes = bytes.NewBuffer(b)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, reqBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %s", err)
	}
	if reqData != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error at http request: %w", err)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error at reading body: %w", err)
	}

	var respData *RESP
	var errRespData *ERRRESP
	if res.StatusCode == http.StatusOK {
		respData = new(RESP)
		if err := json.Unmarshal(b, respData); err != nil {
			return nil, nil, fmt.Errorf("failed to parse json: %w, data=%s", err, b)
		}
	} else {
		errRespData = new(ERRRESP)
		if err := json.Unmarshal(b, errRespData); err != nil {
			return nil, nil, fmt.Errorf("failed to parse json: %w, data=%s", err, b)
		}
	}
	return respData, errRespData, nil
}
