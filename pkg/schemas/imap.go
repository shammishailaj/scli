package schemas

import (
	"encoding/json"
	"fmt"
)

type ImapServer struct {
	Server string `json:"Server"`
	Port   int    `json:"Port"`
	TLS    bool   `json:"TLS"`
}

func (i *ImapServer) String() string {
	return fmt.Sprintf("Server: %s\nPort: %d\nTLS: %t\n", i.Server, i.Port, i.TLS)
}

func (i *ImapServer) ServerPort() string {
	return fmt.Sprintf("%s:%d", i.Server, i.Port)
}

func (i *ImapServer) JSON() ([]byte, error) {
	return json.Marshal(i)
}
