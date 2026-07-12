package lanzou

import (
	"crypto/aes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// ===== JS 挑战求解器 (acw_sc__v2) =====

// reArg1 提取 arg1 的正则
var reArg1 = regexp.MustCompile(`var\s+arg1\s*=\s*'([A-F0-9]{40})'`)

// solveAcwScV2 计算 acw_sc__v2 cookie 值
// 蓝奏云返回一个 JS 挑战页面，需要执行以下算法：
// 1. 提取 arg1（40位十六进制字符串）
// 2. 用置换表 m 对 arg1 重排序
// 3. 与 XOR 密钥异或得到 cookie 值
//
// cfg 为 nil 时使用 DefaultChallengeConfig()
func solveAcwScV2(html string, cfg *ChallengeConfig) (string, error) {
	if cfg == nil {
		cfg = DefaultChallengeConfig()
	}

	matches := reArg1.FindStringSubmatch(html)
	if len(matches) < 2 {
		return "", fmt.Errorf("arg1 not found in challenge page")
	}
	arg1 := matches[1]

	// 置换
	var q [40]byte
	for x := 0; x < len(arg1); x++ {
		for z := 0; z < len(cfg.Perm); z++ {
			if cfg.Perm[z] == x+1 {
				q[z] = arg1[x]
			}
		}
	}
	u := string(q[:])

	// XOR
	xorKey := cfg.XORKey
	var v strings.Builder
	for x := 0; x < len(u) && x < len(xorKey); x += 2 {
		a, _ := strconv.ParseUint(u[x:x+2], 16, 8)
		b, _ := strconv.ParseUint(xorKey[x:x+2], 16, 8)
		xor := a ^ b
		v.WriteString(fmt.Sprintf("%02x", xor))
	}
	return v.String(), nil
}

// isChallengePage 检查是否为 JS 挑战页面
func isChallengePage(html string) bool {
	return strings.Contains(html, "var arg1=") && strings.Contains(html, "acw_sc__v2")
}

// ===== 页面解析正则 =====

// reSign 签名提取正则
var reSign = regexp.MustCompile(`sign\s*[:=]\s*'([a-zA-Z0-9]+)'`)

// reFileID 文件ID提取正则
var reFileID = regexp.MustCompile(`file_id\s*[:=]\s*'?(\d+)`)

// reDom 域名提取正则
var reDom = regexp.MustCompile(`dom\s*[:=]\s*'([^']+)'`)

// reIsNewd 新直链标记正则
var reIsNewd = regexp.MustCompile(`is_newd\s*[:=]\s*(\d+)`)

// reWpSign 从iframe页面提取 wp_sign
var reWpSign = regexp.MustCompile(`wp_sign\s*=\s*'([^']+)'`)

// reIframeSrc 从主页提取 iframe src
var reIframeSrc = regexp.MustCompile(`<iframe[^>]+src="(/fn\?[^"]+)"`)

// reFid 从页面提取 fid
var reFid = regexp.MustCompile(`var\s+fid\s*=\s*(\d+)`)

// reTitle 从页面提取文件名
var reTitle = regexp.MustCompile(`<title>([^<]+)</title>`)

// reFileSizeDesc 从页面提取文件大小描述
var reFileSizeDesc = regexp.MustCompile(`文件大小[：:]\s*([^<\s]+(?:\s*[A-Za-z]+)?)`)

// reFileSizeMeta 从meta description提取文件大小
var reFileSizeMeta = regexp.MustCompile(`文件大小[：:]\s*(\d+\.?\d*\s*[A-Za-z]+)`)

// ===== 数据提取函数 =====

// pageData 从分享页面提取的数据
type pageData struct {
	Fid      int    // 文件ID
	Title    string // 文件名
	FileSize string // 文件大小
	Iframe   string // iframe路径
}

// extractPageData 从主页HTML提取数据
func extractPageData(html string) (*pageData, error) {
	data := &pageData{}

	// 提取 fid
	if m := reFid.FindStringSubmatch(html); len(m) > 1 {
		data.Fid, _ = strconv.Atoi(m[1])
	}

	// 提取标题（文件名）
	if m := reTitle.FindStringSubmatch(html); len(m) > 1 {
		title := strings.TrimSpace(m[1])
		title = strings.TrimSuffix(title, " - 蓝奏云")
		title = strings.TrimSuffix(title, " - 蓝奏云盘")
		data.Title = title
	}

	// 提取文件大小
	if m := reFileSizeDesc.FindStringSubmatch(html); len(m) > 1 {
		data.FileSize = strings.TrimSpace(m[1])
	} else if m := reFileSizeMeta.FindStringSubmatch(html); len(m) > 1 {
		data.FileSize = strings.TrimSpace(m[1])
	}

	// 提取 iframe src
	if m := reIframeSrc.FindStringSubmatch(html); len(m) > 1 {
		data.Iframe = m[1]
	}

	if data.Fid == 0 && data.Iframe == "" {
		return nil, fmt.Errorf("%w: no fid or iframe found", ErrExtractFailed)
	}

	return data, nil
}

// iframeData 从iframe页面提取的数据
type iframeData struct {
	WpSign    string // 签名
	Fid       int    // 文件ID（从ajaxm.php URL提取）
	AjaxmURL  string // ajaxm.php 完整URL
}

// extractIframeData 从iframe页面HTML提取数据
func extractIframeData(html string) (*iframeData, error) {
	data := &iframeData{}

	// 提取 wp_sign
	if m := reWpSign.FindStringSubmatch(html); len(m) > 1 {
		data.WpSign = m[1]
	} else {
		return nil, fmt.Errorf("%w: wp_sign not found", ErrExtractFailed)
	}

	// 提取 ajaxm.php URL 中的 file ID
	reAjaxm := regexp.MustCompile(`url\s*:\s*'(/ajaxm\.php\?file=(\d+))'`)
	if m := reAjaxm.FindStringSubmatch(html); len(m) > 2 {
		data.AjaxmURL = m[1]
		data.Fid, _ = strconv.Atoi(m[2])
	}

	return data, nil
}

// ajaxmResp ajaxm.php 的响应
type ajaxmResp struct {
	Zt   int             `json:"zt"`
	Dom  string          `json:"dom"`
	Url  string          `json:"url"`
	Inf  json.RawMessage `json:"inf"` // 可能是 string 或 int
	Mess string          `json:"mess"`
}

// ===== AES 解密工具 =====

// decryptAES ECB模式解密（用于直链接解密）
func decryptAES(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}
	decrypted := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += aes.BlockSize {
		block.Decrypt(decrypted[i:i+aes.BlockSize], ciphertext[i:i+aes.BlockSize])
	}
	// 去除 PKCS7 padding
	if len(decrypted) > 0 {
		padLen := int(decrypted[len(decrypted)-1])
		if padLen > 0 && padLen <= aes.BlockSize {
			decrypted = decrypted[:len(decrypted)-padLen]
		}
	}
	return decrypted, nil
}

// hexDecode 十六进制解码
func hexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

// ===== 通用工具 =====

// isLanzouURL 检查是否为蓝奏云链接
func isLanzouURL(url string) bool {
	for _, domain := range lanzouDomains {
		if strings.Contains(url, domain) {
			return true
		}
	}
	return strings.Contains(url, "lanzou") || strings.Contains(url, "woozooo")
}

// normalizeURL 标准化蓝奏云URL
func normalizeURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	rawURL = strings.TrimSuffix(rawURL, "/")
	return rawURL
}

// getBaseHost 从URL提取 host
func getBaseHost(rawURL string) string {
	rawURL = strings.TrimPrefix(rawURL, "https://")
	rawURL = strings.TrimPrefix(rawURL, "http://")
	if idx := strings.Index(rawURL, "/"); idx > 0 {
		return rawURL[:idx]
	}
	return rawURL
}

// randomString 生成随机字符串（用于boundary等）
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
