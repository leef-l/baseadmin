package service

import (
	"context"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IAuth interface {
	Login(ctx context.Context, in *model.AuthLoginInput) (out *model.AuthLoginOutput, err error)
	TicketLogin(ctx context.Context, in *model.AuthTicketLoginInput) (out *model.AuthLoginOutput, err error)
	IssueTicket(ctx context.Context, in *model.AuthIssueTicketInput) (out *model.AuthIssueTicketOutput, err error)
	Info(ctx context.Context, userID snowflake.JsonInt64) (out *model.AuthInfoOutput, err error)
	ChangePassword(ctx context.Context, in *model.AuthChangePasswordInput) error
	Menus(ctx context.Context, userID snowflake.JsonInt64) ([]*model.AuthMenuOutput, error)
}

var localAuth IAuth

func Auth() IAuth {
	return localAuth
}

func RegisterAuth(i IAuth) {
	localAuth = i
}
