package generator

import (
	"github.com/google/uuid"
	"github.com/wisle25/be-template/applications/generator"
	"github.com/wisle25/be-template/commons"
)

// UUIDGenerator implements IdGenerator using UUID
type UUIDGenerator struct /* implements IdGenerator */ {

}

func NewUUIDGenerator() generator.IdGenerator {
	return &UUIDGenerator{}
}

func (generator *UUIDGenerator) Generate() string {
	id, err := uuid.NewV7()
	if err != nil {
		commons.ThrowServerError("id_generator_err: generate", err)
	}

	return id.String()
}
