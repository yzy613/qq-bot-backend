// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	ICommand interface {
		TryCommand(ctx context.Context) (catch bool)
	}
)

var (
	localCommand ICommand
)

func Command() ICommand {
	if localCommand == nil {
		panic("implement not found for interface ICommand, forgot register?")
	}
	return localCommand
}

func RegisterCommand(i ICommand) {
	localCommand = i
}