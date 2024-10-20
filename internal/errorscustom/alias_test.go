package errorscustom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "ErrConflict",
			err:  ErrConflict,
		},
		{
			name: "ErrUserIDNotContext",
			err:  ErrUserIDNotContext,
		},
		{
			name: "ErrDeletedURL",
			err:  ErrDeletedURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.err)
		})
	}
}
