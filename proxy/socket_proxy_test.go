package proxy

import (
	"reflect"
	"testing"
	"time"

	"github.com/andrecronje/hashgraph/common"
	"github.com/andrecronje/hashgraph/hashgraph"
	aproxy "github.com/andrecronje/hashgraph/proxy/app"
)

func TestSocketProxyServer(t *testing.T) {
	clientAddr := "127.0.0.1:9990"
	proxyAddr := "127.0.0.1:9991"
	proxy := aproxy.NewSocketAppProxy(clientAddr, proxyAddr, 1*time.Second, common.NewTestLogger(t))
	submitCh := proxy.SubmitCh()

	tx := []byte("the test transaction")

	// Listen for a request
	go func() {
		select {
		case st := <-submitCh:
			// Verify the command
			if !reflect.DeepEqual(st, tx) {
				t.Fatalf("tx mismatch: %#v %#v", tx, st)
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("timeout")
		}
	}()

	// now client part connecting to RPC service
	// and calling methods
	dummyClient, err := NewDummySocketClient(clientAddr, proxyAddr, common.NewTestLogger(t))
	if err != nil {
		t.Fatal(err)
	}
	err = dummyClient.SubmitTx(tx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSocketProxyClient(t *testing.T) {
	clientAddr := "127.0.0.1:9992"
	proxyAddr := "127.0.0.1:9993"
	proxy := aproxy.NewSocketAppProxy(clientAddr, proxyAddr, 1*time.Second, common.NewTestLogger(t))

	dummyClient, err := NewDummySocketClient(clientAddr, proxyAddr, common.NewTestLogger(t))
	if err != nil {
		t.Fatal(err)
	}
	clientCh := dummyClient.hashgraphProxy.CommitCh()

	block := hashgraph.NewBlock(0, 1, [][]byte{[]byte("the test transaction")})

	// Listen for a request
	go func() {
		select {
		case commit := <-clientCh:
			if !reflect.DeepEqual(commit.Block, block) {
				t.Fatalf("block mismatch: %#v %#v", block, commit.Block)
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("timeout")
		}
	}()

	stateHash, err := proxy.CommitBlock(block)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("stateHash: %v", stateHash)
}
