package interfaces

type IScriptContainer interface {
	GetMessage() ([]byte)
	IInteropInterface
}