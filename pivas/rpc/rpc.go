package rpc

import (
	"log"
	"marmota/pivas/cc"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/hashicorp/net-rpc-msgpackrpc"
)

func Start() {
	addr := cc.Config().Listen

	server := rpc.NewServer()
	err := server.Register(new(Agent))
	if err != nil {
		log.Fatalln(err)
	}
	err = server.Register(new(Hbs))
	if err != nil {
		log.Fatalln(err)
	}

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalln("listen error:", e)
	} else {
		log.Println("listening", addr)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("listener accept fail:", err)
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		}
		msgpackrpc.NewServerCodec(conn)
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
