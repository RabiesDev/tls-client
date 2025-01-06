package tls_client

import (
	"crypto/tls"
	"net"

	oohttp "github.com/ooni/oohttp"
	utls "github.com/refraction-networking/utls"
)

type SecureTLSFactory struct {
	ClientProfile *utls.ClientHelloID
}

func (factory *SecureTLSFactory) NewTLSConnection(conn net.Conn, config *tls.Config) oohttp.TLSConn {
	if factory.ClientProfile == nil {
		factory.ClientProfile = &utls.HelloChrome_Auto
	}

	uConfig := &utls.Config{
		RootCAs:                     config.RootCAs,
		NextProtos:                  config.NextProtos,
		ServerName:                  config.ServerName,
		DynamicRecordSizingDisabled: config.DynamicRecordSizingDisabled,
		InsecureSkipVerify:          config.InsecureSkipVerify,
		ClientSessionCache:          utls.NewLRUClientSessionCache(0),
	}

	uConn := utls.UClient(conn, uConfig, *factory.ClientProfile)
	return &UTLSConnectionAdapter{UConn: uConn}
}
