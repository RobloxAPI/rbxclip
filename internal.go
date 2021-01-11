// +build !windows

package rbxclip

import "errors"

var notImplemented = errors.New("not implemented")

func clear() (err error) {
	return notImplemented
}

func has() (available bool, err error) {
	return false, notImplemented
}

func get() (b []byte, err error) {
	return nil, notImplemented
}

func set(b []byte) (err error) {
	return notImplemented
}
