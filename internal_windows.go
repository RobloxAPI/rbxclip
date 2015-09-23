// +build windows

package rbxclip

import (
	"errors"
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

const formatName = `application/x-roblox-studio`

var formatID uint32

var (
	libkernel32 uintptr

	getLastError uintptr
	globalAlloc  uintptr
	globalFree   uintptr
	globalLock   uintptr
	globalUnlock uintptr
	globalSize   uintptr
	moveMemory   uintptr
)

var (
	libuser32 uintptr

	closeClipboard             uintptr
	emptyClipboard             uintptr
	getClipboardData           uintptr
	isClipboardFormatAvailable uintptr
	openClipboard              uintptr
	registerClipboardFormat    uintptr
	setClipboardData           uintptr
)

func mustLoadLibrary(name string) uintptr {
	lib, err := syscall.LoadLibrary(name)
	if err != nil {
		panic(err)
	}

	return uintptr(lib)
}

func mustGetProcAddress(lib uintptr, name string) uintptr {
	addr, err := syscall.GetProcAddress(syscall.Handle(lib), name)
	if err != nil {
		panic(err)
	}

	return uintptr(addr)
}

func init() {
	runtime.LockOSThread()

	libkernel32 = mustLoadLibrary("kernel32.dll")
	getLastError = mustGetProcAddress(libkernel32, "GetLastError")
	globalAlloc = mustGetProcAddress(libkernel32, "GlobalAlloc")
	globalFree = mustGetProcAddress(libkernel32, "GlobalFree")
	globalLock = mustGetProcAddress(libkernel32, "GlobalLock")
	globalUnlock = mustGetProcAddress(libkernel32, "GlobalUnlock")
	globalSize = mustGetProcAddress(libkernel32, "GlobalSize")
	moveMemory = mustGetProcAddress(libkernel32, "RtlMoveMemory")

	libuser32 = mustLoadLibrary("user32.dll")
	closeClipboard = mustGetProcAddress(libuser32, "CloseClipboard")
	emptyClipboard = mustGetProcAddress(libuser32, "EmptyClipboard")
	getClipboardData = mustGetProcAddress(libuser32, "GetClipboardData")
	isClipboardFormatAvailable = mustGetProcAddress(libuser32, "IsClipboardFormatAvailable")
	openClipboard = mustGetProcAddress(libuser32, "OpenClipboard")
	registerClipboardFormat = mustGetProcAddress(libuser32, "RegisterClipboardFormatW")
	setClipboardData = mustGetProcAddress(libuser32, "SetClipboardData")

	formatID = win.RegisterClipboardFormat(syscall.StringToUTF16Ptr(formatName))
}

////////////////////////////////////////////////////////////////

type (
	_HANDLE  uintptr
	_HGLOBAL _HANDLE
	_HWND    _HANDLE
)

type winapi struct{}

var win winapi

func (winapi) GetLastError() uint32 {
	ret, _, _ := syscall.Syscall(getLastError, 0,
		0,
		0,
		0)

	return uint32(ret)
}

func lastError(win32FuncName string) error {
	if errno := win.GetLastError(); errno != 0 {
		return errors.New(fmt.Sprintf("%s: Error %d", win32FuncName, errno))
	}
	return errors.New(win32FuncName)
}

func (winapi) GlobalLock(hMem _HGLOBAL) unsafe.Pointer {
	ret, _, _ := syscall.Syscall(globalLock, 1,
		uintptr(hMem),
		0,
		0)

	return unsafe.Pointer(ret)
}

func (winapi) GlobalUnlock(hMem _HGLOBAL) bool {
	ret, _, _ := syscall.Syscall(globalUnlock, 1,
		uintptr(hMem),
		0,
		0)

	return ret != 0
}

func (winapi) GlobalAlloc(uFlags uint32, dwBytes uintptr) _HGLOBAL {
	ret, _, _ := syscall.Syscall(globalAlloc, 2,
		uintptr(uFlags),
		dwBytes,
		0)

	return _HGLOBAL(ret)
}

func (winapi) GlobalFree(hMem _HGLOBAL) _HGLOBAL {
	ret, _, _ := syscall.Syscall(globalFree, 1,
		uintptr(hMem),
		0,
		0)

	return _HGLOBAL(ret)
}

func (winapi) GlobalSize(hMem _HGLOBAL) int {
	ret, _, _ := syscall.Syscall(globalSize, 1,
		uintptr(hMem),
		0,
		0)

	return int(ret)
}

func (winapi) MoveMemory(destination, source unsafe.Pointer, length uintptr) {
	syscall.Syscall(moveMemory, 3,
		uintptr(unsafe.Pointer(destination)),
		uintptr(source),
		uintptr(length))
}

func (winapi) CloseClipboard() bool {
	ret, _, _ := syscall.Syscall(closeClipboard, 0,
		0,
		0,
		0)

	return ret != 0
}

func (winapi) EmptyClipboard() bool {
	ret, _, _ := syscall.Syscall(emptyClipboard, 0,
		0,
		0,
		0)

	return ret != 0
}

func (winapi) GetClipboardData(uFormat uint32) _HANDLE {
	ret, _, _ := syscall.Syscall(getClipboardData, 1,
		uintptr(uFormat),
		0,
		0)

	return _HANDLE(ret)
}

func (winapi) IsClipboardFormatAvailable(format uint32) bool {
	ret, _, _ := syscall.Syscall(isClipboardFormatAvailable, 1,
		uintptr(format),
		0,
		0)

	return ret != 0
}

func (winapi) OpenClipboard(hWndNewOwner _HWND) bool {
	ret, _, _ := syscall.Syscall(openClipboard, 1,
		uintptr(hWndNewOwner),
		0,
		0)

	return ret != 0
}

func (winapi) RegisterClipboardFormat(lpszFormat *uint16) uint32 {
	ret, _, _ := syscall.Syscall(registerClipboardFormat, 1,
		uintptr(unsafe.Pointer(lpszFormat)),
		0,
		0)

	return uint32(ret)
}

func (winapi) SetClipboardData(uFormat uint32, hMem _HANDLE) _HANDLE {
	ret, _, _ := syscall.Syscall(setClipboardData, 2,
		uintptr(uFormat),
		uintptr(hMem),
		0)

	return _HANDLE(ret)
}

////////////////////////////////////////////////////////////////

func clear() error {
	if !win.OpenClipboard(0) {
		return lastError("OpenClipboard")
	}
	defer win.CloseClipboard()

	if !win.EmptyClipboard() {
		return lastError("EmptyClipboard")
	}
	return nil
}

func has() (available bool, err error) {
	if !win.OpenClipboard(0) {
		err = lastError("OpenClipboard")
		return
	}
	defer win.CloseClipboard()

	available = win.IsClipboardFormatAvailable(formatID)
	return
}

func get() (b []byte, err error) {
	if !win.OpenClipboard(0) {
		err = lastError("OpenClipboard")
		return
	}
	defer win.CloseClipboard()

	hMem := _HGLOBAL(win.GetClipboardData(formatID))
	if hMem == 0 {
		err = lastError("GetClipboardData")
		return
	}
	p := win.GlobalLock(hMem)
	if p == nil {
		err = lastError("GlobalLock")
		return
	}
	defer win.GlobalUnlock(hMem)
	l := win.GlobalSize(hMem)
	if l == 0 {
		err = lastError("GlobalSize")
		return
	}
	b = (*[1 << 30]byte)(unsafe.Pointer((*byte)(p)))[:l:l]
	return
}

func set(b []byte) error {
	if !win.OpenClipboard(0) {
		return lastError("OpenClipboard")
	}
	defer win.CloseClipboard()

	const GMEM_MOVEABLE = 0x0002
	hMem := win.GlobalAlloc(GMEM_MOVEABLE, uintptr(len(b)))
	if hMem == 0 {
		return lastError("GlobalAlloc")
	}
	p := win.GlobalLock(hMem)
	if p == nil {
		return lastError("GlobalLock")
	}
	win.MoveMemory(p, unsafe.Pointer(&b[0]), uintptr(len(b)))
	win.GlobalUnlock(hMem)
	if 0 == win.SetClipboardData(formatID, _HANDLE(hMem)) {
		defer win.GlobalFree(hMem)
		return lastError("SetClipboardData")
	}
	return nil
}
