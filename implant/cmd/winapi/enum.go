package winapi

import (
	"syscall"
	"fmt"
	"unsafe"
)


func (obj *IEnumUnknown) Release() uintptr {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0)
	return ret
}

func (obj *IEnumUnknown) Next(celt uint32, pEnumRuntime *uintptr, pCeltFetched *uint32) uintptr {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.Next,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(celt),
		uintptr(unsafe.Pointer(pEnumRuntime)),
		uintptr(unsafe.Pointer(pCeltFetched)),
		0,
		0)
	return ret
}

func GetInstalledRuntimes(metahost *ICLRMetaHost) ([]string, error) {
	var runtimes []string
	var pInstalledRuntimes uintptr
	hr := metahost.EnumerateInstalledRuntimes(&pInstalledRuntimes)
	err := checkOK(hr, "EnumerateInstalledRuntimes")
	if err != nil {
		return runtimes, err
	}
	installedRuntimes := NewIEnumUnknownFromPtr(pInstalledRuntimes)
	var pRuntimeInfo uintptr
	fetched := uint32(0)
	var versionString string
	versionStringBytes := make([]uint16, 20)
	versionStringSize := uint32(len(versionStringBytes))
	var runtimeInfo *ICLRRuntimeInfo
	for {
		hr = installedRuntimes.Next(1, &pRuntimeInfo, &fetched)
		if hr != 0x00000000 {
			break
		}
		runtimeInfo = NewICLRRuntimeInfoFromPtr(pRuntimeInfo)
		if ret := runtimeInfo.GetVersionString(&versionStringBytes[0], &versionStringSize); ret != 0x00000000 {
			return runtimes, fmt.Errorf("GetVersionString returned 0x%08x", ret)
		}
		versionString = syscall.UTF16ToString(versionStringBytes)
		runtimes = append(runtimes, versionString)
	}
	if len(runtimes) == 0 {
		return runtimes, fmt.Errorf("no runtimes found")
	}
	runtimeInfo.Release()
	return runtimes, err
}

func GetICORRuntimeHost(runtimeInfo *ICLRRuntimeInfo) (*ICORRuntimeHost, error) {
	var pRuntimeHost uintptr
	hr := runtimeInfo.GetInterface(&CLSID_CorRuntimeHost, &IID_ICorRuntimeHost, &pRuntimeHost)
	err := checkOK(hr, "runtimeInfo.GetInterface")
	if err != nil {
		return nil, err
	}
	runtimeHost := NewICORRuntimeHostFromPtr(pRuntimeHost)
	hr = runtimeHost.Start()
	err = checkOK(hr, "runtimeHost.Start")
	return runtimeHost, err
}

func GetRuntimeInfo(metahost *ICLRMetaHost, version string) (*ICLRRuntimeInfo, error) {
	pwzVersion, err := syscall.UTF16PtrFromString(version)
	if err != nil {
		return nil, err
	}
	var pRuntimeInfo uintptr
	hr := metahost.GetRuntime(pwzVersion, &IID_ICLRRuntimeInfo, &pRuntimeInfo)
	err = checkOK(hr, "metahost.GetRuntime")
	if err != nil {
		return nil, err
	}
	return NewICLRRuntimeInfoFromPtr(pRuntimeInfo), nil
}