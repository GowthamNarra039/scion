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

package control

import (
	"net/netip"

	"github.com/scionproto/scion/pkg/private/serrors"
)

func isV6(ap netip.AddrPort) bool{
	return ap.Addr().Is6() && !ap.Addr().Is4In6()
}

//pickSiblingPair selects the (local, remote) internal address pair to use for
// a sibling link between the local router and the owning (sibling) router.
//
// A sibling link is an AS-internal IP underlay connection, so both ends must
// belong to the same address family. When both routers are dual-stack and offer
// both families, IPv6 is preferred. Within the chosen family the first address
// on each side (in canonical, IPv6-first order) is used, which lets operators
// express a within-family preference by ordering their address lists.
//
// If the two address lists share no common family, the AS is misconfigured for
// these two routers (e.g. one is IPv6-only and the other IPv4-only). In that
// case an error is returned so the router fails to start rather than silently
// losing reachability to the sibling.

func pickSiblingPair(local, remote []netip.AddrPort) (netip.AddrPort, netip.AddrPort, error) {
	if len(local) == 0 {
		return netip.AddrPort{}, netip.AddrPort{},
			serrors.New("local router has no internal address")
	}
	if len(remote) == 0 {
		return netip.AddrPort{}, netip.AddrPort{},
			serrors.New("sibling router has no internal address")
	}

	firstOfFamily := func(addrs []netip.AddrPort, v6 bool) (netip.AddrPort, bool) {
		for _, ap := range addrs {
			if isV6(ap) == v6 {
				return ap, true
			}
		}
		return netip.AddrPort{}, false
	}

	// Prefer IPv6, then fall back to IPv4.
	for _, v6 := range []bool{true, false} {
		l, lok := firstOfFamily(local, v6)
		r, rok := firstOfFamily(remote, v6)
		if lok && rok {
			return l, r, nil
		}
	}

	return netip.AddrPort{}, netip.AddrPort{}, serrors.New(
		"no common address family between local and sibling internal addresses")
}

