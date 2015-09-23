package rbxclip_test

import (
	"bytes"
	"encoding/hex"
	"github.com/robloxapi/rbxclip"
	"github.com/robloxapi/rbxfile/bin"
	"testing"
)

const inputModel = `` +
	`3c726f626c6f782189ff0d0a1a0a000002000000020000000000000000000000` +
	`494e53541b0000001900000000000000f00a0000000008000000496e7456616c` +
	`7565000100000000000002494e53541e0000001c00000000000000f00d010000` +
	`000b000000537472696e6756616c756500010000000000000050524f50210000` +
	`001f00000000000000f01000000000040000004e616d65010e000000496e7420` +
	`56616c7565205465737450524f50140000001200000000000000f00300000000` +
	`0500000056616c7565030000005450524f50240000002200000000000000f013` +
	`01000000040000004e616d650111000000537472696e672056616c7565205465` +
	`737450524f501f0000001d00000000000000f00e010000000500000056616c75` +
	`65010b0000007465737420737472696e6750524e541000000015000000000000` +
	`0035000200010090020000000000000100454e44000000000009000000000000` +
	`003c2f726f626c6f783e`

func TestFunctions(t *testing.T) {
	inputBytes, _ := hex.DecodeString(inputModel)
	inputBuf := bytes.NewBuffer(inputBytes)
	inputRoot, err := bin.DeserializeModel(inputBuf, nil)
	if err != nil {
		t.Fatalf("rbxfile/bin failed to decode file: %s", err)
	}
	if !rbxclip.Set(nil) {
		t.Error("clear clipboard failed")
	}
	if rbxclip.Has() {
		t.Error("Has: expected no instances in clipboard")
	}
	if rbxclip.Get() != nil {
		t.Error("Get: expected no root")
	}
	if !rbxclip.Set(inputRoot) {
		t.Error("set clipboard failed")
	}
	if !rbxclip.Has() {
		t.Error("Has: expected instances in clipboard")
	}
	root := rbxclip.Get()
	if root == nil {
		t.Fatal("Get: expected root")
	}
	var buf bytes.Buffer
	if err := bin.SerializeModel(&buf, nil, root); err != nil {
		t.Fatal("failed to serialize output")
	}
	if !bytes.Equal(inputBytes, buf.Bytes()) {
		t.Fatal("output bytes do not equal input bytes")
	}
}
