package hashgraph

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type SocketHashgraphProxy struct {
	nodeAddress string
	bindAddress string

	client *SocketHashgraphProxyClient
	server *SocketHashgraphProxyServer
}

func NewSocketHashgraphProxy(nodeAddr string,
	bindAddr string,
	timeout time.Duration,
	logger *logrus.Logger) (*SocketHashgraphProxy, error) {

	if logger == nil {
		logger = logrus.New()
		logger.Level = logrus.DebugLevel
	}

	client := NewSocketHashgraphProxyClient(nodeAddr, timeout)
	server, err := NewSocketHashgraphProxyServer(bindAddr, timeout, logger)
	if err != nil {
		return nil, err
	}

	proxy := &SocketHashgraphProxy{
		nodeAddress: nodeAddr,
		bindAddress: bindAddr,
		client:      client,
		server:      server,
	}
	go proxy.server.listen()

	return proxy, nil
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//Implement HashgraphProxy interface

func (p *SocketHashgraphProxy) CommitCh() chan Commit {
	return p.server.commitCh
}

func (p *SocketHashgraphProxy) SubmitTx(tx []byte) error {
	ack, err := p.client.SubmitTx(tx)
	if err != nil {
		return err
	}
	if !*ack {
		return fmt.Errorf("Failed to deliver transaction to hashgraph")
	}
	return nil
}
