// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/decred/dcrd/wire"
)

// TestPongLatest tests the MsgPong API against the latest protocol version.
func TestPongLatest(t *testing.T) {
	pver := wire.ProtocolVersion

	nonce, err := wire.RandomUint64()
	if err != nil {
		t.Errorf("RandomUint64: error generating nonce: %v", err)
	}
	msg := wire.NewMsgPong(nonce)
	if msg.Nonce != nonce {
		t.Errorf("NewMsgPong: wrong nonce - got %v, want %v",
			msg.Nonce, nonce)
	}

	// Ensure the command is expected value.
	wantCmd := "pong"
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgPong: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value for latest protocol version.
	wantPayload := uint32(8)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload)
	}

	// Test encode with latest protocol version.
	var buf bytes.Buffer
	err = msg.BtcEncode(&buf, pver)
	if err != nil {
		t.Errorf("encode of MsgPong failed %v err <%v>", msg, err)
	}

	// Test decode with latest protocol version.
	readmsg := wire.NewMsgPong(0)
	err = readmsg.BtcDecode(&buf, pver)
	if err != nil {
		t.Errorf("decode of MsgPong failed [%v] err <%v>", buf, err)
	}

	// Ensure nonce is the same.
	if msg.Nonce != readmsg.Nonce {
		t.Errorf("Should get same nonce for protocol version %d", pver)
	}

	return
}

// TestPongWire tests the MsgPong wire encode and decode for various protocol
// versions.
func TestPongWire(t *testing.T) {
	tests := []struct {
		in   wire.MsgPong // Message to encode
		out  wire.MsgPong // Expected decoded message
		buf  []byte       // Wire encoding
		pver uint32       // Protocol version for wire encoding
	}{
		// Latest protocol version.
		{
			wire.MsgPong{Nonce: 123123}, // 0x1e0f3
			wire.MsgPong{Nonce: 123123}, // 0x1e0f3
			[]byte{0xf3, 0xe0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00},
			wire.ProtocolVersion,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode the message to wire format.
		var buf bytes.Buffer
		err := test.in.BtcEncode(&buf, test.pver)
		if err != nil {
			t.Errorf("BtcEncode #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("BtcEncode #%d\n got: %s want: %s", i,
				spew.Sdump(buf.Bytes()), spew.Sdump(test.buf))
			continue
		}

		// Decode the message from wire format.
		var msg wire.MsgPong
		rbuf := bytes.NewReader(test.buf)
		err = msg.BtcDecode(rbuf, test.pver)
		if err != nil {
			t.Errorf("BtcDecode #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(msg, test.out) {
			t.Errorf("BtcDecode #%d\n got: %s want: %s", i,
				spew.Sdump(msg), spew.Sdump(test.out))
			continue
		}
	}
}

// TestPongWireErrors performs negative tests against wire encode and decode
// of MsgPong to confirm error paths work correctly.
func TestPongWireErrors(t *testing.T) {
	pver := wire.ProtocolVersion

	basePong := wire.NewMsgPong(123123) // 0x1e0f3
	basePongEncoded := []byte{
		0xf3, 0xe0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	tests := []struct {
		in       *wire.MsgPong // Value to encode
		buf      []byte        // Wire encoding
		pver     uint32        // Protocol version for wire encoding
		max      int           // Max size of fixed buffer to induce errors
		writeErr error         // Expected write error
		readErr  error         // Expected read error
	}{
		// Latest protocol version with intentional read/write errors.
		// Force error in nonce.
		{basePong, basePongEncoded, pver, 0, io.ErrShortWrite, io.EOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to wire format.
		w := newFixedWriter(test.max)
		err := test.in.BtcEncode(w, test.pver)
		if reflect.TypeOf(err) != reflect.TypeOf(test.writeErr) {
			t.Errorf("BtcEncode #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// For errors which are not of type wire.MessageError, check
		// them for equality.
		if _, ok := err.(*wire.MessageError); !ok {
			if err != test.writeErr {
				t.Errorf("BtcEncode #%d wrong error got: %v, "+
					"want: %v", i, err, test.writeErr)
				continue
			}
		}

		// Decode from wire format.
		var msg wire.MsgPong
		r := newFixedReader(test.max, test.buf)
		err = msg.BtcDecode(r, test.pver)
		if reflect.TypeOf(err) != reflect.TypeOf(test.readErr) {
			t.Errorf("BtcDecode #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}

		// For errors which are not of type wire.MessageError, check
		// them for equality.
		if _, ok := err.(*wire.MessageError); !ok {
			if err != test.readErr {
				t.Errorf("BtcDecode #%d wrong error got: %v, "+
					"want: %v", i, err, test.readErr)
				continue
			}
		}

	}
}
