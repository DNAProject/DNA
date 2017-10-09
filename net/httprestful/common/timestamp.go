package common

import (
	. "DNA/common"
	"DNA/common/config"
	Err "DNA/net/httprestful/error"
	ts "DNA/net/timestamp"
	"crypto/sha256"
	"fmt"
)

var timestampClients map[string]*ts.TimeStampClient

const DEFAULT_TIMESTAMPSOURCE = "http://timestamp.sheca.com/Timestamp/pdftime.do"

func init() {
	timestampClients = make(map[string]*ts.TimeStampClient, 0)
	for _, url := range config.Parameters.TimestampSources {
		timestampClients[url] = ts.NewClient(url)
	}
	if len(timestampClients) == 0 {
		timestampClients[DEFAULT_TIMESTAMPSOURCE] = ts.NewClient(DEFAULT_TIMESTAMPSOURCE)
	}
}

func GetTimestamp(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.TIMESTAMP_UNAVAILABLE)
	hex, ok := cmd["RecordHash"].(string)
	if !ok {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	hash, err := HexToBytes(hex)
	if err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}

	if len(hash) != sha256.Size {
		resp["Error"] = Err.INVALID_PARAMS
		resp["Desc"] = "wrong RecordHash type"
		return resp
	}

	success := false
	for _, client := range timestampClients {
		token, time, err := client.FetchTimeStampToken(hash)
		if err != nil {
			resp["Error"] = Err.TIMESTAMP_ERROR
			resp["Desc"] = err.Error()
			continue
		}

		resp["Result"] = map[string]interface{}{
			"token":     fmt.Sprintf("%x", token),
			"timestamp": time,
		}
		success = true
		break
	}
	if success {
		resp["Error"] = Err.SUCCESS
		resp["Desc"] = ""
	}

	return resp
}
