package tls_client

import (
	"context"
	"crypto/tls"

	utls "github.com/refraction-networking/utls"
)

type UTLSConnectionAdapter struct {
	*utls.UConn
}

func (adapter *UTLSConnectionAdapter) ConnectionState() tls.ConnectionState {
	uState := adapter.UConn.ConnectionState()

	return tls.ConnectionState{
		Version:                     uState.Version,
		HandshakeComplete:           uState.HandshakeComplete,
		DidResume:                   uState.DidResume,
		CipherSuite:                 uState.CipherSuite,
		NegotiatedProtocol:          uState.NegotiatedProtocol,
		ServerName:                  uState.ServerName,
		PeerCertificates:            uState.PeerCertificates,
		VerifiedChains:              uState.VerifiedChains,
		SignedCertificateTimestamps: uState.SignedCertificateTimestamps,
		OCSPResponse:                uState.OCSPResponse,
		TLSUnique:                   uState.TLSUnique,
	}
}

func (adapter *UTLSConnectionAdapter) HandshakeContext(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- adapter.UConn.Handshake()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
