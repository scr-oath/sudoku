package concur

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadInput(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		validate  func(*testing.T, [][]uint64)
	}{
		{
			name: "given a valid inputFile, " +
				"when ReadInput is invoked, " +
				"then no errors occur, and the data is populated",
			inputFile: "../input.txt",
			validate: func(t *testing.T, i [][]uint64) {
				var nonEmptyCount int
				for _, arr := range i {
					for _, elem := range arr {
						if elem != 0 {
							nonEmptyCount++
						}
					}
				}
				assert.Greater(t, nonEmptyCount, 0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotEmpty(t, tt.inputFile)

			got := ReadInput(tt.inputFile)
			require.NotNil(t, got)

			require.NotNil(t, tt.validate)
			tt.validate(t, got)
		})
	}
}
