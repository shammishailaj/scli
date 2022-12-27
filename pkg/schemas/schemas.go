package schemas

func NewImapServer(server string, port int, tls bool) *ImapServer {
	return &ImapServer{
		Server: server,
		Port:   port,
		TLS:    tls,
	}
}
