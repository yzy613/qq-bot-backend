package event

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"net/http"
	"net/url"
	"qq-bot-backend/internal/service"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	atPrefixRe      = regexp.MustCompile(`^\[CQ:at,qq=([^,]+)(?:,[^,=]+=[^,]*)*]\s*`)
	webhookPrefixRe = regexp.MustCompile(`^webhook(?::([A-Za-z]{3,7}))?(?:#([\s\S]+)#)?(?:<([\s\S]+)>)?(?:@(.+)@)?://(.+)$`)
	commandPrefixRe = regexp.MustCompile(`^(?:command|cmd)://([\s\S]+)$`)
	rewritePrefixRe = regexp.MustCompile(`^rewrite://([\s\S]+)$`)
	placeholderRe   = regexp.MustCompile(`\{(.+?)(\d+)?}`)
)

func decrementPlaceholderIndex(text string) string {
	arr := placeholderRe.FindAllStringSubmatch(text, -1)
	for _, sub := range arr {
		if len(sub) < 3 {
			continue
		}
		num := gconv.Int(sub[2]) - 1
		if num <= 0 {
			text = strings.ReplaceAll(text, sub[0], "{"+sub[1]+"}")
			continue
		}
		text = strings.ReplaceAll(text, sub[0], "{"+sub[1]+gconv.String(num)+"}")
	}
	return text
}

func (s *sEvent) TryKeywordReply(ctx context.Context) (catch bool) {
	// 获取基础信息
	msg := service.Bot().GetMessage(ctx)
	userId := service.Bot().GetUserId(ctx)
	// 匹配 @bot
	if atPrefixRe.MatchString(msg) {
		sub := atPrefixRe.FindStringSubmatch(msg)
		if sub[1] == gconv.String(service.Bot().GetSelfId(ctx)) {
			msg = strings.Replace(msg, sub[0], "", 1)
		}
	}
	// 匹配关键词
	contains, hit, value := service.Util().IsOnKeywordLists(ctx, msg, service.Namespace().GetGlobalNamespaceLists(ctx))
	if !contains || value == "" {
		return
	}
	// 匹配成功，回复
	replyMsg := value
	noReplyPrefix := false
	switch {
	case webhookPrefixRe.MatchString(value):
		replyMsg, noReplyPrefix = s.keywordReplyWebhook(ctx,
			userId, 0, service.Bot().GetNickname(ctx),
			msg, hit, value)
	case rewritePrefixRe.MatchString(value):
		catch = s.keywordReplyRewrite(ctx, s.TryKeywordReply, msg, hit, value)
		replyMsg = ""
	case commandPrefixRe.MatchString(value):
		replyMsg = s.keywordReplyCommand(ctx, msg, hit, value)
	}
	// 内容为空，不回复
	if replyMsg == "" {
		return
	}
	// 限速
	const kind = "replyU"
	uid := gconv.String(userId)
	if limited, _ := service.Util().AutoLimit(ctx, kind, uid, 5, time.Minute); limited {
		g.Log().Notice(ctx, kind, uid, "is limited")
		return
	}
	if !noReplyPrefix {
		replyMsg = "[CQ:reply,id=" + gconv.String(service.Bot().GetMsgId(ctx)) + "]" + replyMsg
	}
	service.Bot().SendMsg(ctx, replyMsg)
	catch = true
	return
}

