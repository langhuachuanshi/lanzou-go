package lanzou

// API 域名常量
const (
	baseURLPC     = "https://pc.woozooo.com"
	baseURLUpload = "https://up.woozooo.com"
	baseURLDown   = "https://down.woozooo.com"
	baseURLMy     = "https://www.woozooo.com"
)

// API 端点路径
const (
	pathLogin        = "/account/loginajax"
	pathLogout       = "/account/logout"
	pathFileList     = "/doupload.php"
	pathDirList      = "/doupload.php"
	pathUpload       = "/up.php"
	pathDownload     = "/ajaxm.php"
	pathFileInfo     = "/ajaxm.php"
	pathGetSign      = "/ajaxm.php"
	pathMoveFiles    = "/doupload.php"
	pathDeleteFiles  = "/doupload.php"
	pathSetPasswd    = "/doupload.php"
	pathNewFolder    = "/doupload.php"
	pathRecycleList  = "/doupload.php"
	pathRestoreFiles = "/doupload.php"
	pathCleanRecycle = "/doupload.php"
)

// 默认 User-Agent
const defaultUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

// 蓝奏云分享链接域名模式
var lanzouDomains = []string{
	"lanzoul.com",
	"lanzoui.com",
	"lanzout.com",
	"lanzouw.com",
	"lanzoux.com",
	"lanzouy.com",
	"lanzouf.com",
	"lanzouh.com",
	"lanzouj.com",
	"lanzouk.com",
	"lanzoup.com",
	"lanzouq.com",
	"lanzous.com",
}

// 默认配置
const (
	defaultMaxSize    = 100 * 1024 * 1024 // 100MB, 蓝奏云单文件限制
	defaultTimeout    = 30                 // 30秒
	defaultMaxDLCount = 3                  // 默认下载并发数
)

// ChallengeConfig acw_sc__v2 JS 挑战参数配置
// 蓝奏云会定期更换混淆JS，届时需要更新这些参数
//
// 配置文件示例 (challenge.json):
//
//	{
//	  "perm": [15,35,29,24,33,16,1,38,10,9,19,31,40,27,22,23,25,13,6,11,39,18,20,8,14,21,32,26,2,30,7,4,17,5,3,28,34,37,12,36],
//	  "xor_key": "3000176000856006061501533003690027800375"
//	}
type ChallengeConfig struct {
	// Perm 置换表（40个整数，值域1-40，不重复）
	Perm [40]int `json:"perm"`
	// XORKey XOR 密钥（40位十六进制字符串）
	XORKey string `json:"xor_key"`
}

// DefaultChallengeConfig 返回默认的挑战参数（2025年有效）
func DefaultChallengeConfig() *ChallengeConfig {
	return &ChallengeConfig{
		Perm: [40]int{
			15, 35, 29, 24, 33, 16, 1, 38, 10, 9,
			19, 31, 40, 27, 22, 23, 25, 13, 6, 11,
			39, 18, 20, 8, 14, 21, 32, 26, 2, 30,
			7, 4, 17, 5, 3, 28, 34, 37, 12, 36,
		},
		XORKey: "3000176000856006061501533003690027800375",
	}
}
