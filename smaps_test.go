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
	Dump0LCIdentifier   Identifier = "/usr/lib/locale/C.utf8/LC_CTYPE"
)

var (
	Dump0BashAddresses = []string{"55f3e5c55000-55f3e5c84000", "55f3e5c84000-55f3e5d6e000", "55f3e5d6e000-55f3e5da3000", "55f3e5da3000-55f3e5da7000", "55f3e5da7000-55f3e5db0000"}
)

func TestParseDump0(t *testing.T) {
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

	assert.Len(t, smaps[Dump0LCIdentifier], 1)
	assert.Equal(t, smaps[Dump0LCIdentifier][0].Size, int64(352*1024))
	assert.Equal(t, smaps[Dump0LCIdentifier][0].Rss, int64(112*1024))
	assert.Equal(t, smaps[Dump0LCIdentifier][0].SharedClean, int64(112*1024))
}

func TestParseDump1(t *testing.T) {
	file, err := os.Open("testdata/dump1.txt")
	require.NoError(t, err)

	smaps, err := Parse(file)
	require.NoError(t, err)

	assert.Equal(t, 346, len(smaps))
}
