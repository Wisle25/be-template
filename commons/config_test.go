package commons_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/commons"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	assert.NotPanics(t, func() {
		commons.LoadConfig("../")

	})
}
