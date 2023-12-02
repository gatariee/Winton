package winapi

import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

func (obj *MethodInfo) Invoke_3(variantObj Variant, parameters uintptr, pRetVal *uintptr) uintptr {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.Invoke_3,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&variantObj)),
		parameters,
		uintptr(unsafe.Pointer(pRetVal)),
		0,
		0,
	)
	return ret
}

func (obj *MethodInfo) GetString(addr *uintptr) error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.get_ToString,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(addr)),
		0,
	)
	return checkOK(ret, "get_ToString")
}

func (obj *IUnknown) QueryInterface(riid *windows.GUID, ppvObject *uintptr) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.QueryInterface,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppvObject)))
	return ret
}
