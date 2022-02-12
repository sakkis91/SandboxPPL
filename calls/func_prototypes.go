package calls
//Run the below command to produce calls.go
//mkwinsyscall.exe -output calls.go func_prototypes.go

//sys SetTokenInformation(tokenHandle syscall.Token, tokenInformationClass uint32, tokenInformation uintptr, tokenInformationLength uint32) (err error) = advapi32.SetTokenInformation
//sys ConvertStringSidToSid(stringSid *uint16, sid **structs.SID) (err error) = advapi32.ConvertStringSidToSidW
//sys GetLengthSid(sid *structs.SID) (len uint32) = advapi32.GetLengthSid
