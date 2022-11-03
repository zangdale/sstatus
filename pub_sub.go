package sstatus

import "context"

type ConnInter interface {
	Send(k string, data []byte) error
}

type dataAndConns struct {
	datas  map[string]chan []byte
	conns  map[string][]ConnInter
	cancel map[string]context.CancelFunc
}
