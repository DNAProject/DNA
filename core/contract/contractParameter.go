package contract

type ContractParameterType byte

const (
	Signature ContractParameterType = iota
	Integer
	Hash160
	Hash256
	ByteArray
)

