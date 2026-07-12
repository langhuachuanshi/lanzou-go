package lanzou

import "errors"

// 业务错误码
var (
	ErrNotLoggedIn    = errors.New("lanzou: not logged in")
	ErrFileExpired    = errors.New("lanzou: file expired or deleted")
	ErrFileNotFound   = errors.New("lanzou: file not found")
	ErrPasswordWrong  = errors.New("lanzou: wrong password")
	ErrPasswordNeeded = errors.New("lanzou: password required")
	ErrUploadFailed   = errors.New("lanzou: upload failed")
	ErrDownloadFailed = errors.New("lanzou: download failed")
	ErrRateLimited    = errors.New("lanzou: rate limited, please try later")
	ErrAPIError       = errors.New("lanzou: api error")
	ErrFileSizeLimit  = errors.New("lanzou: file size exceeds limit")
	ErrFolderNotEmpty = errors.New("lanzou: folder not empty")
	ErrInvalidURL     = errors.New("lanzou: invalid lanzou url")
	ErrExtractFailed  = errors.New("lanzou: failed to extract sign/id from page")
	ErrDecryptFailed  = errors.New("lanzou: failed to decrypt download url")
)
