package group

import (
	"context"
	sj "github.com/bitly/go-simplejson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/service"
)

func (s *sGroup) SetAntiRecallReturnRes(ctx context.Context, groupId int64, enable bool) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if enable {
		if _, ok := settingJson.CheckGet(antiRecallKey); ok {
			service.Bot().SendPlainMsg(ctx, "早已启用 group("+gconv.String(groupId)+") 群反撤回")
			return
		}
		settingJson.Set(antiRecallKey, true)
	} else {
		if _, ok := settingJson.CheckGet(antiRecallKey); ok {
			settingJson.Del(antiRecallKey)
		} else {
			service.Bot().SendPlainMsg(ctx, "并未启用 group("+gconv.String(groupId)+") 群反撤回")
			return
		}
	}
	// 保存数据
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	if enable {
		service.Bot().SendPlainMsg(ctx, "已启用 group("+gconv.String(groupId)+") 群反撤回")
	} else {
		service.Bot().SendPlainMsg(ctx, "已禁用 group("+gconv.String(groupId)+") 群反撤回")
	}
}

func (s *sGroup) SetMessageNotificationReturnRes(ctx context.Context, groupId int64, notificationGroupId int64) {
	// 参数合法性校验
	if groupId < 1 || notificationGroupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if v, ok := settingJson.CheckGet(messageNotificationGroupIdKey); ok {
		service.Bot().SendPlainMsg(ctx,
			"早已设置 group("+gconv.String(groupId)+") 群消息通知群为 group("+
				gconv.String(v.MustInt64())+")")
		return
	}
	// 验证是否存在该群
	_, err = service.Bot().GetGroupInfo(ctx, notificationGroupId)
	if err != nil {
		service.Bot().SendPlainMsg(ctx, "group("+gconv.String(notificationGroupId)+") 未找到")
		return
	}
	settingJson.Set(messageNotificationGroupIdKey, notificationGroupId)
	// 保存数据
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已设置 group("+gconv.String(groupId)+") 群消息通知群为 group("+gconv.String(notificationGroupId)+")")
}

func (s *sGroup) RemoveMessageNotificationReturnRes(ctx context.Context, groupId int64) {
	// 参数合法性校验
	if groupId < 1 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdmin(ctx) {
		return
	}
	// 获取 group
	gEntity := getGroup(ctx, groupId)
	if gEntity == nil {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdmin(ctx, gEntity.Namespace, service.Bot().GetUserId(ctx)) {
		return
	}
	// 数据处理
	settingJson, err := sj.NewJson([]byte(gEntity.SettingJson))
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	if _, ok := settingJson.CheckGet(messageNotificationGroupIdKey); !ok {
		service.Bot().SendPlainMsg(ctx, "并未设置 group("+gconv.String(groupId)+") 群消息通知群")
		return
	}
	settingJson.Del(messageNotificationGroupIdKey)
	// 保存数据
	settingBytes, err := settingJson.Encode()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 数据库更新
	_, err = dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().SettingJson, string(settingBytes)).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	service.Bot().SendPlainMsg(ctx, "已移除 group("+gconv.String(groupId)+") 群消息通知群")
}