// SPDX-License-Identifier: Apache-2.0

package smaps

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	Dump0BashIdentifier Identifier = "/usr/bin/bash"
)

var (
	Dump0BashAddresses = []string{"55f3e5c55000-55f3e5c84000", "55f3e5c84000-55f3e5d6e000", "55f3e5d6e000-55f3e5da3000", "55f3e5da3000-55f3e5da7000", "55f3e5da7000-55f3e5db0000"}
)

func TestParse(t *testing.T) {
	file, err := os.Open("testdata/dump0.txt")
	require.NoError(t, err)

	smaps, err := Parse(file)
	require.NoError(t, err)

	assert.Len(t, smaps[UnknownIdentifier], 5)
	assert.Len(t, smaps[Dump0BashIdentifier], 5)

	for i, smap := range smaps[Dump0BashIdentifier] {
		assert.Equal(t, Dump0BashIdentifier, smap.Identifier)
		assert.Equal(t, Dump0BashAddresses[i], smap.Address)
		assert.True(t, smap.IsPath)
	}
}
