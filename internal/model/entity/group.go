// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Group is the golang structure for table group.
type Group struct {
	GroupId     int64       `json:"group_id"     ` //
	Namespace   string      `json:"namespace"    ` //
	SettingJson string      `json:"setting_json" ` //
	CreatedAt   *gtime.Time `json:"created_at"   ` //
	UpdatedAt   *gtime.Time `json:"updated_at"   ` //
	DeletedAt   *gtime.Time `json:"deleted_at"   ` //
}