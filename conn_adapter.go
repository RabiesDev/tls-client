package tls_client

import (
	"context"
	"crypto/tls"
	utls "github.com/refraction-networking/utls"
)

type ConnAdapter struct {
	*utls.UConn
}

func (adapter *ConnAdapter) ConnectionState() tls.ConnectionState {
	state := adapter.UConn.ConnectionState()
	return tls.ConnectionState{
		Version:                     state.Version,
		HandshakeComplete:           state.HandshakeComplete,
		DidResume:                   state.DidResume,
		CipherSuite:                 state.CipherSuite,
		NegotiatedProtocol:          state.NegotiatedProtocol,
		NegotiatedProtocolIsMutual:  state.NegotiatedProtocolIsMutual,
		ServerName:                  state.ServerName,
		PeerCertificates:            state.PeerCertificates,
		VerifiedChains:              state.VerifiedChains,
		SignedCertificateTimestamps: state.SignedCertificateTimestamps,
		OCSPResponse:                state.OCSPResponse,
		TLSUnique:                   state.TLSUnique,
	}
}

func (adapter *ConnAdapter) HandshakeContext(ctx context.Context) error {
	errch := make(chan error, 1)
	go func() {
		errch <- adapter.UConn.Handshake()
	}()

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
