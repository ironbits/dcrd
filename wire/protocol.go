// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// ProtocolVersion is the latest protocol version this package supports.
	ProtocolVersion uint32 = 1
)

// ServiceFlag identifies services supported by a decred peer.
type ServiceFlag uint64

const (
	// SFNodeNetwork is a flag used to indicate a peer is a full node.
	SFNodeNetwork ServiceFlag = 1 << iota
)

// Map of service flags back to their constant names for pretty printing.
var sfStrings = map[ServiceFlag]string{
	SFNodeNetwork: "SFNodeNetwork",
}

// String returns the ServiceFlag in human-readable form.
func (f ServiceFlag) String() string {
	// No flags are set.
	if f == 0 {
		return "0x0"
	}

	// Add individual bit flags.
	s := ""
	for flag, name := range sfStrings {
		if f&flag == flag {
			s += name + "|"
			f -= flag
		}
	}

	// Add any remaining flags which aren't accounted for as hex.
	s = strings.TrimRight(s, "|")
	if f != 0 {
		s += "|0x" + strconv.FormatUint(uint64(f), 16)
	}
	s = strings.TrimLeft(s, "|")
	return s
}

// CurrencyNet represents which decred network a message belongs to.
type CurrencyNet uint32

// Constants used to indicate the message decred network.  They can also be
// used to seek to the next message when a stream's state is unknown, but
// this package does not provide that functionality since it's generally a
// better idea to simply disconnect clients that are misbehaving over TCP.
const (
	// MainNet represents the main decred network.
	MainNet CurrencyNet = 0xd9b400f9

	// TestNet represents the regression test network.
	RegTest CurrencyNet = 0xdab500fa

	// TestNet represents the test network (version 3).
	TestNet CurrencyNet = 0x0709000b

	// SimNet represents the simulation test network.
	SimNet CurrencyNet = 0x12141c16
)

// bnStrings is a map of decred networks back to their constant names for
// pretty printing.
var bnStrings = map[CurrencyNet]string{
	MainNet: "MainNet",
	TestNet: "TestNet",
	RegTest: "RegNet",
	SimNet:  "SimNet",
}

// String returns the CurrencyNet in human-readable form.
func (n CurrencyNet) String() string {
	if s, ok := bnStrings[n]; ok {
		return s
	}

	return fmt.Sprintf("Unknown CurrencyNet (%d)", uint32(n))
}
