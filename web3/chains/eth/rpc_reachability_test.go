package eth

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Covsj/goTool/web3/chains/basic"
)

type dddd struct {
}

func (d *dddd) ReachabilityDidReceiveNode(tester *basic.ReachMonitor, latency *basic.RpcLatency) {
	fmt.Printf(".... delegate did receive height %v, latency %v, node %v\n", latency.Height, latency.Latency, latency.RpcUrl)
}

func (d *dddd) ReachabilityDidFailNode(tester *basic.ReachMonitor, latency *basic.RpcLatency) {
	fmt.Printf(".... delegate did fail height %v, latency %v, node %v\n", latency.Height, latency.Latency, latency.RpcUrl)
	// tester.StopConnectivity()
}

func (d *dddd) ReachabilityDidFinish(tester *basic.ReachMonitor, overview string) {
	fmt.Printf(".... delegate did finish %v\n", overview)
}

func TestRpcReachability_Test(t *testing.T) {
	reach := NewRpcReachability()
	monitor := basic.NewReachMonitorWithReachability(reach)
	monitor.ReachCount = 3
	monitor.Delay = 3000
	monitor.Timeout = 1500
	t.Log(reach)

	rpcUrls := []string{rpcs.ethereumProd.url, rpcs.binanceTest.url}
	rpcListString := strings.Join(rpcUrls, ",")
	// res := reach.StartConnectivitySync(rpcListString)
	// t.Log(res)

	delegate := &dddd{}
	monitor.StartConnectivityDelegate(rpcListString, delegate)
}
