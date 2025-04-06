package idgen

type Repository interface {
	GenerateID(string) (int64, error)
}
