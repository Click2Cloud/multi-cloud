package ceph_adapter

type CephAdapterInterface interface {
	Remove(name string) error
}
