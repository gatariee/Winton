package winapi

type ICLRRuntimeInfoVtbl struct {
	QueryInterface         uintptr
	AddRef                 uintptr
	Release                uintptr
	GetVersionString       uintptr
	GetRuntimeDirectory    uintptr
	IsLoaded               uintptr
	LoadErrorString        uintptr
	LoadLibrary            uintptr
	GetProcAddress         uintptr
	GetInterface           uintptr
	IsLoadable             uintptr
	SetDefaultStartupFlags uintptr
	GetDefaultStartupFlags uintptr
	BindAsLegacyV2Runtime  uintptr
	IsStarted              uintptr
}

type ICLRRuntimeInfo struct {
	vtbl *ICLRRuntimeInfoVtbl
}

type ICORRuntimeHostVtbl struct {
	QueryInterface                uintptr
	AddRef                        uintptr
	Release                       uintptr
	CreateLogicalThreadState      uintptr
	DeleteLogicalThreadState      uintptr
	SwitchInLogicalThreadState    uintptr
	SwitchOutLogicalThreadState   uintptr
	LocksHeldByLogicalThreadState uintptr
	MapFile                       uintptr
	GetConfiguration              uintptr
	Start                         uintptr
	Stop                          uintptr
	CreateDomain                  uintptr
	GetDefaultDomain              uintptr
	EnumDomains                   uintptr
	NextDomain                    uintptr
	CloseEnum                     uintptr
	CreateDomainEx                uintptr
	CreateDomainSetup             uintptr
	CreateEvidence                uintptr
	UnloadDomain                  uintptr
	CurrentDomain                 uintptr
}

type ICORRuntimeHost struct {
	vtbl *ICORRuntimeHostVtbl
}


type ICLRRuntimeHost struct {
	vtbl *ICLRRuntimeHostVtbl
}

type ICLRRuntimeHostVtbl struct {
	QueryInterface            uintptr
	AddRef                    uintptr
	Release                   uintptr
	Start                     uintptr
	Stop                      uintptr
	SetHostControl            uintptr
	GetCLRControl             uintptr
	UnloadAppDomain           uintptr
	ExecuteInAppDomain        uintptr
	GetCurrentAppDomainId     uintptr
	ExecuteApplication        uintptr
	ExecuteInDefaultAppDomain uintptr
}


type IUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}

type IUnknown struct {
	vtbl *IUnknownVtbl
}


type AppDomain struct {
	vtbl *AppDomainVtbl
}

type AppDomainVtbl struct {
	QueryInterface            uintptr
	AddRef                    uintptr
	Release                   uintptr
	GetTypeInfoCount          uintptr
	GetTypeInfo               uintptr
	GetIDsOfNames             uintptr
	Invoke                    uintptr
	get_ToString              uintptr
	Equals                    uintptr
	GetHashCode               uintptr
	GetType                   uintptr
	InitializeLifetimeService uintptr
	GetLifetimeService        uintptr
	get_Evidence              uintptr
	add_DomainUnload          uintptr
	remove_DomainUnload       uintptr
	add_AssemblyLoad          uintptr
	remove_AssemblyLoad       uintptr
	add_ProcessExit           uintptr
	remove_ProcessExit        uintptr
	add_TypeResolve           uintptr
	remove_TypeResolve        uintptr
	add_ResourceResolve       uintptr
	remove_ResourceResolve    uintptr
	add_AssemblyResolve       uintptr
	remove_AssemblyResolve    uintptr
	add_UnhandledException    uintptr
	remove_UnhandledException uintptr
	DefineDynamicAssembly     uintptr
	DefineDynamicAssembly_2   uintptr
	DefineDynamicAssembly_3   uintptr
	DefineDynamicAssembly_4   uintptr
	DefineDynamicAssembly_5   uintptr
	DefineDynamicAssembly_6   uintptr
	DefineDynamicAssembly_7   uintptr
	DefineDynamicAssembly_8   uintptr
	DefineDynamicAssembly_9   uintptr
	CreateInstance            uintptr
	CreateInstanceFrom        uintptr
	CreateInstance_2          uintptr
	CreateInstanceFrom_2      uintptr
	CreateInstance_3          uintptr
	CreateInstanceFrom_3      uintptr
	Load                      uintptr
	Load_2                    uintptr
	Load_3                    uintptr
	Load_4                    uintptr
	Load_5                    uintptr
	Load_6                    uintptr
	Load_7                    uintptr
	ExecuteAssembly           uintptr
	ExecuteAssembly_2         uintptr
	ExecuteAssembly_3         uintptr
	get_FriendlyName          uintptr
	get_BaseDirectory         uintptr
	get_RelativeSearchPath    uintptr
	get_ShadowCopyFiles       uintptr
	GetAssemblies             uintptr
	AppendPrivatePath         uintptr
	ClearPrivatePath          uintptr
	SetShadowCopyPath         uintptr
	ClearShadowCopyPath       uintptr
	SetCachePath              uintptr
	SetData                   uintptr
	GetData                   uintptr
	SetAppDomainPolicy        uintptr
	SetThreadPrincipal        uintptr
	SetPrincipalPolicy        uintptr
	DoCallBack                uintptr
	get_DynamicDirectory      uintptr
}

