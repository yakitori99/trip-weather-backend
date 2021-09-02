package utils

import (
	"runtime"

	log "github.com/sirupsen/logrus"
)

// errとmessageを受け取り、関数名、ファイル名、行数を追加しエラーログを出力する
func OutErrorLog(msg string, err error) {
	pt, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pt).Name()

	log.WithFields(log.Fields{
		"err":      err,
		"funcName": funcName,
		"file":     file,
		"line":     line,
	}).Error(msg)
}

// errとmessageとdetailを受け取り、関数名、ファイル名、行数を追加しエラーログを出力する
func OutErrorLogDetail(msg string, err error, detail string) {
	pt, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pt).Name()

	log.WithFields(log.Fields{
		"err":      err,
		"funcName": funcName,
		"file":     file,
		"line":     line,
		"detail":   detail,
	}).Error(msg)
}

// messageを受け取り、関数名を追加しログを出力する
func OutInfoLog(msg string) {
	pt, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pt).Name()

	log.WithFields(log.Fields{
		"funcName": funcName,
	}).Info(msg)
}

// messageとrequestURLを受け取り、関数名を追加しログを出力する
func OutInfoLogRequest(msg string, requestURL string) {
	pt, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pt).Name()

	log.WithFields(log.Fields{
		"funcName":   funcName,
		"requestURL": requestURL,
	}).Info(msg)
}

// messageを受け取り、関数名を追加しログを出力する
func OutDebugLog(msg string) {
	pt, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pt).Name()

	log.WithFields(log.Fields{
		"funcName": funcName,
	}).Debug(msg)
}
