package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestApproximatePow(t *testing.T) {
	// Set a common precision for the tests
	precision := "0.000000001"

	// Test cases
	tests := []struct {
		base     string
		exponent string
		expected string
	}{
		{"1.45", "1.5", "1.746031213"},
		{"1.05", "3.5", "1.186212638"},
		{"0.91", "1.85", "0.839898055"},
	}

	for _, tt := range tests {
		base := types.MustNewDecFromStr(tt.base)
		exponent := types.MustNewDecFromStr(tt.exponent)
		expected := types.MustNewDecFromStr(tt.expected)

		result, err := ApproximatePow(base.String(), exponent.String(), precision)
		require.NoError(t, err)
		resultAsDec := types.MustNewDecFromStr(result.String())
		require.Equal(t, expected, resultAsDec)
	}
}
