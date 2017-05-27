package common

import (
	. "DNA/common"
)

type ApiServer interface {
	Start() error
	Stop()
	Push(txHash Uint256, errcode interface{}, result interface{})
}

//Node
func GetNodeInfo(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "GetNodeInfo",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}
func GetNodeCount(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "GetNodeCount",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}

//Block
func PostBlock(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "PostBlock",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}
func GetBlockHeight(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "GetBlockHeight",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}
func GetBlockByHash(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "GetBlockByHash",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}
func GetBlockByHeight(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "GetBlockByHeight",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}

//Smartcode
func PostSmartcodeInvoke(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "PostSmartcodeInvoke",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}
func PostSmartcodePublish(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "PostSmartcodePublish",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}
func GetSmartcodeInfo(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"iesult": "GetSmartcodeInfo",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}

//Asset
func GetAssetByHash(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "GetAssetByHash",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}
func PostAssetTransfer(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "PostAssetTransfer",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}
func PostAssetIssue(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "PostAssetIssue",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}
func PostAssetRegistry(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "PostAssetRegistry",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}

//Transaction
func GetTransactionByHash(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "GetTransactionByHash",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}
func GetTransactionsInMempool(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "GetTransactionMempool",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}

//Record
func GetRecordByHash(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "GetRecordByHash",
		"error":  0,
		"id":     0,
	}
	//TODO   process req and return result
	return resp
}
func PostRecord(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"result": "PostRecord",
		"error":  0,
		"id":     cmd["id"],
	}
	//TODO   process req and return result
	return resp
}

//other
func WsHeartbeat(cmd map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"heartbeat-ack": "",
		"error":         0,
		"result":        "",
		"id":            "",
	}
	return resp
}
func Test(cmd map[string]interface{}) map[string]interface{} {

	resp := map[string]interface{}{
		"error":  0,
		"result": "test",
		"id":     cmd["id"],
	}
	return resp
}
