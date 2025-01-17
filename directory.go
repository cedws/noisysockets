// SPDX-License-Identifier: MPL-2.0
/*
 * Copyright (C) 2024 Damian Peckett <damian@pecke.tt>.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

package noisysockets

import (
	"fmt"
	"net/netip"

	"github.com/noisysockets/noisysockets/types"
)

type peerDirectory struct {
	peerNames       map[string]types.NoisePublicKey
	peerAddresses   map[types.NoisePublicKey][]netip.Addr
	fromPeerAddress map[netip.Addr]types.NoisePublicKey
}

func newPeerDirectory() *peerDirectory {
	return &peerDirectory{
		peerNames:       make(map[string]types.NoisePublicKey),
		peerAddresses:   make(map[types.NoisePublicKey][]netip.Addr),
		fromPeerAddress: make(map[netip.Addr]types.NoisePublicKey),
	}
}

func (pd *peerDirectory) AddPeer(name string, publicKey types.NoisePublicKey, addrs []netip.Addr) error {
	if name != "" {
		pd.peerNames[name] = publicKey
	}
	pd.peerAddresses[publicKey] = addrs
	for _, addr := range addrs {
		if _, ok := pd.fromPeerAddress[addr]; ok {
			return fmt.Errorf("address %s already in use", addr)
		}

		pd.fromPeerAddress[addr] = publicKey
	}

	return nil
}

func (pd *peerDirectory) LookupPeerAddressesByName(name string) ([]netip.Addr, bool) {
	publicKey, ok := pd.peerNames[name]
	if !ok {
		return nil, false
	}
	addrs, ok := pd.peerAddresses[publicKey]
	return addrs, ok
}

func (pd *peerDirectory) LookupPeerByAddress(addr netip.Addr) (types.NoisePublicKey, bool) {
	publicKey, ok := pd.fromPeerAddress[addr]
	return publicKey, ok
}
