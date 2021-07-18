package servicediscovery

import (
	"encoding/json"
	"fmt"
	"github.com/zengjiwen/gameframe/services"
	"strings"
)

type ServiceDiscovery interface {
	GetRandomServer(serverType string) (*ServerInfo, bool)
	GetServer(serverID string) (*ServerInfo, bool)
	AddServerListener(sl ServerListener)
}

type ServerListener interface {
	OnAddServer(serverInfo *ServerInfo)
	OnRemoveServer(serverInfo *ServerInfo)
}

const _sdPrefix = "servers/"

type ServerInfo struct {
	ID             string   `json:"id,omitempty"`
	Type           string   `json:"type,omitempty"`
	Addr           string   `json:"addr,omitempty"`
	ClientHandlers []string `json:"client_handlers,omitempty"`
	ServerHandlers []string `json:"server_handlers,omitempty"`
}

func newServerInfo(serverID, serverType, serviceAddr string) *ServerInfo {
	serverInfo := &ServerInfo{
		ID:   serverID,
		Type: serverType,
		Addr: serviceAddr,
	}

	for ch := range services.ClientHandlers {
		serverInfo.ClientHandlers = append(serverInfo.ClientHandlers, ch)
	}
	for sh := range services.ServerHandlers {
		serverInfo.ServerHandlers = append(serverInfo.ServerHandlers, sh)
	}
	return serverInfo
}

var _serverInfo *ServerInfo

func parseSDKey(sdKey string) (string, string, error) {
	serverMetaData := strings.Split(sdKey, "/")
	if len(serverMetaData) != 3 {
		return "", "", fmt.Errorf("parse sd key error! key:%s", sdKey)
	}

	serverID := serverMetaData[1]
	serverType := serverMetaData[2]
	return serverID, serverType, nil
}

func genSDKey(serverID, serverType string) string {
	return fmt.Sprintf("%s%s/%s", _sdPrefix, serverID, serverType)
}

func parseSDValue(sdValue []byte) (*ServerInfo, error) {
	serverInfo := &ServerInfo{}
	err := json.Unmarshal(sdValue, serverInfo)
	if err != nil {
		return nil, err
	}

	return serverInfo, nil
}

func genSDValue(serverInfo *ServerInfo) (string, error) {
	value, err := json.Marshal(serverInfo)
	if err != nil {
		return "", err
	}

	return string(value), nil
}
