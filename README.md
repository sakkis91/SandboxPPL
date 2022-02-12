# SandboxPPL
Golang PoC that sandboxes Defender (or other PPL) by setting its token integrity to Untrusted, effectively turning it useless.

This is possible due to the fact that PROCESS_QUERY_LIMITED_INFORMATION is enough (in contrast to what MSDN documentation states) to get a handle to an access token of a protected process.

# PoC
MsMpEng.exe (Windows Defender) runs as PPL

<img width="574" alt="PPL-edited" src="https://user-images.githubusercontent.com/23586140/153692170-82a06188-903e-4d9b-a523-7f9a640ba3b9.png">

All sorts of privileges are enabled and token integrity is System

<img width="535" alt="system-integrity" src="https://user-images.githubusercontent.com/23586140/153692254-1a3c4adb-dedc-414f-a98a-e9fca4889ea8.png">

Not anymore

<img width="487" alt="tool-run" src="https://user-images.githubusercontent.com/23586140/153692359-8a1565cb-a034-4840-8261-45622d7b9849.PNG">

<img width="532" alt="untrusted" src="https://user-images.githubusercontent.com/23586140/153692364-5be541f5-71d5-4341-81fc-49cd66feed4f.png">

# Notes
In the original research all the privileges are manually stripped off from the process, besides changing the token integrity. 
It seems that this first step is not necessary, since only by setting the integrity level to Untrusted the same goal is achieved.

The program needs to run with SYSTEM privileges, otherwise the OpenProcessToken call will fail because the Owner of the target token is NT AUTHORITY\SYSTEM.

# Reference 

https://elastic.github.io/security-research/whitepapers/2022/02/02.sandboxing-antimalware-products-for-fun-and-profit/article/
