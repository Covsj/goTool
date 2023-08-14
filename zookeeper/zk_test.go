package zookeeper

import (
	"fmt"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var p *MyZookeeper

func init() {
	p = NewMyZookeeper([]string{"43.138.39.90"}, time.Second*2)
}

func TestMyZookeeper_CreateNode(t *testing.T) {
	fmt.Println(p.CreateNode("/test", []byte("test"), 0, zk.WorldACL(zk.PermAll)))
}

func TestMyZookeeper_GetNode(t *testing.T) {
	res, err := p.GetNode("/test")
	fmt.Println(string(res), err)
}

func TestMyZookeeper_Exists(t *testing.T) {
	exists, err := p.Exists("/test")
	fmt.Println(exists, err)

	exists, err = p.Exists("/test2")
	fmt.Println(exists, err)
}

func TestMyZookeeper_DeleteNode(t *testing.T) {
	err := p.DeleteNode("/test", -1)
	fmt.Println(err)
}
