// Copyright 2026 SCION Association
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topology

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseInternalAddrs(t *testing.T) {
	v4 := "10.0.0.1:30042"
	v6 := "[2001:db8::1]:30042"

	t.Run("empty errors", func(t *testing.T) {
		_, err := parseInternalAddrs("br", nil)
		assert.Error(t, err)
	})
	t.Run("single v4", func(t *testing.T) {
		got, err := parseInternalAddrs("br", []string{v4})
		require.NoError(t, err)
		assert.Equal(t, []netip.AddrPort{netip.MustParseAddrPort(v4)}, got)
	})
	t.Run("dual stack canonicalized v6 first", func(t *testing.T) {
		got, err := parseInternalAddrs("br", []string{v4, v6})
		require.NoError(t, err)
		// Input order is v4 then v6, but the canonical order is v6 first.
		require.Len(t, got, 2)
		assert.True(t, got[0].Addr().Is6() && !got[0].Addr().Is4In6())
		assert.True(t, got[1].Addr().Is4())
	})
	t.Run("duplicate errors", func(t *testing.T) {
		_, err := parseInternalAddrs("br", []string{v4, v4})
		assert.Error(t, err)
	})
	t.Run("same family same port errors", func(t *testing.T) {
		_, err := parseInternalAddrs("br", []string{"10.0.0.1:30042", "10.0.0.2:30042"})
		assert.Error(t, err)
	})
	t.Run("different family same port allowed", func(t *testing.T) {
		got, err := parseInternalAddrs("br", []string{v4, v6})
		require.NoError(t, err)
		assert.Len(t, got, 2)
	})
}

func TestSameInternalAddrSet(t *testing.T) {
	a := netip.MustParseAddrPort("10.0.0.1:1")
	b := netip.MustParseAddrPort("[2001:db8::1]:1")
	assert.True(t, sameInternalAddrSet(
		[]netip.AddrPort{a, b}, []netip.AddrPort{b, a}), "order should not matter")
	assert.False(t, sameInternalAddrSet(
		[]netip.AddrPort{a}, []netip.AddrPort{a, b}), "different length differs")
	assert.False(t, sameInternalAddrSet(
		[]netip.AddrPort{a}, []netip.AddrPort{b}), "different element differs")
}
