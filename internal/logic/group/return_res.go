package group

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"qq-bot-backend/internal/dao"
	"qq-bot-backend/internal/model/entity"
	"qq-bot-backend/internal/service"
	"strings"
	"time"
)

func (s *sGroup) BindNamespaceReturnRes(ctx context.Context,
	groupId int64, namespace string) (retMsg string) {
	// 参数合法性校验
	if groupId == 0 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(ctx) ||
		(!service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, namespace, service.Bot().GetUserId(ctx)) &&
			!service.Namespace().IsNamespacePropertyPublic(ctx, namespace)) {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	var err error
	if groupE == nil {
		// 初始化 group 对象
		groupE = &entity.Group{
			GroupId:     groupId,
			Namespace:   namespace,
			SettingJson: "{}",
		}
		// 数据库插入
		_, err = dao.Group.Ctx(ctx).
			Data(groupE).
			OmitEmptyData().
			Insert()
	} else {
		if groupE.Namespace != "" {
			retMsg = "group(" + gconv.String(groupId) + ") 已经绑定了 namespace(" + groupE.Namespace + ")"
			return
		}
		// 重置 setting
		groupE = &entity.Group{
			Namespace:   namespace,
			SettingJson: "{}",
		}
		// 数据库更新
		_, err = dao.Group.Ctx(ctx).
			Where(dao.Group.Columns().GroupId, groupId).
			Data(groupE).
			OmitEmptyData().
			Update()
	}
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已绑定 group(" + gconv.String(groupId) + ") 到 namespace(" + namespace + ")"
	return
}

func (s *sGroup) UnbindReturnRes(ctx context.Context, groupId int64) (retMsg string) {
	// 参数合法性校验
	if groupId == 0 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(ctx) {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) &&
		!service.Namespace().IsNamespacePropertyPublic(ctx, groupE.Namespace) {
		return
	}
	// 数据库更新
	_, err := dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().Namespace, "").
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已解除 group(" + gconv.String(groupId) + ") 的 namespace 绑定"
	return
}

func (s *sGroup) CloneReturnRes(ctx context.Context, groupId, srcGroupId int64) (retMsg string) {
	// 参数合法性校验
	if groupId == 0 || srcGroupId == 0 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(ctx) {
		return
	}
	// 获取 group srcGroup
	groupE := getGroup(ctx, groupId)
	srcGroupE := getGroup(ctx, srcGroupId)
	if groupE == nil || srcGroupE == nil || groupE.Namespace == "" || srcGroupE.Namespace == "" {
		return
	}
	if groupE.Namespace != srcGroupE.Namespace {
		retMsg = "group(" + gconv.String(groupId) + ") 和 group(" + gconv.String(srcGroupId) + ") 的 namespace 不一致"
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) &&
		!service.Namespace().IsNamespacePropertyPublic(ctx, groupE.Namespace) {
		return
	}
	// 数据库更新
	_, err := dao.Group.Ctx(ctx).
		Where(dao.Group.Columns().GroupId, groupId).
		Data(dao.Group.Columns().SettingJson, srcGroupE.SettingJson).
		Update()
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 回执
	retMsg = "已克隆 group(" + gconv.String(srcGroupId) + ") 的配置到 group(" + gconv.String(groupId) + ")"
	return
}

func (s *sGroup) QueryGroupReturnRes(ctx context.Context, groupId int64) (retMsg string) {
	// 参数合法性校验
	if groupId == 0 {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) &&
		!service.Namespace().IsNamespacePropertyPublic(ctx, groupE.Namespace) {
		return
	}
	// 回执
	retMsg = dao.Group.Columns().Namespace + ": " + groupE.Namespace + "\n" +
		dao.Group.Columns().SettingJson + ": " + groupE.SettingJson + "\n" +
		dao.Group.Columns().UpdatedAt + ": " + groupE.UpdatedAt.String()
	return
}

