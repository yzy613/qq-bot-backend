package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	CommandReq struct {
		g.Meta      `path:"/command" method:"post" tags:"api" summary:"命令"`
		Token       string `json:"token" v:"required"`
		Command     string `json:"command" v:"required"`
		GroupId     int64  `json:"group_id" v:"min:0"`
		MessageSync bool   `json:"message_sync" description:"同步发送信息"`
		Async       bool   `json:"async" description:"提前返回，异步执行命令"`
		Timestamp   int64  `json:"timestamp" v:"required" description:"单位：秒；超过 5 秒的请求会被拒绝"`
		Signature   string `json:"signature" v:"required" description:"以 token+command+group_id+timestamp+message_sync+async 为原文，以 token_name 为 key 的 HmacSHA256 值的 base64 值"`
	}
	CommandRes struct {
		Message string `json:"message"`
	}
)
