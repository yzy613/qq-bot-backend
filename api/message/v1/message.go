package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	MessageReq struct {
		g.Meta  `path:"/message" method:"post" tags:"api" summary:"消息"`
		Token   string `json:"token" description:"support Authorization header"`
		Message string `json:"message" v:"required"`
		UserId  int64  `json:"user_id" v:"min:0|required-without:GroupId"`
		GroupId int64  `json:"group_id" v:"min:0|required-without:UserId"`
	}
	MessageRes struct{}
)
