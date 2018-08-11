package hashgraph

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

type SocketHashgraphProxyClient struct {
	nodeAddr string
	timeout  time.Duration
}

func NewSocketHashgraphProxyClient(nodeAddr string, timeout time.Duration) *SocketHashgraphProxyClient {
	return &SocketHashgraphProxyClient{
		nodeAddr: nodeAddr,
		timeout:  timeout,
	}
}

func (p *SocketHashgraphProxyClient) getConnection() (*rpc.Client, error) {
	conn, err := net.DialTimeout("tcp", p.nodeAddr, p.timeout)
	if err != nil {
		return nil, err
	}
	return jsonrpc.NewClient(conn), nil
}

func (p *SocketHashgraphProxyClient) SubmitTx(tx []byte) (*bool, error) {
	rpcConn, err := p.getConnection()
	if err != nil {
		return nil, err
	}
	var ack bool
	err = rpcConn.Call("hashgraph.SubmitTx", tx, &ack)
	if err != nil {
		return nil, err
	}
	return &ack, nil
}
