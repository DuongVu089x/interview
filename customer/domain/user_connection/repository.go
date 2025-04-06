package userconnection

type Repository interface {
	GetUserConnection(query *UserConnection) (*UserConnection, error)
	GetUserConnections(query *UserConnection, offset, limit int64) ([]*UserConnection, error)
	CreateUserConnection(userConn *UserConnection) (*UserConnection, error)
	UpdateUserConnection(query *UserConnection, updating *UserConnection) error
	DeleteUserConnection(userConn *UserConnection) error
}
