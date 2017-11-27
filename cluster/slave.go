package cluster

import (
	"os"

	"github.com/Sirupsen/logrus"

	"google.golang.org/grpc"
)

// Slave contains a gRPC client pointing to the master
type Slave struct {
	*grpc.ClientConn
}

// New returns a new instance of Slave with an active gRPC connection
func New() Slave {
	masterAddr := os.Getenv("master_addr")
	conn, err := grpc.Dial(masterAddr)
	if err != nil {
		logrus.Fatalf("could not dial to master: %v", err)
	}

	return Slave{conn}
}
