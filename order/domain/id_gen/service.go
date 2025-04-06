package idgen

type Service interface {
	GenerateID(string) (int64, string, error)
}