func (s *sGroup) KickFromListReturnRes(ctx context.Context,
	groupId int64, listName string) (retMsg string) {
	// 参数合法性校验
	if groupId == 0 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(ctx) {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) &&
		!service.Namespace().IsNamespacePropertyPublic(ctx, groupE.Namespace) {
		return
	}
	// 是否存在 list
	lists := service.Namespace().GetNamespaceLists(ctx, groupE.Namespace)
	if _, ok := lists[listName]; !ok {
		retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + listName + ")"
		return
	}
	// 获取 list
	listMap := service.List().GetListData(ctx, listName)
	listMapLen := len(listMap)
	// 获取群成员列表
	membersArr, err := service.Bot().GetGroupMemberList(ctx, groupId, true)
	if err != nil {
		retMsg = "获取群成员列表失败"
		return
	}
	// 局部变量
	kickMap := make(map[string]any)
	// 解析数组
	for _, v := range membersArr {
		// map 断言
		if kk, ok := v.(map[string]any); ok {
			if kk["role"] != "member" {
				continue
			}
			// 放入待踢出 map
			userId := gconv.String(kk["user_id"])
			if _, ok := listMap[userId]; ok {
				kickMap[userId] = nil
			}
		}
	}
	// 不需要执行踢出
	if len(kickMap) == 0 {
		retMsg = "list(" + listName + ") 中的 group(" + gconv.String(groupId) + ") member 无需踢出"
		return
	}
	// 踢人过程
	retMsg = "正在踢出 list(" + listName + ") 中的 group(" + gconv.String(groupId) +
		") member\n共 " + gconv.String(listMapLen) + " 条，有 " + gconv.String(len(kickMap)) + " 条需要踢出"
	service.Bot().SendPlainMsgIfNotApiReq(ctx, retMsg)
	for k := range kickMap {
		// 踢人
		service.Bot().Kick(ctx, groupId, gconv.Int64(k))
		// 随机延时
		time.Sleep(time.Duration(grand.N(1000, 10000)) * time.Millisecond)
	}
	// 检查有没有踢出失败的
	// 获取群成员列表
	membersArr, err = service.Bot().GetGroupMemberList(ctx, groupId, true)
	if err != nil {
		retMsg = "获取群成员列表失败"
		return
	}
	// 局部变量
	kickMap = make(map[string]any)
	// 解析数组
	for _, v := range membersArr {
		// map 断言
		if kk, ok := v.(map[string]any); ok {
			if kk["role"] != "member" {
				continue
			}
			// 放入待踢出 map
			userId := gconv.String(kk["user_id"])
			if _, ok := listMap[userId]; ok {
				kickMap[userId] = nil
			}
		}
	}
	// 回执
	retMsg = "已踢出 list(" + listName + ") 中的 group(" + gconv.String(groupId) +
		") member\n共 " + gconv.String(listMapLen) + " 条，有 " + gconv.String(len(kickMap)) + " 条未踢出"
	return
}

