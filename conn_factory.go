package tls_client

import (
	"crypto/tls"
	oohttp "github.com/ooni/oohttp"
	utls "github.com/refraction-networking/utls"
	"net"
)

type FactoryWithParrot struct {
	Parrot *utls.ClientHelloID
}

func (factory *FactoryWithParrot) NewTLSConn(conn net.Conn, config *tls.Config) oohttp.TLSConn {
	if factory.Parrot == nil {
		factory.Parrot = &utls.HelloChrome_Auto
	}

	return &ConnAdapter{utls.UClient(conn, &utls.Config{
		RootCAs:                     config.RootCAs,
		NextProtos:                  config.NextProtos,
		ServerName:                  config.ServerName,
		DynamicRecordSizingDisabled: config.DynamicRecordSizingDisabled,
		InsecureSkipVerify:          config.InsecureSkipVerify,
		ClientSessionCache:          utls.NewLRUClientSessionCache(0),
	}, *factory.Parrot)}
}
