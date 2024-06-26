// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GroupDao is the data access object for table group.
type GroupDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of current DAO.
	columns GroupColumns // columns contains all the column names of Table for convenient usage.
}

// GroupColumns defines and stores column names for table group.
type GroupColumns struct {
	GroupId     string //
	Namespace   string //
	SettingJson string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// groupColumns holds the columns for table group.
var groupColumns = GroupColumns{
	GroupId:     "group_id",
	Namespace:   "namespace",
	SettingJson: "setting_json",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewGroupDao creates and returns a new DAO object for table data access.
func NewGroupDao() *GroupDao {
	return &GroupDao{
		group:   "default",
		table:   "group",
		columns: groupColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *GroupDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *GroupDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *GroupDao) Columns() GroupColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *GroupDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *GroupDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *GroupDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
