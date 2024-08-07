v1.8

\<required> [optional]

[TOC]

## List

| Command                                           | Description                                                                                     | Comment                                                        |
|---------------------------------------------------|-------------------------------------------------------------------------------------------------|----------------------------------------------------------------|
| /list add \<list_name> \<namespace>               | 在 namespace 下新增 list_name                                                                       | 需要 list_name 对应的 namespace admin 或更高权限或者 namespace 有 public 属性 |
| /list join \<list_name> \<key> [value]            | 将 key[:value] 添加到 list_name 中（key value 可作为双因子认证使用）（key 包含**空格**请用`%20`转义替换，包含**%**请用`%25`转义替换） | 同上                                                             |
| /list copy-key \<list_name> \<src_key> \<dst_key> | 将 list_name 中的 src_key 复制到 dst_key 中（src_key dst_key 包含**空格**请用`%20`转义替换，包含**%**请用`%25`转义替换）    | 同上                                                             |
| /list leave \<list_name> \<key>                   | 将 key 从 list_name 中移除（key 包含**空格**请用`%20`转义替换，包含**%**请用`%25`转义替换）                               | 同上                                                             |
| /list len \<list_name>                            | 查询 list_name 里共有多少条数据                                                                           | 同上                                                             |
| /list query \<list_name> [key]                    | 查询 list_name 或者 list_name 内的 key（key 包含**空格**请用`%20`转义替换，包含**%**请用`%25`转义替换）                    | 同上                                                             |
| /list export \<list_name>                         | 将 list_name 以文件的方式导出                                                                            | 同上                                                             |
| /list append \<list_name> \<json>                 | 用 json 追加 list_name 的数据                                                                         | 同上                                                             |
| /list set \<list_name> \<json>                    | 用 json 覆盖 list_name 的数据                                                                         | 同上                                                             |
| /list reset \<list_name>                          | 重置 list_name 的数据                                                                                | 同上                                                             |
| /list rm \<list_name>                             | 删除 list_name（删除后原 list_name 不可使用，只能恢复）                                                          | 同上                                                             |
| /list recover \<list_name>                        | 恢复 list_name                                                                                    | 同上                                                             |
| /list glance \<list_name>                         | 快速查询 list_name 里的所有 key                                                                         | 同上或者 list_name 处于 global namespace 中                           |

## List Operation

| Command                   | Description                        | Comment                  |
|---------------------------|------------------------------------|--------------------------|
| /list op \<A> U \<B> \<C> | 并集运算 `A` Union `B` equals `C`      | 需要 namespace admin 或更高权限 |
| /list op \<A> I \<B> \<C> | 交集运算 `A` Intersect `B` equals `C`  | 同上                       |
| /list op \<A> D \<B> \<C> | 差集运算 `A` Difference `B` equals `C` | 同上                       |

## Group

| Command                  | Description                                  | Comment                                                           |
|--------------------------|----------------------------------------------|-------------------------------------------------------------------|
| /group bind \<namespace> | 将当前 group 绑定到 namespace 中（会重置当前 group 的所有配置） | 需要 group admin 和 namespace admin 或更高权限或者 namespace 具有 public 属性   |
| /group unbind            | 解除当前 group 的绑定                               | 同上                                                                |
| /group clone \<group_id> | 从 group_id 克隆配置到当前 group（需要先绑定 namespace）    | 同上                                                                |
| /group query             | 查询当前 group 的配置                               | 需要 namespace admin 或更高权限或者对应 namespace 具有 public 属性               |
| /group query \<group_id> | 查询指定 group 的配置                               | 同上                                                                |
| /group kick \<list_name> | 把当前 group 在 list_name 中的成员踢出                 | 需要 group admin 和 namespace admin 或更高权限或者对应 namespace 具有 public 属性 |
| /group keep \<list_name> | 把当前 group 不在 list_name 中的成员踢出                | 同上                                                                |

## Group Approval

