// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IThirdParty interface {
		QueryMinecraftGenuineUser(ctx context.Context, name string) (genuine bool, realName string, uuid string, err error)
	}
)

var (
	localThirdParty IThirdParty
)

func ThirdParty() IThirdParty {
	if localThirdParty == nil {
		panic("implement not found for interface IThirdParty, forgot register?")
	}
	return localThirdParty
}

func RegisterThirdParty(i IThirdParty) {
	localThirdParty = i
}
