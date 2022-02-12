package structs

type TOKEN_MANDATORY_LABEL struct {
	Label SID_AND_ATTRIBUTES
}

type SID struct {
	Revision byte
	SubAuthorityCount byte
	IdentifierAuthority SID_IDENTIFIER_AUTHORITY
	SubAuthority [1]int
}

type SID_AND_ATTRIBUTES struct {
	Sid        *SID
	Attributes uint32
}

type SID_IDENTIFIER_AUTHORITY struct {
	Value [6]byte
}