func (s *sGroup) KeepFromListReturnRes(ctx context.Context,
	groupId int64, listName string) (retMsg string) {
	// 参数合法性校验
	if groupId == 0 {
		return
	}
	// 权限校验
	if !service.Bot().IsGroupOwnerOrAdminOrSysTrusted(ctx) {
		return
	}
	// 获取 group
	groupE := getGroup(ctx, groupId)
	if groupE == nil || groupE.Namespace == "" {
		return
	}
	// 权限校验
	if !service.Namespace().IsNamespaceOwnerOrAdminOrOperator(ctx, groupE.Namespace, service.Bot().GetUserId(ctx)) &&
		!service.Namespace().IsNamespacePropertyPublic(ctx, groupE.Namespace) {
		return
	}
	// 是否存在 list
	lists := service.Namespace().GetNamespaceLists(ctx, groupE.Namespace)
	if _, ok := lists[listName]; !ok {
		retMsg = "在 namespace(" + groupE.Namespace + ") 中未找到 list(" + listName + ")"
		return
	}
	// 获取 list
	listMap := service.List().GetListData(ctx, listName)
	listMapLen := len(listMap)
	// 获取群成员列表
	membersArr, err := service.Bot().GetGroupMemberList(ctx, groupId, true)
	if err != nil {
		retMsg = "获取群成员列表失败"
		return
	}
	// 局部变量
	kickMap := make(map[string]any)
	// 解析数组
	for _, v := range membersArr {
		// map 断言
		if kk, ok := v.(map[string]any); ok {
			if kk["role"] != "member" {
				continue
			}
			// 放入待踢出 map
			userId := gconv.String(kk["user_id"])
			if _, ok := listMap[userId]; !ok {
				kickMap[userId] = nil
			}
		}
	}
	// 不需要执行踢出
	if len(kickMap) == 0 {
		retMsg = "list(" + listName + ") 中的 group(" + gconv.String(groupId) + ") member 无需踢出"
		return
	}
	// 踢人过程
	retMsg = "正在踢出不在 list(" + listName + ") 中的 group(" + gconv.String(groupId) +
		") member\n共 " + gconv.String(listMapLen) + " 条，有 " + gconv.String(len(kickMap)) + " 条需要踢出"
	service.Bot().SendPlainMsgIfNotApiReq(ctx, retMsg)
	for k := range kickMap {
		// 踢人
		service.Bot().Kick(ctx, groupId, gconv.Int64(k))
		// 随机延时
		time.Sleep(time.Duration(grand.N(1000, 10000)) * time.Millisecond)
	}
	// 检查有没有踢出失败的
	// 获取群成员列表
	membersArr, err = service.Bot().GetGroupMemberList(ctx, groupId, true)
	if err != nil {
		retMsg = "获取群成员列表失败"
		return
	}
	// 局部变量
	kickMap = make(map[string]any)
	// 解析数组
	for _, v := range membersArr {
		// map 断言
		if kk, ok := v.(map[string]any); ok {
			if kk["role"] != "member" {
				continue
			}
			// 放入待踢出 map
			userId := gconv.String(kk["user_id"])
			if _, ok := listMap[userId]; !ok {
				kickMap[userId] = nil
			}
		}
	}
	// 回执
	retMsg = "已踢出不在 list(" + listName + ") 中的 group(" + gconv.String(groupId) +
		") member\n共 " + gconv.String(listMapLen) + " 条，有 " + gconv.String(len(kickMap)) + " 条未踢出"
	return
}

func (s *sGroup) CheckExistReturnRes(ctx context.Context) (retMsg string) {
	pageNum, pageSize := 1, 10
	var msgBuilder strings.Builder
	waitForDelGroups := make([]int64, 0)
	loginUserId, _ := service.Bot().GetLoginInfo(ctx)
	for {
		var groupE []*entity.Group
		err := dao.Group.Ctx(ctx).
			Fields(dao.Group.Columns().GroupId).
			Page(pageNum, pageSize).
			Scan(&groupE)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		for _, v := range groupE {
			vGroupId := v.GroupId
			_, err = service.Bot().GetGroupInfo(ctx, vGroupId, true)
			// 判断群是否已经解散或者登录账号不在群内
			if err == nil {
				_, err = service.Bot().GetGroupMemberInfo(ctx, vGroupId, loginUserId)
				if err == nil {
					continue
				}
			} else if err.Error() != "群聊不存在" {
				retMsg = "获取群信息失败"
				return
			}
			// 记录信息
			msgBuilder.WriteString("\ngroup(" + gconv.String(vGroupId) + ")")
			// 放入待删除数组
			waitForDelGroups = append(waitForDelGroups, vGroupId)
		}
		// 结束条件
		if len(groupE) < pageSize {
			break
		}
		pageNum++
	}
	if len(waitForDelGroups) > 0 {
		// 删除 groups
		_, err := dao.Group.Ctx(ctx).Where(dao.Group.Columns().GroupId, waitForDelGroups).Delete()
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	// 回执
	retMsg = "已删除 " + gconv.String(len(waitForDelGroups)) + " 条不适用的 group" + msgBuilder.String()
	return
}
