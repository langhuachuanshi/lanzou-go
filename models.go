package lanzou

import "encoding/json"

// apiResp 蓝奏云 API 通用响应结构
type apiResp struct {
	Zt   int             `json:"zt"`
	Info json.RawMessage `json:"info"`
	Dom  string          `json:"dom"`
	Text json.RawMessage `json:"text"`
}

// FolderInfo 文件夹信息
type FolderInfo struct {
	FolID      string `json:"fol_id"`     // 文件夹ID
	Name       string `json:"name"`       // 文件夹名
	FolderDesc string `json:"folder_des"` // 描述
	Onof       string `json:"onof"`       // 密码开关
	FolderLock string `json:"folderlock"` // 锁定
	IsLock     string `json:"is_lock"`    // 状态
	IsCopyr    string `json:"is_copyright"` // 版权
}

// FileInfo 文件列表中的文件信息
type FileInfo struct {
	ID       string `json:"id"`        // 文件ID
	Name     string `json:"name"`      // 文件名
	NameAll  string `json:"name_all"`  // 完整文件名
	Size     string `json:"size"`      // 文件大小
	Icon     string `json:"icon"`      // 图标
	Time     string `json:"time"`      // 上传时间
	Downs    string `json:"downs"`     // 下载次数
	Onof     string `json:"onof"`      // 密码开关
	IsDes    string `json:"is_des"`    // 是否有描述
	FileLock string `json:"filelock"`  // 文件锁定
	IsNewd   string `json:"is_newd"`   // 新版标记
	FID      string `json:"f_id"`      // FID
}

// FileDetail 文件详情（含直链）
type FileDetail struct {
	FileID      string `json:"file_id"`
	NameAll     string `json:"name_all"`
	Size        string `json:"size"`
	UploadTime  string `json:"time"`
	DownloadURL string `json:"url"`
	DURL        string `json:"durl"`
	Description string `json:"des"`
	IsNewd      int    `json:"is_newd"`
	IsLock      int    `json:"is_locked"`
	IsFolder    int    `json:"is_folderd"`
}

// FolderList 文件夹列表响应
type FolderList struct {
	Zt   int             `json:"zt"`
	Info json.RawMessage `json:"info"` // 面包屑导航（array）
	Text []*FolderInfo   `json:"text"`
}

// FileList 文件列表响应
type FileList struct {
	Zt   int             `json:"zt"`
	Info json.RawMessage `json:"info"` // 可能 string 或 int
	Text []*FileInfo     `json:"text"`
}

// RecycleList 回收站列表响应
type RecycleList struct {
	Zt   int             `json:"zt"`
	Info string          `json:"info"`
	Text []*RecycledFile `json:"text"`
}

// RecycledFile 回收站中的文件
type RecycledFile struct {
	FileID     string `json:"id"`
	FileName   string `json:"name_all"`
	FilePath   string `json:"path"`
	UploadTime string `json:"time"`
}

// UploadResult 上传结果
type UploadResult struct {
	Zt       int    `json:"zt"`
	Text     string `json:"text"`
	Info     string `json:"info"`
	FileID   string `json:"file_id"`
	FileName string `json:"name_all"`
}

// UserInfo 用户信息
type UserInfo struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"nickname"`
	VIP      int    `json:"vip"`
}

// AccountInfo 帐号详细信息
type AccountInfo struct {
	UserID    int    `json:"user_id"`
	UserName  string `json:"nickname"`
	TotalSize string `json:"total_size"`
	UsedSize  string `json:"used_size"`
	VIP       int    `json:"vip"`
}

// filePageInfo 文件详情页提取的数据（内部使用）
type filePageInfo struct {
	Sign   string
	FileID string
	Dom    string
	IsNewd int
}