func (s *sEvent) keywordReplyWebhook(ctx context.Context, userId, groupId int64, nickname,
	message, hit, value string) (replyMsg string, noReplyPrefix bool) {
	// 必须以 hit 开头
	if !strings.HasPrefix(message, hit) {
		return
	}
	// Url
	subMatch := webhookPrefixRe.FindStringSubmatch(service.Codec().DecodeCqCode(value))
	method := strings.ToUpper(subMatch[1])
	if method == "" {
		method = http.MethodGet
	}
	headers := subMatch[2]
	payload := subMatch[3]
	bodyPath := strings.Split(subMatch[4], ".")
	urlLink := subMatch[5]
	// Arguments
	var err error
	message = service.Codec().DecodeCqCode(message)
	hit = service.Codec().DecodeCqCode(hit)
	remain := strings.Replace(message, hit, "", 1)
	// Headers
	if headers != "" {
		headers = strings.ReplaceAll(headers, "\\n", "\n")
		headers = strings.ReplaceAll(headers, "\r", "\n")
		headers = strings.ReplaceAll(headers, "{message}", message)
		headers = strings.ReplaceAll(headers, "{remain}", remain)
		headers = strings.ReplaceAll(headers, "{nickname}", nickname)
		headers = strings.ReplaceAll(headers, "{userId}", gconv.String(userId))
		headers = strings.ReplaceAll(headers, "{groupId}", gconv.String(groupId))
	}
	// Url escape
	urlLink = strings.ReplaceAll(urlLink, "{message}", url.QueryEscape(message))
	urlLink = strings.ReplaceAll(urlLink, "{remain}", url.QueryEscape(remain))
	urlLink = strings.ReplaceAll(urlLink, "{nickname}", url.QueryEscape(nickname))
	urlLink = strings.ReplaceAll(urlLink, "{userId}", gconv.String(userId))
	urlLink = strings.ReplaceAll(urlLink, "{groupId}", gconv.String(groupId))
	// Call webhook
	var body []byte
	var statusCode int
	var contentType string
	switch method {
	case http.MethodGet:
		statusCode, contentType, body, err = service.Util().WebhookGetHeadConnectOptionsTrace(ctx, headers, method, urlLink)
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		// Payload
		msg, _ := sonic.ConfigDefault.MarshalToString(message)
		r, _ := sonic.ConfigDefault.MarshalToString(remain)
		nick, _ := sonic.ConfigDefault.MarshalToString(nickname)
		// 占位符替换
		payload = strings.ReplaceAll(payload, "{message}", msg)
		payload = strings.ReplaceAll(payload, "{remain}", r)
		payload = strings.ReplaceAll(payload, "{nickname}", nick)
		payload = strings.ReplaceAll(payload, "{userId}", gconv.String(userId))
		payload = strings.ReplaceAll(payload, "{groupId}", gconv.String(groupId))
		statusCode, contentType, body, err = service.Util().WebhookPostPutPatchDelete(ctx, headers, method, urlLink, payload)
	default:
		return
	}
	if err != nil {
		g.Log().Notice(ctx, "webhook", statusCode, method, urlLink, message, err)
		return
	}
	// Log
	if statusCode != http.StatusOK {
		g.Log().Notice(ctx,
			nickname+"("+gconv.String(userId)+") in group("+gconv.String(service.Bot().GetGroupId(ctx))+
				") call webhook", statusCode, method, urlLink, message)
	}
	// 媒体文件
	{
		var mediumUrl string
		// 如果是图片
		if strings.HasPrefix(contentType, "image/") {
			mediumUrl, err = service.File().CacheFile(ctx, body, 5*time.Minute)
			if err != nil {
				replyMsg = "Image cache failed"
				return
			}
			replyMsg = "[CQ:image,file=" + mediumUrl + "]"
			return
		}
		// 如果是音频
		if strings.HasPrefix(contentType, "audio/") {
			mediumUrl, err = service.File().CacheFile(ctx, body, 5*time.Minute)
			if err != nil {
				replyMsg = "Audio cache failed"
				return
			}
			replyMsg = "[CQ:record,file=" + mediumUrl + "]"
			noReplyPrefix = true
			return
		}
		// 如果是视频
		if strings.HasPrefix(contentType, "video/") {
			mediumUrl, err = service.File().CacheFile(ctx, body, 5*time.Minute)
			if err != nil {
				replyMsg = "Video cache failed"
				return
			}
			replyMsg = "[CQ:video,file=" + mediumUrl + "]"
			noReplyPrefix = true
			return
		}
	}
	// 没有 bodyPath，直接返回 body
	if len(bodyPath) == 1 && bodyPath[0] == "" {
		replyMsg = string(body)
		return
	}
	// 默认视为 JSON 数据
	path := make([]any, len(bodyPath))
	for i, v := range bodyPath {
		index, e := strconv.Atoi(v)
		if e == nil {
			path[i] = index
			continue
		}
		path[i] = v
	}
	// 解析 body 获取数据
	node, err := sonic.Get(body, path...)
	if err != nil {
		replyMsg = "Wrong JSON path"
		return
	}
	if node.TypeSafe() != ast.V_STRING {
		err = node.LoadAll()
		if err != nil {
			replyMsg = "Wrong JSON format"
			return
		}
		replyMsg, _ = node.Raw()
		return
	}
	replyMsg, _ = node.StrictString()
	return
}

func (s *sEvent) keywordReplyCommand(ctx context.Context, message, hit, text string) (replyMsg string) {
	// 必须以 hit 开头
	if !strings.HasPrefix(message, hit) {
		return
	}
	// 解码提取
	subMatch := commandPrefixRe.FindStringSubmatch(service.Codec().DecodeCqCode(text))
	// 占位符替换
	remain := strings.Replace(message, hit, "", 1)
	subMatch[1] = strings.ReplaceAll(subMatch[1], "{message}", message)
	subMatch[1] = strings.ReplaceAll(subMatch[1], "{remain}", remain)
	// 转换占位符
	subMatch[1] = decrementPlaceholderIndex(subMatch[1])
	// 为什么是 " &&"？因为 " &&" 后可能是换行符，需要替换为 " "
	subMatch[1] = strings.ReplaceAll(subMatch[1], " &&\r", " && ")
	subMatch[1] = strings.ReplaceAll(subMatch[1], " &&\n", " && ")
	// 切分命令
	commands := strings.Split(subMatch[1], " && ")
	var replyBuilder strings.Builder
	for _, command := range commands {
		catch, tmp := service.Command().TryCommand(ctx, strings.TrimSpace(command))
		if !catch {
			return
		}
		replyBuilder.WriteString(tmp + "\n")
	}
	replyMsg = strings.TrimRight(replyBuilder.String(), "\n")
	return
}

func (s *sEvent) keywordReplyRewrite(ctx context.Context, try func(context.Context) bool, message, hit, text string) (catch bool) {
	// 必须以 hit 开头
	if !strings.HasPrefix(message, hit) {
		return
	}
	// 防止循环递归
	err := service.Bot().SetHistory(ctx, hit)
	if err != nil {
		// rewrite loop detected
		g.Log().Notice(ctx, "rewrite loop detected: "+hit)
		return
	}
	// 解码提取
	subMatch := rewritePrefixRe.FindStringSubmatch(service.Codec().DecodeCqCode(text))
	// 占位符替换
	remain := strings.Replace(message, hit, "", 1)
	subMatch[1] = strings.ReplaceAll(subMatch[1], "{message}", message)
	subMatch[1] = strings.ReplaceAll(subMatch[1], "{remain}", remain)
	// 为什么是 " &"？因为 " &" 后可能是换行符，需要替换为 " "
	subMatch[1] = strings.ReplaceAll(subMatch[1], " &\r", " & ")
	subMatch[1] = strings.ReplaceAll(subMatch[1], " &\n", " & ")
	// 切分
	rewrites := strings.Split(subMatch[1], " & ")
	for _, rewrite := range rewrites {
		service.Bot().RewriteMessage(ctx, strings.TrimSpace(rewrite))
		// callback
		catch = try(ctx)
	}
	return
}
