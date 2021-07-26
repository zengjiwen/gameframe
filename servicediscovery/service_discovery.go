package servicediscovery

import (
	"encoding/json"
	"fmt"
	"strings"
)

var _sd ServiceDiscovery

type ServiceDiscovery interface {
	Start() error
	GetRandomServer(serverType string) (*ServerInfo, bool)
	GetServer(serverID string) (*ServerInfo, bool)
	AddServerWatcher(watcher ServerWatcher)
	Close() error
}

func Get() ServiceDiscovery {
	return _sd
}

func Set(sd ServiceDiscovery) {
	_sd = sd
}

type ServerWatcher interface {
	OnAddServer(serverInfo *ServerInfo)
	OnRemoveServer(serverInfo *ServerInfo)
}

const _sdPrefix = "servers/"

var _serverInfo *ServerInfo

type ServerInfo struct {
	ID             string   `json:"id,omitempty"`
	Type           string   `json:"type,omitempty"`
	Addr           string   `json:"addr,omitempty"`
	ClientHandlers []string `json:"client_handlers,omitempty"`
	ServerHandlers []string `json:"server_handlers,omitempty"`
}

func InitServerInfo(serverID, serverType, serviceAddr string) {
	_serverInfo = &ServerInfo{
		ID:   serverID,
		Type: serverType,
		Addr: serviceAddr,
	}
}

func GetServerInfo() *ServerInfo {
	return _serverInfo
}

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
