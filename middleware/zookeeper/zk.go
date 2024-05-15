package zookeeper

import (
	"time"

	"golang.org/x/xerrors"

	"github.com/samuel/go-zookeeper/zk"
)

var defaultZk *ZooKeeper

type ZooKeeper struct {
	DefaultConn *zk.Conn
}

func InitZooKeeper(servers []string, timeout time.Duration) *ZooKeeper {
	c, err := ConnectZk(servers, timeout)
	if err != nil {
		panic(err)
	}
	defaultZk = &ZooKeeper{
		DefaultConn: c,
	}
	return defaultZk
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

// Create  flags  0 持久节点 1 临时节点 2 有序节点
func (p *ZooKeeper) Create(nodePath string, data []byte, flags int32, acl []zk.ACL) error {
	if err := p.checkZooKeeper(); err != nil {
		return err
	}
	_, err := p.DefaultConn.Create(nodePath, data, flags, acl)
	if err != nil {
		return xerrors.Errorf("ZooKeeper CreateNode : %w", err)
	}
	return nil
}

func (p *ZooKeeper) Set(path string, data []byte, version int32) (*zk.Stat, error) {
	return p.DefaultConn.Set(path, data, version)
}

func (p *ZooKeeper) Children(path string) ([]string, error) {
	if err := p.checkZooKeeper(); err != nil {
		return []string{}, err
	}
	res, _, err := p.DefaultConn.Children(path)
	if err != nil {
		return []string{}, xerrors.Errorf("ZooKeeper GetNodeChildren : %w", err)
	}
	return res, nil
}

func (p *ZooKeeper) GetW(path string) ([]byte, *zk.Stat, <-chan zk.Event, error) {
	return p.DefaultConn.GetW(path)
}

func (p *ZooKeeper) ChildrenW(path string) ([]string, <-chan zk.Event, error) {
	if err := p.checkZooKeeper(); err != nil {
		return []string{}, nil, err
	}
	res, _, watcher, err := p.DefaultConn.ChildrenW(path)
	if err != nil {
		return []string{}, nil, xerrors.Errorf("ZooKeeper GetNodeChildrenAndWatcher : %w", err)
	}
	return res, watcher, nil
}

func (p *ZooKeeper) Get(path string) ([]byte, error) {
	if err := p.checkZooKeeper(); err != nil {
		return nil, err
	}
	data, _, err := p.DefaultConn.Get(path)
	if err != nil {
		return nil, xerrors.Errorf("ZooKeeper GetNode : %w", err)
	}
	return data, nil

}

func (p *ZooKeeper) Delete(nodePath string, version int32) error {
	if err := p.checkZooKeeper(); err != nil {
		return err
	}
	err := p.DefaultConn.Delete(nodePath, version)
	if err != nil {
		return xerrors.Errorf("ZooKeeper DeleteNode : %w", err)
	}
	return nil
}

func (p *ZooKeeper) Exists(path string) (bool, error) {
	if err := p.checkZooKeeper(); err != nil {
		return false, err
	}
	res, _, err := p.DefaultConn.Exists(path)
	if err != nil {
		return false, xerrors.Errorf("ZooKeeper Exists : %w", err)
	}
	return res, nil
}

func (p *ZooKeeper) checkZooKeeper() error {
	if p == nil {
		return xerrors.Errorf("ZooKeeper checkZooKeeper ZooKeeper is nil")
	}
	if p.DefaultConn == nil {
		return xerrors.Errorf("ZooKeeper checkZooKeeper DefaultConn is nil")
	}
	return nil
}
