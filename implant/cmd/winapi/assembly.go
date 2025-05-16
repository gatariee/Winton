package winapi

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
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

func ExecuteByteArray(targetRuntime string, rawBytes []byte, params []string) error {
	if targetRuntime == "" {
		targetRuntime = "v4"
	}
	metahost, err := GetICLRMetaHost()
	if err != nil {
		return err
	}
	defer metahost.Release()

	runtimes, err := GetInstalledRuntimes(metahost)
	if err != nil {
		return err
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
		return err
	}
	defer runtimeInfo.Release()

	var isLoadable bool
	hr := runtimeInfo.IsLoadable(&isLoadable)
	if err = checkOK(hr, "runtimeInfo.IsLoadable"); err != nil {
		return err
	}
	if !isLoadable {
		return fmt.Errorf("%s is not loadable for some reason", latestRuntime)
	}

	runtimeHost, err := GetICORRuntimeHost(runtimeInfo)
	if err != nil {
		return err
	}
	defer runtimeHost.Release()

	appDomain, err := GetAppDomain(runtimeHost)
	if err != nil {
		return err
	}
	defer appDomain.Release()

	safeArrayPtr, err := CreateSafeArray(rawBytes)
	if err != nil {
		return err
	}

	var pAssembly uintptr
	hr = appDomain.Load_3(uintptr(safeArrayPtr), &pAssembly)
	if err = checkOK(hr, "appDomain.Load_3"); err != nil {
		return err
	}

	assembly := NewAssemblyFromPtr(pAssembly)

	var pEntryPointInfo uintptr
	hr = assembly.GetEntryPoint(&pEntryPointInfo)
	if err = checkOK(hr, "assembly.GetEntryPoint"); err != nil {
		return err
	}

	methodInfo := NewMethodInfoFromPtr(pEntryPointInfo)

	var methodSignaturePtr, paramPtr uintptr
	if err = methodInfo.GetString(&methodSignaturePtr); err != nil {
		return err
	}

	methodSignature := readUnicodeStr(unsafe.Pointer(methodSignaturePtr))

	if expectsParams(methodSignature) {
		if paramPtr, err = PrepareParameters(params); err != nil {
			return err
		}
	}

	nullVariant := Variant{
		VT:  1,
		Val: uintptr(0),
	}

	fmt.Printf("[*] Executing %s with params: %v\n", methodSignature, params)

	hr = methodInfo.Invoke_3(
		nullVariant,
		paramPtr,
		nil)
	
	fmt.Printf("[*] hr: 0x%08x\n", hr)
	return checkOK(hr, "methodInfo.Invoke_3")
}
