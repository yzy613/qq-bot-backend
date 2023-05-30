// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"time"
)

type (
	IFile interface {
		GetCachedFileFromId(ctx context.Context, id string) (content string, err error)
		SetCachedFile(ctx context.Context, content string, duration time.Duration) (id string, err error)
		GetCachedFileUrl(ctx context.Context, id string) (url string, err error)
	}
)

var (
	localFile IFile
)

func File() IFile {
	if localFile == nil {
		panic("implement not found for interface IFile, forgot register?")
	}
	return localFile
}

func RegisterFile(i IFile) {
	localFile = i
}
