package generator

import (
	"github.com/google/uuid"
	"github.com/wisle25/be-template/applications/generator"
)

type UUIDGenerator struct {
}

func NewUUIDGenerator() generator.IdGenerator {
	return &UUIDGenerator{}
}

func (generator *UUIDGenerator) Generate() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	return id.String()
}
