package winapi

import (
	"syscall"
	"unsafe"
)

func (obj *AppDomain) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *AppDomain) Load_3(pRawAssembly uintptr, asmbly *uintptr) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Load_3,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pRawAssembly)),
		uintptr(unsafe.Pointer(asmbly)))
	return ret
}

func NewAppDomainFromPtr(ppv uintptr) *AppDomain {
	return (*AppDomain)(unsafe.Pointer(ppv))
}

func GetAppDomain(runtimeHost *ICORRuntimeHost) (appDomain *AppDomain, err error) {
	var pAppDomain uintptr
	var pIUnknown uintptr
	hr := runtimeHost.GetDefaultDomain(&pIUnknown)
	err = checkOK(hr, "runtimeHost.GetDefaultDomain")
	if err != nil {
		return
	}
	iu := NewIUnknownFromPtr(pIUnknown)
	hr = iu.QueryInterface(&IID_AppDomain, &pAppDomain)
	err = checkOK(hr, "IUnknown.QueryInterface")
	return NewAppDomainFromPtr(pAppDomain), err
}

