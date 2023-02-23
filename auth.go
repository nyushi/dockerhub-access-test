package dockerhub

import (
	"context"
	"encoding/json"
	"fmt"
)

type PostUsersLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type PostUsersLoginResponse struct {
	Token string `json:"token"`
}
type PostUsersLoginErrorResponse struct {
	Detail        string `json:"detail"`
	Login2FAToken string `json:"login_2fa_token"`
}

func (p *PostUsersLoginErrorResponse) Error() string {
	b, _ := json.MarshalIndent(p, "", "  ")
	return string(b)
}

func (c *Client) UsersLogin(ctx context.Context, user, pass string) error {
	url := fmt.Sprintf("%s/v2/users/login", c.baseURL)
	resp, errResp, err := callAPI[PostUsersLoginRequest, PostUsersLoginResponse, PostUsersLoginErrorResponse](c.httpClient, ctx, nil, url, &PostUsersLoginRequest{Username: user, Password: pass})
	if err != nil {
		return err
	}
	if errResp != nil {
		return err
	}
	c.token = resp.Token
	return nil
}
