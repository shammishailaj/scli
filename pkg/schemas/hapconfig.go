package schemas

import (
	"fmt"
	"sort"
)

/*
	backend backend.ssh.hocrm002
		mode	tcp
		server	hocrm002	10.0.1.65:22	check
		timeout	server	2h

	frontend frontend.ssh.hocrm002
		bind	*:16002
		mode	tcp
		timeout	client	2h
		default_backend	backend.ssh.hocrm002
*/

type HapconfigBackend struct {
	ForService    string
	Mode          string
	ServerName    string
	ServerIP      string
	ServerPort    int64
	ServerTimeOut string
}

type HapconfigFrontend struct {
	BackendName   string
	BindPort      int64
	ClientTimeOut string
	ForService    string
	Mode          string
	ServerName    string
}

func (h *HapconfigBackend) BackendName() string {
	return fmt.Sprintf("backend.%s.%s", h.ForService, h.ServerName)
}

func (h *HapconfigBackend) TCPBackendString() string {
	return fmt.Sprintf("backend %s\n\tmode\ttcp\n\tserver\t%s\t%s:%d\tcheck\n\ttimeout\tserver\t%s", h.BackendName(), h.ServerName, h.ServerIP, h.ServerPort, h.ServerTimeOut)
}

func (h *HapconfigFrontend) FrontendName() string {
	return fmt.Sprintf("frontend.%s.%s", h.ForService, h.ServerName)
}

func (h *HapconfigFrontend) TCPFrontendString() string {
	return fmt.Sprintf("frontend %s\n\tbind\t*:%d\n\tmode\ttcp\n\ttimeout\tclient\t%s\n\tdefault_backend\t%s", h.FrontendName(), h.BindPort, h.ClientTimeOut, h.BackendName)
}

type HapconfigSSHProxy struct {
	Mode          string
	ServerName    string
	ServerIP      string
	ServerPort    int64
	ServerTimeOut string
	BindPort      int64
	ClientTimeOut string
}

func (h *HapconfigSSHProxy) BackendName() string {
	return fmt.Sprintf("backend.ssh.%s", h.ServerName)
}

func (h *HapconfigSSHProxy) TCPBackendString() string {
	return fmt.Sprintf("backend %s\n\tmode\ttcp\n\tserver\t%s\t%s:%d\tcheck\n\ttimeout\tserver\t%s", h.BackendName(), h.ServerName, h.ServerIP, h.ServerPort, h.ServerTimeOut)
}

func (h *HapconfigSSHProxy) FrontendName() string {
	return fmt.Sprintf("frontend.ssh.%s", h.ServerName)
}

func (h *HapconfigSSHProxy) TCPFrontendString() string {
	return fmt.Sprintf("frontend %s\n\tbind\t*:%d\n\tmode\ttcp\n\ttimeout\tclient\t%s\n\tdefault_backend\t%s", h.FrontendName(), h.BindPort, h.ClientTimeOut, h.BackendName())
}

func (h *HapconfigSSHProxy) ProxyString() string {
	return fmt.Sprintf("%s\n\n%s", h.TCPBackendString(), h.TCPFrontendString())
}

// HapconfigSSHProxySorted Implementing a sorted struct as per: https://procrypt.github.io/post/2017-06-01-sorting-structs-in-golang/
type HapconfigSSHProxySorted []HapconfigSSHProxy

func (h HapconfigSSHProxySorted) Len() int {
	return len(h)
}

func (h HapconfigSSHProxySorted) Less(i, j int) bool {
	return h[i].BindPort < h[j].BindPort
}

func (h HapconfigSSHProxySorted) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h HapconfigSSHProxySorted) Sort() {
	sort.Sort(h)
}
