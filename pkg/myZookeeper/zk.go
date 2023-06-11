package myZookeeper

import (
	"time"

	"golang.org/x/xerrors"

	"github.com/samuel/go-zookeeper/zk"
)

var myZookeeper *MyZookeeper

type MyZookeeper struct {
	DefaultConn *zk.Conn
}

func NewMyZookeeper(servers []string, timeout time.Duration) *MyZookeeper {
	c, err := ConnectZk(servers, timeout)
	if err != nil {
		panic(err)
	}
	myZookeeper = &MyZookeeper{
		DefaultConn: c,
	}
	return myZookeeper
}

func ConnectZk(servers []string, timeout time.Duration) (*zk.Conn, error) {
	formatServers := zk.FormatServers(servers)

	if timeout == 0 {
		timeout = time.Second * 2
	}
	c, _, err := zk.Connect(formatServers, timeout)
	// c, session, err := zk.Connect([]string{"127.0.0.1:2182"}, time.Second*1)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// CreateNode flags  0 持久节点 1 临时节点 2 有序节点
func (p *MyZookeeper) CreateNode(nodePath string, data []byte, flags int32, acl []zk.ACL) error {
	if err := p.checkMyZookeeper(); err != nil {
		return err
	}
	_, err := p.DefaultConn.Create(nodePath, data, flags, acl)
	if err != nil {
		return xerrors.Errorf("MyZookeeper CreateNode : %w", err)
	}
	return nil
}

func (p *MyZookeeper) GetNodeChildren(path string) ([]string, error) {
	if err := p.checkMyZookeeper(); err != nil {
		return []string{}, err
	}
	res, _, err := p.DefaultConn.Children(path)
	if err != nil {
		return []string{}, xerrors.Errorf("MyZookeeper GetNodeChildren : %w", err)
	}
	return res, nil
}

func (p *MyZookeeper) GetNodeChildrenAndWatcher(path string) ([]string, <-chan zk.Event, error) {
	if err := p.checkMyZookeeper(); err != nil {
		return []string{}, nil, err
	}
	res, _, watcher, err := p.DefaultConn.ChildrenW(path)
	if err != nil {
		return []string{}, nil, xerrors.Errorf("MyZookeeper GetNodeChildrenAndWatcher : %w", err)
	}
	return res, watcher, nil
}

func (p *MyZookeeper) GetNode(path string) ([]byte, error) {
	if err := p.checkMyZookeeper(); err != nil {
		return nil, err
	}
	data, _, err := p.DefaultConn.Get(path)
	if err != nil {
		return nil, xerrors.Errorf("MyZookeeper GetNode : %w", err)
	}
	return data, nil

}

func (p *MyZookeeper) DeleteNode(nodePath string, version int32) error {
	if err := p.checkMyZookeeper(); err != nil {
		return err
	}
	err := p.DefaultConn.Delete(nodePath, version)
	if err != nil {
		return xerrors.Errorf("MyZookeeper DeleteNode : %w", err)
	}
	return nil
}

func (p *MyZookeeper) Exists(path string) (bool, error) {
	if err := p.checkMyZookeeper(); err != nil {
		return false, err
	}
	res, _, err := p.DefaultConn.Exists(path)
	if err != nil {
		return false, xerrors.Errorf("MyZookeeper Exists : %w", err)
	}
	return res, nil
}

func (p *MyZookeeper) checkMyZookeeper() error {
	if p == nil {
		return xerrors.Errorf("MyZookeeper checkMyZookeeper MyZookeeper is nil")
	}
	if p.DefaultConn == nil {
		return xerrors.Errorf("MyZookeeper checkMyZookeeper DefaultConn is nil")
	}
	return nil
}
