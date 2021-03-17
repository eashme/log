package coord_log

import "net"

type KafkaCli struct {
	addr net.Addr

}


func NewKafkaCli() *KafkaCli{
	cli := new(KafkaCli)
	// to produce messages


	cli.addr = &net.TCPAddr{
		IP:   net.IP("127.0.0.1"),
		Port: 9092,
	}


	return cli
}