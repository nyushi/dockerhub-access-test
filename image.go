package dockerhub

import (
	"context"
	"encoding/json"
	"fmt"
)

type GetImageSummaryRequest struct{}
type GetImageSummaryResponse struct {
	ActiveFrom string `json:"active_from"`
	Statistics struct {
		Total    int `json:"total"`
		Active   int `json:"active"`
		Inactive int `json:"inactive"`
	} `json:"statistics"`
}
type GetImageSummaryErrorResponse struct {
	TxnID   string `json:"txnid"`
	Message string `json:"message"`
	ErrInfo struct {
		APICallDockerID string `json:"api_call_docker_id"`
		APICallName     string `json:"api_call_name"`
		APICallStart    string `json:"api_call_start"`
		APICallTxnID    string `json:"api_call_txnid"`
	} `json:"errinfo"`
}

func (g *GetImageSummaryErrorResponse) Error() string {
	b, _ := json.MarshalIndent(g, "", "  ")
	return string(b)
}

func (c *Client) GetImageSummary(ctx context.Context, namespace, repository string) (*GetImageSummaryResponse, error) {
	url := fmt.Sprintf("%s/v2/namespaces/%s/repositories/%s/images-summary", c.baseURL, namespace, repository)
	resp, errResp, err := callAPI[GetImageSummaryRequest, GetImageSummaryResponse, GetImageSummaryErrorResponse](c.httpClient, ctx, &c.token, url, &GetImageSummaryRequest{})
	if err != nil {
		return nil, err
	}
	if errResp != nil {
		return nil, errResp
	}
	return resp, nil
}

type GetImageDetailsRequest struct{}
type GetImageDetailsResponse struct {
	c        *Client
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Namespace  string `json:"namespace"`
		Repository string `json:"repository"`
		Digest     string `json:"digest"`
		Tags       []struct {
			Tag       string `json:"tag"`
			IsCurrent bool   `json:"is_current"`
		} `json:"tags"`
		LastPushed *string `json:"last_pushed"`
		LastPulled *string `json:"last_pulled"`
		Status     string  `json:"status"`
	} `json:"results"`
}

func (g *GetImageDetailsResponse) GetNext(ctx context.Context) (*GetImageDetailsResponse, error) {
	if g.Next == nil {
		return nil, nil
	}
	return g.c.getImageDetails(ctx, *g.Next)
}

func (g *GetImageDetailsResponse) GetPrevious(ctx context.Context) (*GetImageDetailsResponse, error) {
	if g.Previous == nil {
		return nil, nil
	}
	return g.c.getImageDetails(ctx, *g.Previous)
}

type GetImageDetailsErrorResponse struct {
	TxnID   string `json:"txnid"`
	Message string `json:"message"`
	ErrInfo struct {
		APICallDockerID string `json:"api_call_docker_id"`
		APICallName     string `json:"api_call_name"`
		APICallStart    string `json:"api_call_start"`
		APICallTxnID    string `json:"api_call_txnid"`
	} `json:"errinfo"`
}

func (g *GetImageDetailsErrorResponse) Error() string {
	b, _ := json.MarshalIndent(g, "", "  ")
	return string(b)
}

func (c *Client) getImageDetails(ctx context.Context, url string) (*GetImageDetailsResponse, error) {
	resp, errResp, err := callAPI[GetImageDetailsRequest, GetImageDetailsResponse, GetImageDetailsErrorResponse](c.httpClient, ctx, &c.token, url, &GetImageDetailsRequest{})
	if err != nil {
		return nil, err
	}
	if resp != nil {
		resp.c = c
		return resp, nil
	}

	return nil, errResp
}

func (c *Client) GetImageDetails(ctx context.Context, namespace, repository string) (*GetImageDetailsResponse, error) {
	url := fmt.Sprintf("%s/v2/namespaces/%s/repositories/%s/images?ordering=-last_activity", c.baseURL, namespace, repository)
	return c.getImageDetails(ctx, url)
}

type GetImageTagsRequest struct{}
type GetImageTagsResponse struct {
	c        *Client
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Tag       string `json:"tag"`
		IsCurrent bool   `json:"is_current"`
	} `json:"results"`
}
type GetImageTagsErrorResponse struct {
	TxnID   string `json:"txnid"`
	Message string `json:"message"`
	ErrInfo struct {
		APICallDockerID string `json:"api_call_docker_id"`
		APICallName     string `json:"api_call_name"`
		APICallStart    string `json:"api_call_start"`
		APICallTxnID    string `json:"api_call_txnid"`
	} `json:"errinfo"`
}

func (g *GetImageTagsErrorResponse) Error() string {
	b, _ := json.MarshalIndent(g, "", "  ")
	return string(b)
}

func (c *Client) getImageTags(ctx context.Context, url string) (*GetImageTagsResponse, error) {
	resp, errResp, err := callAPI[GetImageTagsRequest, GetImageTagsResponse, GetImageTagsErrorResponse](c.httpClient, ctx, &c.token, url, &GetImageTagsRequest{})
	if err != nil {
		return nil, err
	}
	if resp != nil {
		resp.c = c
		return resp, nil
	}

	return nil, errResp
}

func (c *Client) GetImageTags(ctx context.Context, namespace, repository, digest string) (*GetImageTagsResponse, error) {
	url := fmt.Sprintf("%s/v2/namespaces/%s/repositories/%s/images/%s/tags", c.baseURL, namespace, repository, digest)
	return c.getImageTags(ctx, url)
}
