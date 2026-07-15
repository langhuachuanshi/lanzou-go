package lanzou

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Login 登录蓝奏云帐号
func (c *Client) Login(user, pwd string) error {
	data := map[string]string{
		"action": "login",
		"task":   "login",
		"user":   user,
		"pass":   pwd,
		"setdir": "0",
	}
	body, _, err := c.post(baseURLPC+pathLogin, data, nil)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}

	var resp struct {
		Zt   int    `json:"zt"`
		Info string `json:"info"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("%w: invalid login response", ErrAPIError)
	}

	if resp.Zt != 1 {
		if strings.Contains(resp.Info, "密码") || strings.Contains(resp.Text, "密码") {
			return ErrPasswordWrong
		}
		return fmt.Errorf("%w: login failed: %s", ErrAPIError, resp.Info)
	}

	c.logged = true
	c.initUID()
	return nil
}

// Logout 登出
func (c *Client) Logout() error {
	_, _, err := c.get(baseURLPC+pathLogout, nil)
	if err != nil {
		return fmt.Errorf("logout request failed: %w", err)
	}
	c.logged = false
	c.cookies = nil
	return nil
}

// GetUserInfo 获取用户信息
func (c *Client) GetUserInfo() (*UserInfo, error) {
	if !c.isLoggedIn() {
		return nil, ErrNotLoggedIn
	}
	c.initUID()
	return &UserInfo{
		UserName: c.uid,
	}, nil
}

// GetAccountInfo 获取帐号详细信息
func (c *Client) GetAccountInfo() (*AccountInfo, error) {
	if !c.isLoggedIn() {
		return nil, ErrNotLoggedIn
	}
	c.initUID()

	// 通过个人中心页面获取帐号信息
	body, _, err := c.get(baseURLPC+"/mydisk.php?item=profile&action=mypower", nil)
	if err != nil {
		return nil, fmt.Errorf("get account info failed: %w", err)
	}
	html := string(body)

	// 从页面提取用户名
	reName := regexp.MustCompile(`(\d{11,})`)
	info := &AccountInfo{
		UserInfo: UserInfo{
			UserName: c.uid,
		},
	}
	if m := reName.FindStringSubmatch(html); len(m) > 1 {
		info.UserName = m[1]
	}

	// 提取容量信息
	reSize := regexp.MustCompile(`(\d+\.?\d*)\s*(GB|MB|KB|TB)`)
	matches := reSize.FindAllStringSubmatch(html, 2)
	if len(matches) >= 2 {
		info.TotalSize = matches[0][0]
		info.UsedSize = matches[1][0]
	} else if len(matches) >= 1 {
		info.TotalSize = matches[0][0]
	}

	return info, nil
}
