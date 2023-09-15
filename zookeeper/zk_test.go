package zookeeper

import (
	"fmt"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var p *ZooKeeper

func init() {
	p = InitZooKeeper([]string{"43.138.39.90:2181", "43.138.39.90:2182", "43.138.39.90:2183"}, time.Second*2)
}

func TestZooKeeper_Create(t *testing.T) {
	err := p.Create("/test/001", []byte("hello2"), 0, zk.WorldACL(zk.PermAll))
	fmt.Println(err)
}

func TestZooKeeper_Set(t *testing.T) {
	stat, err := p.Set("/test", []byte("hello2"), -1)
	fmt.Println(stat, err)
}

func TestGetWater(t *testing.T) {
	w, stat, events, err := p.GetW("/test")
	if err != nil {
		panic(err)
	}
	fmt.Println("Get data", w)
	fmt.Println("Get stat", stat)
	for {
		select {
		case r, _ := <-events:
			fmt.Println(r.Err, r.Type.String(), r.Path, r.Server, r.State.String())
		}
	}
}
