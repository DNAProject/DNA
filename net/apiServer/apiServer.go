package apiServer

import (
	. "DNA/common"
	. "DNA/common/config"
	. "DNA/net/apiServer/common"
	. "DNA/net/apiServer/restful"
	. "DNA/net/apiServer/websocket"
	"strings"
)

var pushServer ApiServer

func StartServers() {
	servers := func() map[string]string {
		serverMap := make(map[string]string)
		strList := strings.Split(Parameters.HttpServers, ",")
		for _, v := range strList {
			temp := strings.TrimSpace(v)
			serverMap[temp] = temp
		}
		return serverMap
	}()
	for _, v := range servers {
		switch v {
		case "local":
			//TODO
			//go StartLocalServer()
			break
		case "rpc":
			//TODO
			//go StartRPCServer()
			break
		case "http":
			func() ApiServer {
				rest := InitRestServer()
				go rest.Start()
				return rest
			}()
			break
		case "ws":
			pushServer = func() ApiServer {
				ws := InitWsServer()
				go ws.Start()
				return ws
			}()
			break
		default:
			break
		}
	}
}

func Push(txHash Uint256, errcode interface{}, result interface{}) {
	if pushServer != nil {
		pushServer.Push(txHash, errcode, result)
	}

}
