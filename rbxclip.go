// Package rbxclip reads and writes Roblox instances to and from the
// clipboard. Data written to the clipboard is compatible with Roblox Studio.
package rbxclip

import (
	"bytes"
	"github.com/robloxapi/rbxfile"
	"github.com/robloxapi/rbxfile/bin"
)

// Get reads instances from the clipboard, decoding them to a rbxfile.Root. A
// nil value is returned if the clipboard contains no instances, or is unable
// to read them.
func Get() *rbxfile.Root {
	b, err := get()
	if err != nil {
		return nil
	}
	buf := bytes.NewBuffer(b)
	root, err := bin.DeserializeModel(buf, nil)
	if err != nil {
		return nil
	}
	return root
}

// Has returns whether or not there are instances in the clipboard. Returns
// false if the call fails.
func Has() bool {
	b, _ := has()
	return b
}

// Set writes instances to the clipboard, encoding from a rbxfile.Root.
// Passing nil will simply clear the clipboard. Returns whether or not the
// call was successful.
func Set(root *rbxfile.Root) bool {
	if root == nil {
		return clear() == nil
	}
	var buf bytes.Buffer
	if err := bin.SerializeModel(&buf, nil, root); err != nil {
		return false
	}
	clear()
	return set(buf.Bytes()) == nil
}