| Command                                      | Description                                       | Comment                                |
|----------------------------------------------|---------------------------------------------------|----------------------------------------|
| /group approval enable mc                    | 入群审核启用 mc 正版用户名验证（将使用正版 UUID 作为双因子认证的输入）          | 需要 group admin 和 namespace admin 或更高权限 |
| /group approval enable regexp                | 入群审核启用正则表达式（将使用匹配结果作为双因子认证的输入）                    | 同上                                     |
| /group approval enable whitelist             | 入群审核启用白名单                                         | 同上                                     |
| /group approval enable blacklist             | 入群审核启用黑名单                                         | 同上                                     |
| /group approval enable notify-only           | 入群审核启用仅通知                                         | 同上                                     |
| /group approval enable auto-pass             | 入群审核启用自动通过（默认启用）                                  | 同上                                     |
| /group approval enable auto-reject           | 入群审核启用自动拒绝（默认启用）                                  | 同上                                     |
| /group approval set regexp \<regexp>         | 指定入群审核的正则表达式（若有子表达式，则会使用第一个子表达式的匹配结果）             | 同上                                     |
| /group approval set notification \<group_id> | 指定入群审核通知群                                         | 同上                                     |
| /group approval add whitelist \<list_name>   | 新增入群审核白名单 list_name（可以多次指定不同的 list_name 最终采用并集查找） | 同上                                     |
| /group approval add blacklist \<list_name>   | 新增入群审核黑名单 list_name（可以多次指定不同的 list_name 最终采用并集查找） | 同上                                     |
| /group approval rm whitelist \<list_name>    | 移除入群审核白名单 list_name                               | 同上                                     |
| /group approval rm blacklist \<list_name>    | 移除入群审核黑名单 list_name                               | 同上                                     |
| /group approval rm notification              | 移除入群审核通知群                                         | 同上                                     |
| /group approval disable mc                   | 入群审核禁用 mc 正版用户名验证                                 | 同上                                     |
| /group approval disable regexp               | 入群审核禁用正则表达式                                       | 同上                                     |
| /group approval disable whitelist            | 入群审核禁用白名单                                         | 同上                                     |
| /group approval disable blacklist            | 入群审核禁用黑名单                                         | 同上                                     |
| /group approval disable notify-only          | 入群审核禁用仅通知                                         | 同上                                     |
| /group approval disable auto-pass            | 入群审核禁用自动通过（言下之意，符合通过条件的申请不自动处理，需要手动处理）            | 同上                                     |
| /group approval disable auto-reject          | 入群审核禁用自动拒绝（言下之意，不符合通过条件的申请不自动处理，需要手动处理）           | 同上                                     |

## Group Keyword

| Command                                   | Description                                        | Comment                                |
|-------------------------------------------|----------------------------------------------------|----------------------------------------|
| /group keyword enable blacklist           | 关键词检查启用黑名单                                         | 需要 group admin 和 namespace admin 或更高权限 |
| /group keyword enable whitelist           | 关键词检查启用白名单                                         | 同上                                     |
| /group keyword add blacklist \<list_name> | 新增关键词检查黑名单 list_name（可以多次指定不同的 list_name 最终采用并集查找） | 同上                                     |
| /group keyword add whitelist \<list_name> | 新增关键词检查白名单 list_name（可以多次指定不同的 list_name 最终采用并集查找） | 同上                                     |
| /group keyword add reply \<list_name>     | 新增关键词回复列表 list_name（可以多次指定不同的 list_name 最终采用并集查找）  | 同上                                     |
| /group keyword rm blacklist \<list_name>  | 移除关键词检查黑名单 list_name                               | 同上                                     |
| /group keyword rm whitelist \<list_name>  | 移除关键词检查白名单 list_name                               | 同上                                     |
| /group keyword rm reply \<list_name>      | 移除关键词回复列表 list_name                                | 同上                                     |
| /group keyword disable blacklist          | 关键词检查禁用黑名单                                         | 同上                                     |
| /group keyword disable whitelist          | 关键词检查禁用白名单                                         | 同上                                     |

## Group Card

| Command                                                | Description                                               | Comment                                |
|--------------------------------------------------------|-----------------------------------------------------------|----------------------------------------|
| /group card check \<list_name> with \<regexp>          | 使用 regexp 正则匹配群名片，将不匹配的成员写入 list_name 中                   | 需要 group admin 和 namespace admin 或更高权限 |
| /group card check \<to_list_name> by \<from_list_name> | 使用 from_list_name uid:card 匹配群名片，将不匹配的成员写入 to_list_name 中 | 同上                                     |
| /group card set auto-set \<list_name>                  | 设置入群自动修改群名片 list_name；若不包含，那么不修改                          | 同上                                     |
| /group card rm auto-set                                | 清除设置的自动修改群名片 list_name                                    | 同上                                     |
| /group card lock                                       | 锁定群名片，不让修改                                                | 同上                                     |
| /group card unlock                                     | 解锁群名片                                                     | 同上                                     |

