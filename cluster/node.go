package cluster

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/consul/api"

	pb "github.com/c-Brooks/go-2pc/protobuf"
)

// Node represents the state of a replica in the cluster
type Node struct {
	*api.Client                                // manages consul communication
	Hostname    string                         // os hostname
	IPAddr      string                         // ip address
	SDAddr      string                         // Service Discovery address
	SDKV        *api.KV                        // Local Key-Value store
	Clients     map[string]pb.RPCServiceClient // clients connected to RPC server
}

// RegisterService registers with the consul service discovery module.
// This implementation simply uses the key-value store. One major drawback is that when nodes crash. nothing is updated on the key-value store. Services are a better fit and should be used eventually.
func (n *Node) RegisterService(isMaster bool) {
	config := api.DefaultConfig()
	config.Address = n.SDAddr
	consul, err := api.NewClient(config)
	if err != nil {
		logrus.Fatalf("Unable to contact Service Discovery: %v", err)
	}

	kv := consul.KV()
	p := &api.KVPair{Key: n.Hostname, Value: []byte(n.IPAddr)}
	_, err = kv.Put(p, nil)
	if err != nil {
		logrus.Fatalf("Unable to register with Service Discovery: %v", err)
	}

	// embed consul.KV for future use
	n.SDKV = kv
	n.Client = consul

	logrus.Infoln("Successfully registered with Consul.")
}

// CheckServiceDiscoveryAndConnect checks the consul service for any new client connections.
// If any are found, current node registers as a client with the new node's gRPC server
func (n *Node) CheckServiceDiscoveryAndConnect() {
	logrus.Infof("listing k/v's from consul from host %s", n.Hostname)
	kvpairs, _, err := n.SDKV.List("Node", nil)
	if err != nil {
		logrus.Fatalf("could not fetch k/v pairs: %v", err)
		return
	}

	logrus.Infoln("Found nodes: ")
	for _, kventry := range kvpairs {
		if strings.Compare(kventry.Key, n.Hostname) == 0 {
			// we found ourselves
			continue
		}

		logrus.Infof("hostname: %s", kventry.Key)
		if n.Clients[kventry.Key] == nil {
			logrus.Infof("New member: %s", kventry.Key)

			// connection not established previously
			n.SetupGRPCClient(kventry.Key, string(kventry.Value))

		}
	}
}

// SetupGRPCClient initializes a gRPC client for contacting the server at addr.
func (n *Node) SetupGRPCClient(name string, addr string) {

	// setup connection with other node
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Save the client in the local k/v store
	n.Clients[name] = pb.NewRPCServiceClient(conn)

	// r, err := n.Clients[name].SayHello(context.Background(), &hs.HelloRequest{Name: n.Name})

	// todo: ping the new client

}

// RegisterGRPCService creates an agent and service entry in Consul
func (n *Node) RegisterGRPCService() error {
	rand.Seed(time.Now().UnixNano())
	sid := rand.Intn(65534)
	serviceID := n.Hostname + "-grpc-" + strconv.Itoa(sid)

	consulService := &api.AgentServiceRegistration{
		ID:   serviceID,
		Name: n.Hostname + "-grpc",
		Port: 50051,
	}

	logrus.Infof("registered gRPC service for %s", n.Hostname)

	ag := n.Agent()
	err := ag.ServiceRegister(consulService)
	if err != nil {
		logrus.Errorf("FUCK %v", err)
	}
	return err
}

// Healthcheck reports health of the cluster
func (n *Node) Healthcheck() error {
	_, err := n.Agent().Services()
	if err != nil {
		return err
	}

	// ping each service
	// for _, service := range as {
	// 	logrus.Infoln(service.Service)
	// 	logrus.Infoln(service.Address)
	// 	logrus.Infoln(service.Port)
	// }

	hcr := &pb.HealthCheckReq{}
	ctx := context.Background()
	for _, c := range n.Clients {
		resp, err := c.HealthCheck(ctx, hcr)
		// idk what sort of codes to expect
		if err != nil || resp.GetStatusCode() == 500 {
			return nil
		}
	}

	return nil
}
