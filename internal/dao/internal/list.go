// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ListDao is the data access object for the table list.
type ListDao struct {
	table   string      // table is the underlying table name of the DAO.
	group   string      // group is the database configuration group name of the current DAO.
	columns ListColumns // columns contains all the column names of Table for convenient usage.
}

// ListColumns defines and stores column names for the table list.
type ListColumns struct {
	ListName  string //
	Namespace string //
	ListJson  string //
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
}

// listColumns holds the columns for the table list.
var listColumns = ListColumns{
	ListName:  "list_name",
	Namespace: "namespace",
	ListJson:  "list_json",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewListDao creates and returns a new DAO object for table data access.
func NewListDao() *ListDao {
	return &ListDao{
		group:   "default",
		table:   "list",
		columns: listColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ListDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ListDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ListDao) Columns() ListColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ListDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ListDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *ListDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
