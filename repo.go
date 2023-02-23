package dockerhub

import (
	"context"
	"encoding/json"
	"fmt"
)

type LayerResponse struct {
	Digest      string `json:"digest"`
	Size        int    `json:"size"`
	Instruction string `json:"instruction"`
}

type ImageResponse struct {
	Architecture string           `json:"architecture"`
	Features     string           `json:"features"`
	Variant      string           `json:"variant"`
	Digest       *string          `json:"digest"`
	Layers       []*LayerResponse `json:"layers"`
	OS           string           `json:"os"`
	OSFeatures   string           `json:"os_features"`
	OSVersion    string           `json:"os_version"`
	Size         int              `json:"size"`
	Status       string           `json:"Status"`
	LastPulled   string           `json:"last_pulled"`
	LastPushed   string           `json:"last_pushed"`
}

type TagResponse struct {
	ID                  int              `json:"id"`
	Images              []*ImageResponse `json:"images"`
	Cerator             int              `json:"creator"`
	LastUpdated         *string          `json:"last_updated"`
	LastUpdater         int              `json:"last_updater"`
	LastUpdaterUsername string           `json:"last_updater_username"`
	Name                string           `json:"name"`
	Repository          int              `json:"repository"`
	FullSize            int              `json:"full_size"`
	V2                  bool             `json:"v2"`
	Status              string           `json:"status"`
	TagLastPulled       *string          `json:"tag_last_pulled"`
	TagLastPushed       *string          `json:"tag_last_pushed"`
}

type ListRepositoryRequest struct{}
type ListRepositoryResponse struct {
	c        *Client
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []*TagResponse `json:"results"`
}

func (g *ListRepositoryResponse) GetNext(ctx context.Context) (*ListRepositoryResponse, error) {
	if g.Next == nil {
		return nil, nil
	}
	return g.c.listRepository(ctx, *g.Next)
}

func (g *ListRepositoryResponse) GetPrevious(ctx context.Context) (*ListRepositoryResponse, error) {
	if g.Previous == nil {
		return nil, nil
	}
	return g.c.listRepository(ctx, *g.Previous)
}

type ListRepositoryErrorResponse struct {
	Detail  string `json:"detail"`
	Message string `json:"message"`
}

func (g *ListRepositoryErrorResponse) Error() string {
	b, _ := json.MarshalIndent(g, "", "  ")
	return string(b)
}

func (c *Client) listRepository(ctx context.Context, url string) (*ListRepositoryResponse, error) {
	resp, errResp, err := callAPI[ListRepositoryRequest, ListRepositoryResponse, ListRepositoryErrorResponse](c.httpClient, ctx, &c.token, url, &ListRepositoryRequest{})
	if err != nil {
		return nil, err
	}
	if errResp != nil {
		return nil, errResp
	}
	resp.c = c
	return resp, nil

}
func (c *Client) ListRepository(ctx context.Context, namespace, repository string) (*ListRepositoryResponse, error) {
	url := fmt.Sprintf("%s/v2/namespaces/%s/repositories/%s/tags?page_size=50", c.baseURL, namespace, repository)
	return c.listRepository(ctx, url)
}
