package idgen

import (
	"math"

	domainidgen "github.com/DuongVu089x/interview/order/domain/id_gen"
)

type IDGenService struct {
	idGenRepository domainidgen.Repository
}

func NewIDGenService(idGenRepository domainidgen.Repository) domainidgen.Service {
	return &IDGenService{
		idGenRepository: idGenRepository,
	}
}

// convertToCode convert id from int to string
func convertToCode(number int64, length int64, template string) string {
	var result = ""
	var i = int64(0)
	var ln = int64(len(template))
	var capacity = int64(math.Pow(float64(ln), float64(length)))
	number = number % capacity
	for i < length {
		var cur = number % ln
		if i > 0 {
			cur = (cur + int64(result[i-1])) % ln
		}
		result = result + string(template[cur])
		number = number / ln
		i++
	}
	return result
}

func (s *IDGenService) GenerateID(key string) (int64, string, error) {
	id, err := s.idGenRepository.GenerateID(key)
	if err != nil {
		return 0, "", err
	}
	return id, convertToCode(id, 10, "0123456789"), nil
}
