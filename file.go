package lanzou

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetFileList 获取文件夹内的文件列表
// fid: 文件夹ID，根目录传 0
func (c *Client) GetFileList(fid int) (*FileList, error) {
	if !c.isLoggedIn() {
		return nil, ErrNotLoggedIn
	}
	c.initUID()
	data := map[string]string{
		"task":      "5",
		"folder_id": fmt.Sprintf("%d", fid),
		"pg":        "1",
		"vei":       c.vei,
	}
	body, _, err := c.post(c.apiURL(pathTaskAPI), data, nil)
	if err != nil {
		return nil, fmt.Errorf("get file list failed: %w", err)
	}

	var resp FileList
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("%w: invalid file list response", ErrAPIError)
	}
	if resp.Zt != 1 && resp.Zt != 2 {
		return nil, fmt.Errorf("%w: zt=%d info=%s", ErrAPIError, resp.Zt, string(resp.Info))
	}
	return &resp, nil
}

// GetFileInfoByURL 通过分享链接获取文件信息（委托给 GetFileInfo）
func (c *Client) GetFileInfoByURL(shareURL, pwd string) (*FileDetail, error) {
	return c.GetFileInfo(shareURL, pwd)
}

// MoveFiles 移动文件到指定文件夹
// fids: 文件ID列表，fid: 目标文件夹ID
func (c *Client) MoveFiles(fids []string, fid int) error {
	if !c.isLoggedIn() {
		return ErrNotLoggedIn
	}
	c.initUID()
	data := map[string]string{
		"task":      "17",
		"folder_id": fmt.Sprintf("%d", fid),
		"file_id":   strings.Join(fids, "-"),
	}
	body, _, err := c.post(c.apiURL(pathTaskAPI), data, nil)
	if err != nil {
		return fmt.Errorf("move files failed: %w", err)
	}
	return c.checkRespBody(body)
}

// DeleteFiles 删除文件
// fids: 文件ID列表
func (c *Client) DeleteFiles(fids []string) error {
	if !c.isLoggedIn() {
		return ErrNotLoggedIn
	}
	c.initUID()
	data := map[string]string{
		"task":    "6",
		"file_id": strings.Join(fids, "-"),
	}
	body, _, err := c.post(c.apiURL(pathTaskAPI), data, nil)
	if err != nil {
		return fmt.Errorf("delete files failed: %w", err)
	}
	return c.checkRespBody(body)
}

// SetPassword 设置文件或文件夹密码
// entityID: 文件或文件夹ID
func (c *Client) SetPassword(entityID int, pwd string) error {
	if !c.isLoggedIn() {
		return ErrNotLoggedIn
	}
	c.initUID()
	data := map[string]string{
		"task":      "5",
		"file_id":   fmt.Sprintf("%d", entityID),
		"shows":     "2",
		"shownames": pwd,
	}
	body, _, err := c.post(c.apiURL(pathTaskAPI), data, nil)
	if err != nil {
		return fmt.Errorf("set password failed: %w", err)
	}
	return c.checkRespBody(body)
}

// checkRespBody 检查通用响应体
func (c *Client) checkRespBody(body []byte) error {
	var resp struct {
		Zt   int    `json:"zt"`
		Info string `json:"info"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("%w: invalid response", ErrAPIError)
	}
	if resp.Zt == 0 {
		return fmt.Errorf("%w: %s", ErrAPIError, resp.Info)
	}
	return nil
}
