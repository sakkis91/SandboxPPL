package funcs

import (
   "fmt"
   "unsafe"
   "strings"
   "unicode/utf16"
   "log"
   "syscall"
   "golang.org/x/sys/windows"
)

const procEntrySize = (uint32)(unsafe.Sizeof(windows.ProcessEntry32{}))

func EnablePrivilegeOnCurrentProcess(priv string) error {
	var tHandle windows.Token
	cHandle, err := windows.GetCurrentProcess()
	if err != nil {
		return err
	}
	windows.OpenProcessToken(cHandle, windows.TOKEN_ADJUST_PRIVILEGES, &tHandle)
	var luid windows.LUID
	err = windows.LookupPrivilegeValue(nil, UTF16PtrFromString(priv), &luid)
	if err != nil {
		log.Println("LookupPrivilegeValue failed", err)
		return err
	}
	tp := windows.Tokenprivileges{}
	tp.PrivilegeCount = 1
	tp.Privileges[0].Attributes = windows.SE_PRIVILEGE_ENABLED
	tp.Privileges[0].Luid = luid
	err = windows.AdjustTokenPrivileges(tHandle, false, &tp, 0, nil, nil)
	if err != nil {
		log.Println("AdjustTokenPrivileges failed", err)
		return err
	}
	return nil
}

func UTF16FromString(s string) ([]uint16, error) {
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			return nil, syscall.EINVAL
		}
	}
	return utf16.Encode([]rune(s + "\x00")), nil
}

func UTF16PtrFromString(s string) (*uint16) {
	a, err := UTF16FromString(s)
	if err != nil {
		return nil
	}
	return &a[0]
}

func GetProcessID(name string) (uint32, error) {
   backslash := "\\"
   trimmedName := after(name,backslash)
   h, e := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
   if e != nil { return 0, e }
   p := windows.ProcessEntry32{Size: procEntrySize}
   for {
      e := windows.Process32Next(h, &p)
      if e != nil { return 0, e }
      if windows.UTF16ToString(p.ExeFile[:]) == trimmedName {
         return p.ProcessID, nil
      }
   }
   return 0, fmt.Errorf("%q not found", name)
}

func after(value string, backslash string) string {
    //Check if characters exist in string
    if (strings.Contains(value, backslash)) {
    // Get substring after a string.
    pos := strings.LastIndex(value, backslash)
    if pos == -1 {
        return ""
    }
    adjustedPos := pos + len(backslash)
    if adjustedPos >= len(value) {
        return ""
    }
    return value[adjustedPos:len(value)]
} else {
    return value
}
}
