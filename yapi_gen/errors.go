package main

import "errors"

var (
	yapiReqError         = errors.New("yapi请求失败")
	yapiStructParseError = errors.New("yapi结构体生成失败")
)
