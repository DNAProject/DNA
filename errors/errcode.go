package errors

import (
	"fmt"
)

type ErrCoder interface {
	GetErrCode() ErrCode
}

type ErrCode int32

const (
	ErrNoCode               ErrCode = -2
	ErrNoError              ErrCode = 0
	ErrUnknown              ErrCode = -1
	ErrDuplicatedTx         ErrCode = 1
	ErrDuplicateInput       ErrCode = 45003
	ErrAssetPrecision       ErrCode = 45004
	ErrTransactionBalance   ErrCode = 45005
	ErrAttributeProgram     ErrCode = 45006
	ErrTransactionContracts ErrCode = 45007
	ErrTransactionPayload   ErrCode = 45008
	ErrDoubleSpend          ErrCode = 45009
	ErrTxHashDuplicate      ErrCode = 45010
	ErrStateUpdaterVaild    ErrCode = 45011
	ErrSummaryAsset         ErrCode = 45012
	ErrXmitFail             ErrCode = 45013
	ErrTooEarly             ErrCode = 45014
	ErrExpired              ErrCode = 45015
	ErrInternal             ErrCode = 45016
)

func (err ErrCode) Error() string {
	switch err {
	case ErrNoCode:
		return "No error code"
	case ErrNoError:
		return "Not an error"
	case ErrUnknown:
		return "Unknown error"
	case ErrDuplicatedTx:
		return "There are duplicated Transactions"
	case ErrTooEarly:
		return "Too early to be packed in the block"
	case ErrExpired:
		return "Expired"
	case ErrInternal:
		return "Internal error"
	}

	return fmt.Sprintf("Unknown error? Error code = %d", err)
}

func ErrerCode(err error) ErrCode {
	if err, ok := err.(ErrCoder); ok {
		return err.GetErrCode()
	}
	return ErrUnknown
}
