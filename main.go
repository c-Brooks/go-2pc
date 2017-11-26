package main

import (
	"os"

	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/c-Brooks/zookeeper-demo/persistence"
)

const xactQueueBuffSize = 10

func main() {

	pwd, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("could not get PWD: %v", err)
	}

	path := pwd + "/tmp/zk"

	// store
	store, err := persistence.NewKeyValueStore("file", map[string]string{"filepath": path})
	if err != nil {
		logrus.Warnf("could not create store: %v", err)
	}

	// http server
	go serve(store)

	isMaster := flag.Bool("master", false, "whether or not this replica is the master")
	if *isMaster {
		// TODO
		// master server
	}

	// keep process alive
	for {
	}

}

// func transactionQueue(xactQueue chan persistence.Tuple) {
// 	var t persistence.Tuple

// 	logrus.Infof("creating k/v store at /tmp/zk")

// 	store, err := persistence.NewKeyValueStore("file", map[string]string{"filepath": "/tmp/zk"})
// 	if err != nil {
// 		logrus.Fatalf("could not create key-value store: %v", err)
// 	}

// 	select {

// 	case t = <-xactQueue:
// 		store.EditResource(t)
// 	}

// }

// func (s *grpcServer) startGrpcServer() {
// 	for {

// 	}
// }

// func serveGrpc(port int) grpc.Server {

// 	// conn, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
// 	s := pb.NewRPCServiceClient()

// 	grpcServer := grpc.NewServer()
// 	pb.RegisterRPCServiceServer(grpcServer, s)

// 	grpcServer.Serve(conn)
// 	return grpcServer
// }
