// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

type (
	ICodec interface {
		DecodeCqCode(src string) (dest string)
		EncodeCqCode(src string) (dest string)
		IsIncludeCqCode(str string) bool
		DecodeBlank(src string) (dest string)
	}
)

var (
	localCodec ICodec
)

func Codec() ICodec {
	if localCodec == nil {
		panic("implement not found for interface ICodec, forgot register?")
	}
	return localCodec
}

func RegisterCodec(i ICodec) {
	localCodec = i
}
