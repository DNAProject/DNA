package httpjsonrpc

import (
	"DNA/core/transaction"
	"DNA/core/code"
)

func NewSamplePublish(fc *code.FunctionCode,icFunctionCode []byte,name string,codeversion string,
author string,email string,desp string) *transaction.Transaction {

	// generate transaction
	tx, _ := transaction.NewPublishTransaction(fc,icFunctionCode,name,codeversion,author,email,desp)
	return tx
}

func NewSampleInvoke(fc []byte) *transaction.Transaction {

	// generate transaction
	tx, _ := transaction.NewInvokeTransaction(fc)
	return tx
}

//func NewSampleInvoke(code []byte) *transaction.Transaction {
//	r :=bytes.NewBuffer(code)
//	FunctionCode := new(code.FunctionCode)
//	FunctionCode.Deserialize(r)
//
//	// generate transaction
//	tx, _ := transaction.NewInvokeTransaction(FunctionCode)
//	return tx
//}