package contract

//parameter defined type.
type ContractParameterType byte

const (
	Signature ContractParameterType = iota
	Integer
	Hash160
	Hash256
	ByteArray
)

