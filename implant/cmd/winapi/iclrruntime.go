package winapi

import (
	"syscall"
	"unsafe"	
	"golang.org/x/sys/windows"
)

func (obj *ICORRuntimeHost) Start() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Start,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICORRuntimeHost) GetDefaultDomain(pAppDomain *uintptr) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.GetDefaultDomain,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pAppDomain)),
		0)
	return ret
}

func (obj *ICORRuntimeHost) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func NewICLRRuntimeInfoFromPtr(ppv uintptr) *ICLRRuntimeInfo {
	return (*ICLRRuntimeInfo)(unsafe.Pointer(ppv))
}

func (obj *ICLRRuntimeInfo) IsLoadable(pbLoadable *bool) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.IsLoadable,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pbLoadable)),
		0)
	return ret
}

func (obj *ICLRRuntimeInfo) BindAsLegacyV2Runtime() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.BindAsLegacyV2Runtime,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return ret
}

func (obj *ICLRRuntimeInfo) GetInterface(rclsid *windows.GUID, riid *windows.GUID, ppUnk *uintptr) uintptr {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.GetInterface,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(rclsid)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppUnk)),
		0,
		0)
	return ret
}


func (obj *ICLRRuntimeHost) Start() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Start,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *ICLRRuntimeInfo) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}


func (obj *ICLRRuntimeInfo) GetVersionString(pcchBuffer *uint16, pVersionstringSize *uint32) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.GetVersionString,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pcchBuffer)),
		uintptr(unsafe.Pointer(&pVersionstringSize)))
	return ret
}







