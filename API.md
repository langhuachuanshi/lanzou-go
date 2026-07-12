# API 参考

> lanzou-go 完整 API 文档，按功能模块组织。

---

## NewClient

```go
func NewClient(opts ...Option) *Client
```

| Option | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `WithTimeout(sec)` | `int` | `30` | HTTP 超时（秒） |
| `WithMaxSize(bytes)` | `int` | `104857600` | 单文件大小限制 |
| `WithMaxDownloadCount(n)` | `int` | `3` | 并发下载数 |
| `WithUploadDelay(min, max)` | `int, int` | `0, 0` | 上传延迟（毫秒） |
| `WithHTTPClient(c)` | `*http.Client` | 默认 | 自定义 HTTP 客户端 |
| `WithChallengeConfig(c)` | `*ChallengeConfig` | 内置默认值 | 反爬挑战参数 |

### 运行时配置

```go
client.SetTimeout(60)
client.SetMaxSize(50 * 1024 * 1024)
client.SetMaxDownloadCount(5)
client.SetChallengeConfig(&cfg)
// cookie 注入
client.SetCookiesFromMap(map[string]string{...})
```

---

## 一、直链解析（无需登录）

### GetDurlByURL

```go
func (c *Client) GetDurlByURL(shareURL, pwd string) (string, error)
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| shareURL | string | ✅ | 蓝奏云分享链接 |
| pwd | string | ❌ | 密码，无密码传 `""` |

返回：直链 URL

```go
durl, err := client.GetDurlByURL("https://36cq.lanzouo.com/i5qJz2z14zoh", "")
```

### GetFileInfo

```go
func (c *Client) GetFileInfo(shareURL, pwd string) (*FileDetail, error)
```

返回文件详情（含直链）。

```go
type FileDetail struct {
    FileID      string `json:"file_id"`
    NameAll     string `json:"name_all"`   // 文件名
    Size        string `json:"size"`       // 文件大小
    DURL        string `json:"durl"`       // 直链
    DownloadURL string `json:"url"`        // 分享链接
}
```

---

## 二、账号操作

### Login

```go
func (c *Client) Login(user, pwd string) error
```

### Logout

```go
func (c *Client) Logout() error
```

### GetUserInfo

```go
func (c *Client) GetUserInfo() (*UserInfo, error)
```

### GetAccountInfo

```go
func (c *Client) GetAccountInfo() (*AccountInfo, error)
```

---

## 三、文件操作

### GetFileList

```go
func (c *Client) GetFileList(fid int) (*FileList, error)
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| fid | int | ✅ | 文件夹ID，根目录传 `0` |

```go
type FileInfo struct {
    ID      string `json:"id"`       // 文件ID
    NameAll string `json:"name_all"` // 文件名
    Size    string `json:"size"`     // 大小
    Time    string `json:"time"`     // 上传时间
    Icon    string `json:"icon"`     // 图标类型
    Downs   string `json:"downs"`    // 下载次数
}
```

### GetFileInfoByURL

```go
func (c *Client) GetFileInfoByURL(shareURL, pwd string) (*FileDetail, error)
```

等价于 `GetFileInfo`。

### DownloadFile

```go
func (c *Client) DownloadFile(savePath, shareURL string, pwd ...string) error
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| savePath | string | ✅ | 保存路径 |
| shareURL | string | ✅ | 分享链接 |
| pwd | ...string | ❌ | 密码 |

### DownloadFile2

```go
func (c *Client) DownloadFile2(saveDir, shareURL, pwd string) error
```

自动获取文件名，保存到指定目录。

### DownloadDir

```go
func (c *Client) DownloadDir(saveDir string, fid int) error
```

递归下载整个文件夹（含子文件夹）。

### UploadFile

```go
func (c *Client) UploadFile(filePath string, fid int, desc ...string) (*UploadResult, error)
```

上传本地文件。

### MoveFiles

```go
func (c *Client) MoveFiles(fids []string, fid int) error
```

### DeleteFiles

```go
func (c *Client) DeleteFiles(fids []string) error
```

### SetPassword

```go
func (c *Client) SetPassword(entityID int, pwd string) error
```

---

## 四、文件夹操作

### GetDirList

```go
func (c *Client) GetDirList(fid int) (*FolderList, error)
```

```go
type FolderInfo struct {
    FolID  string `json:"fol_id"`  // 文件夹ID
    Name   string `json:"name"`    // 名称
    Onof   string `json:"onof"`    // 密码开关
}
```

### NewFolder

```go
func (c *Client) NewFolder(name string, parentID int) (*FolderInfo, error)
```

### DeleteFolder / MoveFolder

```go
func (c *Client) DeleteFolder(fids []string) error
func (c *Client) MoveFolder(fids []string, fid int) error
```

---

## 五、回收站

| 方法 | 说明 |
|------|------|
| `GetRecycleList(page int)` | 获取回收站列表 |
| `MoveToTrash(fids []string)` | 移入回收站 |
| `RestoreFiles(fids []string)` | 恢复文件 |
| `CleanRecycle()` | 清空回收站 |

---

## 错误码

```go
ErrNotLoggedIn    // 未登录
ErrFileExpired    // 文件已过期
ErrPasswordWrong  // 密码错误
ErrPasswordNeeded // 需要密码
ErrFileSizeLimit  // 超过大小限制
ErrInvalidURL     // 无效链接
ErrExtractFailed  // 页面解析失败
ErrUploadFailed   // 上传失败
ErrDownloadFailed // 下载失败
ErrAPIError       // API 错误
```

---

## 挑战参数更新

```go
// 蓝奏云换混淆时更新
client.SetChallengeConfig(&lanzou.ChallengeConfig{
    Perm:   [40]int{/* 新置换表 */},
    XORKey: "/* 新XOR密钥 */",
})
```

详见 [README.md](./README.md#蓝奏云换混淆时如何更新)