type SafeArrayBound struct {
	cElements uint32
	lLbound   int32
}

type SafeArray struct {
	cDims      uint16
	fFeatures  uint16
	cbElements uint32
	cLocks     uint32
	pvData     uintptr
	rgsabound  [1]SafeArrayBound
}

type Assembly struct {
	vtbl *AssemblyVtbl
}

type AssemblyVtbl struct {
	QueryInterface              uintptr
	AddRef                      uintptr
	Release                     uintptr
	GetTypeInfoCount            uintptr
	GetTypeInfo                 uintptr
	GetIDsOfNames               uintptr
	Invoke                      uintptr
	get_ToString                uintptr
	Equals                      uintptr
	GetHashCode                 uintptr
	GetType                     uintptr
	get_CodeBase                uintptr
	get_EscapedCodeBase         uintptr
	GetName                     uintptr
	GetName_2                   uintptr
	get_FullName                uintptr
	get_EntryPoint              uintptr
	GetType_2                   uintptr
	GetType_3                   uintptr
	GetExportedTypes            uintptr
	GetTypes                    uintptr
	GetManifestResourceStream   uintptr
	GetManifestResourceStream_2 uintptr
	GetFile                     uintptr
	GetFiles                    uintptr
	GetFiles_2                  uintptr
	GetManifestResourceNames    uintptr
	GetManifestResourceInfo     uintptr
	get_Location                uintptr
	get_Evidence                uintptr
	GetCustomAttributes         uintptr
	GetCustomAttributes_2       uintptr
	IsDefined                   uintptr
	GetObjectData               uintptr
	add_ModuleResolve           uintptr
	remove_ModuleResolve        uintptr
	GetType_4                   uintptr
	GetSatelliteAssembly        uintptr
	GetSatelliteAssembly_2      uintptr
	LoadModule                  uintptr
	LoadModule_2                uintptr
	CreateInstance              uintptr
	CreateInstance_2            uintptr
	CreateInstance_3            uintptr
	GetLoadedModules            uintptr
	GetLoadedModules_2          uintptr
	GetModules                  uintptr
	GetModules_2                uintptr
	GetModule                   uintptr
	GetReferencedAssemblies     uintptr
	get_GlobalAssemblyCache     uintptr
}
type MethodInfo struct {
	vtbl *MethodInfoVtbl
}

type MethodInfoVtbl struct {
	QueryInterface                 uintptr
	AddRef                         uintptr
	Release                        uintptr
	GetTypeInfoCount               uintptr
	GetTypeInfo                    uintptr
	GetIDsOfNames                  uintptr
	Invoke                         uintptr
	get_ToString                   uintptr
	Equals                         uintptr
	GetHashCode                    uintptr
	GetType                        uintptr
	get_MemberType                 uintptr
	get_name                       uintptr
	get_DeclaringType              uintptr
	get_ReflectedType              uintptr
	GetCustomAttributes            uintptr
	GetCustomAttributes_2          uintptr
	IsDefined                      uintptr
	GetParameters                  uintptr
	GetMethodImplementationFlags   uintptr
	get_MethodHandle               uintptr
	get_Attributes                 uintptr
	get_CallingConvention          uintptr
	Invoke_2                       uintptr
	get_IsPublic                   uintptr
	get_IsPrivate                  uintptr
	get_IsFamily                   uintptr
	get_IsAssembly                 uintptr
	get_IsFamilyAndAssembly        uintptr
	get_IsFamilyOrAssembly         uintptr
	get_IsStatic                   uintptr
	get_IsFinal                    uintptr
	get_IsVirtual                  uintptr
	get_IsHideBySig                uintptr
	get_IsAbstract                 uintptr
	get_IsSpecialName              uintptr
	get_IsConstructor              uintptr
	Invoke_3                       uintptr
	get_returnType                 uintptr
	get_ReturnTypeCustomAttributes uintptr
	GetBaseDefinition              uintptr
}

type Variant struct {
	VT         uint16 // VARTYPE
	wReserved1 uint16
	wReserved2 uint16
	wReserved3 uint16
	Val        uintptr
	_          [8]byte
}

type ICLRMetaHost struct {
	vtbl *ICLRMetaHostVtbl
}

type ICLRMetaHostVtbl struct {
	QueryInterface                   uintptr
	AddRef                           uintptr
	Release                          uintptr
	GetRuntime                       uintptr
	GetVersionFromFile               uintptr
	EnumerateInstalledRuntimes       uintptr
	EnumerateLoadedRuntimes          uintptr
	RequestRuntimeLoadedNotification uintptr
	QueryLegacyV2RuntimeBinding      uintptr
	ExitProcess                      uintptr
}


type IEnumUnknown struct {
	vtbl *IEnumUnknownVtbl
}

type IEnumUnknownVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	Next           uintptr
	Skip           uintptr
	Reset          uintptr
	Clone          uintptr
}
