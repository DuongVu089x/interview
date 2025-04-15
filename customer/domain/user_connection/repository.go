package userconnection

import "context"

type Repository interface {
	GetUserConnection(ctx context.Context, query *UserConnection) (*UserConnection, error)
	GetUserConnections(ctx context.Context, query *UserConnection, offset, limit int64) ([]*UserConnection, error)
	CreateUserConnection(ctx context.Context, userConn *UserConnection) (*UserConnection, error)
	UpdateUserConnection(ctx context.Context, query *UserConnection, updating *UserConnection) error
	DeleteUserConnection(ctx context.Context, userConn *UserConnection) error
}
