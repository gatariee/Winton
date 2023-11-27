package main 

import (
	"syscall"
	"unsafe"
	"golang.org/x/sys/windows"
	"fmt"
	"strings"
	"bytes"
	"unicode/utf16"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func NewICLRMetaHostFromPtr(ppv uintptr) *ICLRMetaHost {
	return (*ICLRMetaHost)(unsafe.Pointer(ppv))
}

func NewICORRuntimeHostFromPtr(ppv uintptr) *ICORRuntimeHost {
	return (*ICORRuntimeHost)(unsafe.Pointer(ppv))
}

func NewIUnknownFromPtr(ppv uintptr) *IUnknown {
	return (*IUnknown)(unsafe.Pointer(ppv))
}

func CreateEmptySafeArray(arrayType int, size int) (unsafe.Pointer, error) {
	modOleAuto := syscall.MustLoadDLL("OleAut32.dll")
	procSafeArrayCreate := modOleAuto.MustFindProc("SafeArrayCreate")

	sab := SafeArrayBound{
		cElements: uint32(size),
		lLbound:   0,
	}
	vt := uint16(arrayType)
	ret, _, err := procSafeArrayCreate.Call(
		uintptr(vt),
		uintptr(1),
		uintptr(unsafe.Pointer(&sab)))

	if err != syscall.Errno(0) {
		return nil, err
	}

	return unsafe.Pointer(ret), nil
}

func CreateSafeArray(rawBytes []byte) (unsafe.Pointer, error) {
	saPtr, err := CreateEmptySafeArray(0x11, len(rawBytes)) // VT_UI1
	if err != nil {
		return nil, err
	}
	modNtDll := syscall.MustLoadDLL("ntdll.dll")
	procRtlCopyMemory := modNtDll.MustFindProc("RtlCopyMemory")
	sa := (*SafeArray)(saPtr)
	_, _, err = procRtlCopyMemory.Call(
		sa.pvData,
		uintptr(unsafe.Pointer(&rawBytes[0])),
		uintptr(len(rawBytes)))
	if err != syscall.Errno(0) {
		return nil, err
	}
	return saPtr, nil
}


func NewAssemblyFromPtr(ppv uintptr) *Assembly {
	return (*Assembly)(unsafe.Pointer(ppv))
}

func NewMethodInfoFromPtr(ppv uintptr) *MethodInfo {
	return (*MethodInfo)(unsafe.Pointer(ppv))
}


func CLRCreateInstance(clsid, riid *windows.GUID, ppInterface *uintptr) uintptr {
	ret, _, _ := procCLRCreateInstance.Call(
		uintptr(unsafe.Pointer(clsid)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppInterface)))
	return ret
}

func checkOK(hr uintptr, msg string) error {
	if hr != 0x00000000 {
		return fmt.Errorf("%s returned 0x%08x", msg, hr)
	}
	return nil
}

func GetICLRMetaHost() (metahost *ICLRMetaHost, err error) {
	var pMetaHost uintptr
	hr := CLRCreateInstance(&CLSID_CLRMetaHost, &IID_ICLRMetaHost, &pMetaHost)
	err = checkOK(hr, "CLRCreateInstance")
	if err != nil {
		return
	}
	metahost = NewICLRMetaHostFromPtr(pMetaHost)
	return
}

func NewIEnumUnknownFromPtr(ppv uintptr) *IEnumUnknown {
	return (*IEnumUnknown)(unsafe.Pointer(ppv))
}

func expectsParams(input string) bool {
	return !strings.Contains(input, "Void Main()")
}

func readUnicodeStr(ptr unsafe.Pointer) string {
	var byteVal uint16
	out := make([]uint16, 0)
	for i := 0; ; i++ {
		byteVal = *(*uint16)(unsafe.Pointer(ptr))
		if byteVal == 0x0000 {
			break
		}
		out = append(out, byteVal)
		ptr = unsafe.Pointer(uintptr(ptr) + 2)
	}
	return string(utf16.Decode(out))
}

func utf16Le(s string) []byte {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	var buf bytes.Buffer
	t := transform.NewWriter(&buf, enc)
	t.Write([]byte(s))
	return buf.Bytes()
}

func SysAllocString(str string) (unsafe.Pointer, error) {
	modOleAuto := syscall.MustLoadDLL("OleAut32.dll")
	sysAllocString := modOleAuto.MustFindProc("SysAllocString")
	input := utf16Le(str)
	ret, _, err := sysAllocString.Call(
		uintptr(unsafe.Pointer(&input[0])),
	)
	if err != syscall.Errno(0) {
		return nil, err
	}
	return unsafe.Pointer(ret), nil
}

func SafeArrayPutElement(array, btsr unsafe.Pointer, index int) (err error) {
	modOleAuto := syscall.MustLoadDLL("OleAut32.dll")
	safeArrayPutElement := modOleAuto.MustFindProc("SafeArrayPutElement")
	_, _, err = safeArrayPutElement.Call(
		uintptr(array),
		uintptr(unsafe.Pointer(&index)),
		uintptr(btsr),
	)
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func PrepareParameters(params []string) (uintptr, error) {
	listStrSafeArrayPtr, err := CreateEmptySafeArray(0x0008, len(params)) 
	if err != nil {
		return 0, err
	}
	for i, p := range params {
		bstr, _ := SysAllocString(p)
		SafeArrayPutElement(listStrSafeArrayPtr, bstr, i)
	}

	paramVariant := Variant{
		VT:  0x0008 | 0x2000, 
		Val: uintptr(listStrSafeArrayPtr),
	}

	paramsSafeArrayPtr, err := CreateEmptySafeArray(0x000C, 1) 
	if err != nil {
		return 0, err
	}
	err = SafeArrayPutElement(paramsSafeArrayPtr, unsafe.Pointer(&paramVariant), 0)
	if err != nil {
		return 0, err
	}
	return uintptr(paramsSafeArrayPtr), nil
}
