package lanzou

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// UploadFile 上传本地文件到蓝奏云
// filePath: 本地文件路径, fid: 目标文件夹ID（根目录传 0）, desc: 文件描述（可选）
func (c *Client) UploadFile(filePath string, fid int, desc ...string) (*UploadResult, error) {
	if !c.isLoggedIn() {
		return nil, ErrNotLoggedIn
	}

	// 检查文件大小
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("stat file failed: %w", err)
	}
	if info.Size() > int64(c.maxsize) {
		return nil, fmt.Errorf("%w: file size %d exceeds limit %d", ErrFileSizeLimit, info.Size(), c.maxsize)
	}

	// 打开文件
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file failed: %w", err)
	}
	defer f.Close()

	// 计算文件MD5
	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return nil, fmt.Errorf("calc md5 failed: %w", err)
	}
	md5Str := hex.EncodeToString(hash.Sum(nil))
	f.Seek(0, 0)

	// 获取上传参数
	uploadInfo, err := c.getUploadParams(fid)
	if err != nil {
		return nil, fmt.Errorf("get upload params failed: %w", err)
	}

	// 上传延迟
	if c.uploadDelay[1] > c.uploadDelay[0] {
		delay := c.uploadDelay[0] + rand.Intn(c.uploadDelay[1]-c.uploadDelay[0])
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	// 构造上传字段
	descStr := ""
	if len(desc) > 0 {
		descStr = desc[0]
	}
	fields := map[string]string{
		"task":          "1",
		"folder_id":     fmt.Sprintf("%d", fid),
		"id":            "WU_FILE_0",
		"name":          filepath.Base(filePath),
		"upload_type":   "file",
		"t_a":           uploadInfo.TA,
		"t_b":           uploadInfo.TB,
		"t_c":           uploadInfo.TC,
		"ve":            "1",
		"des":           descStr,
		"ss":            md5Str,
	}

	// 上传
	body, _, err := c.postMultipart(
		baseURLUpload+pathUpload,
		fields,
		"upload_file",
		filepath.Base(filePath),
		f,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("upload request failed: %w", err)
	}

	var resp UploadResult
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("%w: invalid upload response", ErrAPIError)
	}
	if resp.Zt == 0 {
		return nil, fmt.Errorf("%w: %s", ErrUploadFailed, resp.Info)
	}
	return &resp, nil
}

// UploadFileByURL 上传网盘已有文件（非本地文件）
// fileURL: 文件URL, fid: 目标文件夹ID
func (c *Client) UploadFileByURL(fileURL string, fid int, desc ...string) (*UploadResult, error) {
	if !c.isLoggedIn() {
		return nil, ErrNotLoggedIn
	}

	descStr := ""
	if len(desc) > 0 {
		descStr = desc[0]
	}

	data := map[string]string{
		"task":      "42",
		"folder_id": fmt.Sprintf("%d", fid),
		"url":       fileURL,
		"name":      filepath.Base(fileURL),
		"des":       descStr,
	}
	body, _, err := c.post(baseURLPC+pathTaskAPI, data, nil)
	if err != nil {
		return nil, fmt.Errorf("upload by url failed: %w", err)
	}

	var resp UploadResult
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("%w: invalid upload response", ErrAPIError)
	}
	if resp.Zt == 0 {
		return nil, fmt.Errorf("%w: %s", ErrUploadFailed, resp.Info)
	}
	return &resp, nil
}

// uploadInfo 上传参数
type uploadInfo struct {
	TA string `json:"t_a"`
	TB string `json:"t_b"`
	TC string `json:"t_c"`
}

// getUploadParams 获取上传所需的动态参数
func (c *Client) getUploadParams(fid int) (*uploadInfo, error) {
	// 先获取上传页面，提取动态参数
	data := map[string]string{
		"task":      "1",
		"folder_id": fmt.Sprintf("%d", fid),
	}
	body, _, err := c.post(baseURLUpload+pathUpload, data, nil)
	if err != nil {
		return nil, err
	}

	// 尝试从响应中提取参数
	var resp struct {
		Zt   int `json:"zt"`
		Text struct {
			TA string `json:"t_a"`
			TB string `json:"t_b"`
			TC string `json:"t_c"`
		} `json:"text"`
	}
	if err := json.Unmarshal(body, &resp); err == nil && resp.Text.TA != "" {
		return &uploadInfo{
			TA: resp.Text.TA,
			TB: resp.Text.TB,
			TC: resp.Text.TC,
		}, nil
	}

	// 回退：使用默认参数
	return &uploadInfo{
		TA: fmt.Sprintf("%d", time.Now().UnixMilli()),
		TB: fmt.Sprintf("%d", time.Now().UnixMilli()+1000),
		TC: fmt.Sprintf("%d", time.Now().UnixMilli()+2000),
	}, nil
}
