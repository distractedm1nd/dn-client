package client

import (
	"context"
	"fmt"
	"net/http"

	"cosmossdk.io/math"
	"github.com/filecoin-project/go-jsonrpc"

	"github.com/celestiaorg/nmt/namespace"
)

type State struct {
	SubmitPayForBlob func(
		ctx context.Context,
		nID namespace.ID,
		data []byte,
		fee math.Int,
		gasLim uint64,
	) (*TxResponse, error)
}

type Share struct {
	GetSharesByNamespace func(
		ctx context.Context,
		root *DataAvailabilityHeader,
		namespace namespace.ID,
	) (NamespacedShares, error)
}

type Header struct {
	GetByHeight func(context.Context, uint64) (*ExtendedHeader, error)
	Head        func(context.Context) (*ExtendedHeader, error)
}

type Client struct {
	State
	Header
	Share

	closer multiClientCloser
}

// multiClientCloser is a wrapper struct to close clients across multiple namespaces.
type multiClientCloser struct {
	closers []jsonrpc.ClientCloser
}

// register adds a new closer to the multiClientCloser
func (m *multiClientCloser) register(closer jsonrpc.ClientCloser) {
	m.closers = append(m.closers, closer)
}

// closeAll closes all saved clients.
func (m *multiClientCloser) closeAll() {
	for _, closer := range m.closers {
		closer()
	}
}

// Close closes the connections to all namespaces registered on the client.
func (c *Client) Close() {
	c.closer.closeAll()
}

// NewPublicClient creates a new Client with one connection per namespace.
func NewPublicClient(ctx context.Context, addr string) (*Client, error) {
	return newClient(ctx, addr, nil)
}

// NewClient creates a new Client with one connection per namespace with the
// given token as the authorization token.
func NewClient(ctx context.Context, addr string, token string) (*Client, error) {
	authHeader := http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", token)}}
	return newClient(ctx, addr, authHeader)
}

func newClient(ctx context.Context, addr string, authHeader http.Header) (*Client, error) {
	var client Client
	var multiCloser multiClientCloser

	// TODO: this duplication of strings many times across the codebase can be avoided with issue #1176
	var modules = map[string]interface{}{
		"share":  &client.Share,
		"state":  &client.State,
		"header": &client.Header,
	}
	for name, module := range modules {
		closer, err := jsonrpc.NewClient(ctx, addr, name, module, authHeader)
		if err != nil {
			return nil, err
		}
		multiCloser.register(closer)
	}

	return &client, nil
}
