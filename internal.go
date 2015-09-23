// +build !windows

package rbxclip

func notImplemented() {
	panic("not implemented")
}

func clear() (err error) {
	notImplemented()
	return
}

func has() (available bool, err error) {
	notImplemented()
	return
}

func get() (b []byte, err error) {
	notImplemented()
	return
}

func set(b []byte) (err error) {
	notImplemented()
	return
}
