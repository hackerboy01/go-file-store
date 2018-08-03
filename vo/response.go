package vo

type ResponseMessage struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data map[string]interface{}	`json:"data"`
}

const ClientCreateFailed = 2001
const ClientCreateSuccess = 2000
const TokenCreateFailed = 2011
const TokenCreateSuccess = 2012
const WrongClientOdAndSecret = 2021
const UploadFileFailed = 2031
const UploadFileSuccess = 2032
const DownloadFailed = 2041
const DownloadSuccess = 2042
const DownloadFileLost = 2043
const InvalIdSlug = 2044
const InvalidToken = 2051
const ExpiredToken = 2052
const RequestError = 4000