## Group Message

| Command                                     | Description | Comment                                |
|---------------------------------------------|-------------|----------------------------------------|
| /group message enable anti-recall           | 启用群反撤回      | 需要 group admin 和 namespace admin 或更高权限 |
| /group message disable anti-recall          | 禁用群反撤回      | 同上                                     |
| /group message set notification \<group_id> | 指定群消息通知群    | 同上                                     |
| /group message rm notification              | 移除群消息通知群    | 同上                                     |
| /group message set only-anti-recall-member  | 设置仅反撤回群成员   | 同上                                     |
| /group message rm only-anti-recall-member   | 取消仅反撤回群成员   | 同上                                     |

## Group Log

| Command                              | Description        | Comment                                |
|--------------------------------------|--------------------|----------------------------------------|
| /group log set approval \<list_name> | 设置记录入群审核 list_name | 需要 group admin 和 namespace admin 或更高权限 |
| /group log rm approval               | 移除记录入群审核           | 同上                                     |
| /group log set leave \<list_name>    | 设置记录离群 list_name   | 同上                                     |
| /group log rm leave                  | 移除记录离群             | 同上                                     |

## Group Export

| Command                           | Description                 | Comment                                |
|-----------------------------------|-----------------------------|----------------------------------------|
| /group export member \<list_name> | 导出 group member 到 list_name | 需要 group admin 和 namespace admin 或更高权限 |

## User

| Command                             | Description                      | Comment                    |
|-------------------------------------|----------------------------------|----------------------------|
| /user join \<namespace> \<user_id>  | 将 user_id 添加到 namespace admin 名单 | 需要 namespace owner 权限或更高权限 |
| /user leave \<namespace> \<user_id> | 将 user_id 从 namespace admin 名单移除 | 同上                         |

## Namespace

| Command                                           | Description                      | Comment                                                                 |
|---------------------------------------------------|----------------------------------|-------------------------------------------------------------------------|
| /namespace add \<namespace>                       | 新建 namespace                     | 需要系统授予的操作 namespace 权限                                                  |
| /namespace rm \<namespace>                        | 删除 namespace                     | 需要 namespace owner 权限                                                   |
| /namespace query                                  | 查询自己所有的和具有 public 属性的 namespace  | 需要 namespace admin 或更高权限或者是具有 public 属性的 namespace 或者是 global namespace |
| /namespace \<namespace>                           | 查询 namespace 配置                  | 同上                                                                      |
| /namespace \<namespace> reset admin               | 重置 namespace 的 admin             | 需要 namespace owner 权限或更高权限                                              |
| /namespace chown \<owner_id> \<namespace>         | 修改 namespace 的 owner             | 同上                                                                      |
| /namespace \<namespace> set public \<true\|false> | 设置 namespace 的 public 属性         | 同上                                                                      |
| /namespace \<namespace> load list \<list_name>    | 加载属于public namespace 的 list_name | 需要 namespace admin 或更高权限                                                |
| /namespace \<namespace> unload list \<list_name>  | 卸载属于public namespace 的 list_name | 同上                                                                      |

## Extra

| Command                               | Description        | Comment                  |
|---------------------------------------|--------------------|--------------------------|
| /raw \<message>                       | 获取 message 的原始信息   | 需要系统授予的获取 raw 的权限        |
| /broadcast group <group_id> <content> | 广播消息到群             | 需要 namespace admin 或更高权限 |
| /model set \<model>                   | 设置机型               | 需要受系统信任                  |
| /token add \<name> \<token>           | 添加可让 user 接入本系统的令牌 | 需要系统授予的操作 token 的权限      |
| /token rm \<name>                     | 删除令牌               | 需要系统授予的操作 token 的权限      |
| /token query                          | 查询自己所有的令牌          | 需要系统授予的操作 token 的权限      |
| /token query \<name>                  | 查询令牌               | 需要系统授予的操作 token 的权限      |
| /token chown \<owner_id> \<name>      | 修改令牌 owner         | 需要系统授予的操作 token 的权限      |
| /token bind \<bot_id> \<name>         | 绑定令牌使用的机器人账号       | 需要系统授予的操作 token 的权限      |
| /token unbind \<name>                 | 解绑令牌使用的机器人账号       | 需要系统授予的操作 token 的权限      |

