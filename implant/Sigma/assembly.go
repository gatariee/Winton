package main 

import (
	"unsafe"
	"fmt"
	"strings"
	"syscall"
)

func (obj *Assembly) GetEntryPoint(pMethodInfo *uintptr) uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.get_EntryPoint,
		2,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(pMethodInfo)),
		0)
	return ret
}

func ExecuteByteArray(targetRuntime string, rawBytes []byte, params []string) (retCode int32, err error) {
	retCode = -1
	if targetRuntime == "" {
		targetRuntime = "v4"
	}
	metahost, err := GetICLRMetaHost()
	if err != nil {
		return
	}

	runtimes, err := GetInstalledRuntimes(metahost)
	if err != nil {
		return
	}
	var latestRuntime string
	for _, r := range runtimes {
		if strings.Contains(r, targetRuntime) {
			latestRuntime = r
			break
		} else {
			latestRuntime = r
		}
	}
	runtimeInfo, err := GetRuntimeInfo(metahost, latestRuntime)
	if err != nil {
		return
	}
	var isLoadable bool
	hr := runtimeInfo.IsLoadable(&isLoadable)
	err = checkOK(hr, "runtimeInfo.IsLoadable")
	if err != nil {
		return
	}
	if !isLoadable {
		return -1, fmt.Errorf("%s is not loadable for some reason", latestRuntime)
	}
	runtimeHost, err := GetICORRuntimeHost(runtimeInfo)
	if err != nil {
		return
	}
	appDomain, err := GetAppDomain(runtimeHost)
	if err != nil {
		return
	}
	safeArrayPtr, err := CreateSafeArray(rawBytes)
	if err != nil {
		return
	}
	var pAssembly uintptr
	hr = appDomain.Load_3(uintptr(safeArrayPtr), &pAssembly)
	err = checkOK(hr, "appDomain.Load_3")
	if err != nil {
		return
	}
	assembly := NewAssemblyFromPtr(pAssembly)
	var pEntryPointInfo uintptr
	hr = assembly.GetEntryPoint(&pEntryPointInfo)
	err = checkOK(hr, "assembly.GetEntryPoint")
	if err != nil {
		return
	}
	methodInfo := NewMethodInfoFromPtr(pEntryPointInfo)

	var methodSignaturePtr, paramPtr uintptr
	err = methodInfo.GetString(&methodSignaturePtr)
	if err != nil {
		return
	}
	methodSignature := readUnicodeStr(unsafe.Pointer(methodSignaturePtr))

	if expectsParams(methodSignature) {
		if paramPtr, err = PrepareParameters(params); err != nil {
			return
		}
	}

	var pRetCode uintptr
	nullVariant := Variant{
		VT:  1,
		Val: uintptr(0),
	}
	hr = methodInfo.Invoke_3(
		nullVariant,
		paramPtr,
		&pRetCode)
	err = checkOK(hr, "methodInfo.Invoke_3")
	if err != nil {
		return
	}
	appDomain.Release()
	runtimeHost.Release()
	runtimeInfo.Release()
	metahost.Release()
	return int32(pRetCode), nil
}