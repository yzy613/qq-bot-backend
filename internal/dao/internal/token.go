// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TokenDao is the data access object for table token.
type TokenDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of current DAO.
	columns TokenColumns // columns contains all the column names of Table for convenient usage.
}

// TokenColumns defines and stores column names for table token.
type TokenColumns struct {
	Name         string //
	Token        string //
	OwnerId      string //
	CreatedAt    string //
	UpdatedAt    string //
	DeletedAt    string //
	LastLoginAt  string //
	BindingBotId string //
}

// tokenColumns holds the columns for table token.
var tokenColumns = TokenColumns{
	Name:         "name",
	Token:        "token",
	OwnerId:      "owner_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
	LastLoginAt:  "last_login_at",
	BindingBotId: "binding_bot_id",
}

// NewTokenDao creates and returns a new DAO object for table data access.
func NewTokenDao() *TokenDao {
	return &TokenDao{
		group:   "default",
		table:   "token",
		columns: tokenColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TokenDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TokenDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TokenDao) Columns() TokenColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TokenDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TokenDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TokenDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