## Advanced features

1. 当 group keyword 使用的 key 有**空格**（转义后为`%20`）时，会对 key 的字符串以**空格**拆分。当检查的文本全部包含拆分字符串，则匹配成功。

   e.g. key: `砍 并夕夕` 文本1：`帮我并夕夕砍一刀`（匹配成功）文本2：`是兄弟就来砍我`（匹配失败）

2. 关键字回复支持 Webhook

   Webhook 支持 `GET`、`POST`、`PUT`、`DELETE` 和 `PATCH` 方法。如果返回的内容是 JSON 数据，可以使用 JSON Path
   提取内容。如果返回的内容是媒体，将会直接发送。

   **协议头**

   * `GET`：`webhook://` 或 `webhook:get://`
   * `POST`：`webhook:post://`
   * `PUT`：`webhook:put://`
   * `DELETE`：`webhook:delete://`
   * `PATCH`：`webhook:patch://`

   使用协议头作为回复内容的开头，将会执行协议头之后的 url。

   **可用的占位符**

   * `{message}` 用户发送的消息
   * `{remain}` 除去 Keyword 剩下的字符串
   * `{nickname}` 用户昵称
   * `{userId}` 用户 ID
   * `{groupId}` 群 ID

   **Method**

   * `GET`：

     template: `webhook#headers#@response.json.path@://url`

     `#headers#` 接受换行输入；可以省略。

     `@response.json.path@` 可以省略。

     e.g. `/list join example get webhook@data.text@://https://example.com/{groupId}/{userId}/{message}`

   * `POST`、`PUT`、`DELETE`、`PATCH`：

     template: `webhook:post#headers#<request body>@response.json.path@://url`

     `#headers#` 接受换行输入；可以省略。

     `<request body>` 接受换行输入；可以省略。

     `@response.json.path@` 可以省略。

     e.g. `/list join example post webhook:post#Authorization: Bearer\nX-Id: {userId}#<{"group_id":{groupId},"user_id":{userId},"message":{message},"remain":{remain}}>@data.text@://https://example.com/{groupId}/{userId}/{message}`

   **注意**：要触发 Webhook，消息必须以 Keyword 开头。

   例如：`bb` 是 Keyword。用户发送 `bb 123`，触发 Webhook；用户发送 `aa bb`，不会触发 Webhook。

3. 关键字回复支持 Commands 命令组合

   **协议头**

   `command://` 或 `cmd://`

   使用协议头作为回复内容的开头，将会执行协议头之后的 Commands。

   多个命令用 `&&` 分隔。**注意！**`&&`前后必须各有一个*空格*！除非`&&`后紧跟换行。

   **可用的占位符**

   * `{message}` 用户发送的消息
   * `{remain}` 除去 Keyword 剩下的字符串

   若需要多层占位符，请手动换成 `{<key><number>}` 格式。程序会查找所有符合 `{<key><number>}` 格式的占位符，并将其中的数字递减
   1。当数字递减到 0 时，移除数字后缀。

   e.g. `/list join example cmd cmd:///list glance example && /list len example`

   **注意**：要触发 Commands 命令组合，消息必须以 Keyword 开头。

4. 关键字回复支持 Rewrite 消息

   **协议头**

   `rewrite://`

   使用协议头作为回复内容的开头，将会执行协议头之后的重写。

   多个重写用 `&` 分隔。**注意！**`&`前后必须各有一个*空格*！除非`&`后紧跟换行。

   **可用的占位符**

   * `{message}` 用户发送的消息
   * `{remain}` 除去 Keyword 剩下的字符串

   e.g. `/list join example rewrite rewrite://{remain} & {message}`

   **注意**：要触发 Rewrite，消息必须以 Keyword 开头。

   例如：`bb` 是 Keyword。用户发送 `bb 123`，触发 Rewrite；用户发送 `aa bb`，不会触发 Rewrite。