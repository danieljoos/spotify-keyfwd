package main

import (
	"syscall"
	"unsafe"

	"github.com/TheTitanrain/w32"
)

// KeyboardCapture holds the Win32 handle and configured keys to forward.
type KeyboardCapture struct {
	keyboardHook w32.HHOOK
	captureKeys  []int

	KeyPressed chan int
}

// NewKeyboardCapture creates a new KeyboardCapture object.
// The object will only 'capture' the keys, specified in the given integer array.
// Specify keys by using the VK_* constants of Windows:
// http://msdn.microsoft.com/en-us/library/windows/desktop/dd375731(v=vs.85).aspx
func NewKeyboardCapture(captureKeys []int) *KeyboardCapture {
	return &KeyboardCapture{
		captureKeys: captureKeys,
		KeyPressed:  make(chan int),
	}
}

// SyncReceive creates a low-level keyboard hook using the SetWindowsHookEx function:
// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644990(v=vs.85).aspx
//
// Each intercepted key, which was included in the 'captureKeys' configuration
// variable (see NewKeyboardCapture), will be pushed to the 'KeyPressed' channel field.
// Returns an error in case the initialization of the hook failed.
// Calls to this function will block until KeyboardCapture.Stop() was called or the
// WM_QUIT message was sent to the current process.
func (t *KeyboardCapture) SyncReceive() error {
	isValidKey := func(key w32.DWORD) bool {
		for _, e := range t.captureKeys {
			if e == int(key) {
				return true
			}
		}
		return false
	}
	t.keyboardHook = w32.SetWindowsHookEx(w32.WH_KEYBOARD_LL,
		(w32.HOOKPROC)(func(code int, wparam w32.WPARAM, lparam w32.LPARAM) w32.LRESULT {
			if wparam == w32.WM_KEYDOWN {
				vkCode := (*w32.KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam)).VkCode
				if isValidKey(vkCode) {
					select {
					case t.KeyPressed <- int(vkCode):
					default:
					}
				}
			}
			return w32.CallNextHookEx(t.keyboardHook, code, wparam, lparam)
		}), 0, 0)
	if t.keyboardHook == 0 {
		return syscall.GetLastError()
	}
	var msg w32.MSG
	for w32.GetMessage(&msg, 0, 0, 0) != 0 {
	}
	w32.UnhookWindowsHookEx(t.keyboardHook)
	t.keyboardHook = 0

	return nil
}

// Stop stops the key interception by sending the quit message (WM_QUIT) to the current
// process.
func (t *KeyboardCapture) Stop() {
	w32.PostQuitMessage(0)
}
