package bot

import (
	"context"
	"errors"
	"github.com/bytedance/sonic/ast"
	"qq-bot-backend/utility/segment"
)

func (s *sBot) RewriteMessage(ctx context.Context, message string) {
	_, _ = s.reqJsonFromCtx(ctx).Set("raw_message", ast.NewString(message))
}

func (s *sBot) SetHistory(ctx context.Context, history string) error {
	const historyKey = "_history"
	node := s.reqJsonFromCtx(ctx)
	if !node.Get(historyKey).Valid() {
		_, _ = node.Set(historyKey, ast.NewNull())
	}
	if node.Get(historyKey).Get(history).Valid() {
		return errors.New("history already exists")
	}
	_, _ = node.Get(historyKey).Set(history, ast.NewNull())
	return nil
}

func (s *sBot) tryMessageSegmentToString(ctx context.Context) {
	node := s.reqJsonFromCtx(ctx)

	messageNode := node.Get("message")

	if messageNode.Exists() && messageNode.TypeSafe() != ast.V_ARRAY {
		return
	}

	jsonBytes, err := messageNode.MarshalJSON()
	if err != nil {
		return
	}

	mss, err := segment.ParseJSON(jsonBytes)
	if err != nil {
		return
	}

	_, _ = node.Set("message", ast.NewString(mss.String()))
	_, _ = node.Set("_is_message_segment", ast.NewObject([]ast.Pair{}))
}
