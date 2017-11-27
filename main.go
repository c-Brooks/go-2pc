package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/c-Brooks/go-2pc/cluster"
	"github.com/c-Brooks/go-2pc/storage"
)

func main() {

	pwd, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("could not get PWD: %v", err)
	}

	dataDir := pwd + "/data/"
	tmpDir := pwd + "/tmp/"

	// store
	store, err := storage.NewKeyValueStore("file", map[string]string{
		"dataDir": dataDir,
		"tmpDir":  tmpDir,
	})
	if err != nil {
		logrus.Warnf("could not create store: %v", err)
	}

	_ = store

	hostname, _ := os.Hostname()
	listenaddr := os.Getenv("addr")
	sdaddress := os.Getenv("CONSUL_ADDRESS")

	currentNode := cluster.Node{
		Hostname: hostname,
		IPAddr:   listenaddr,
		SDAddr:   sdaddress,
		Clients:  nil,
	}

	// http server
	go ServeHTTP(store, 8080, &currentNode)

	isMaster := flag.Bool("master", false, "whether or not this replica is the master")
	currentNode.RegisterService(*isMaster)

	currentNode.CheckServiceDiscoveryAndConnect()

	err = currentNode.RegisterGRPCService()
	if err != nil {
		logrus.Fatalf("could not register service: %v", err)
	}

	// keep process alive
	for {
	}

}

//
//
// ========================= HTTP ================== //
//
//

// ServeHTTP listens and serves HTTP on port 8080
func ServeHTTP(kvrw storage.TwoPhaseReadWriter, port int, currentNode *cluster.Node) {
	// mux := http.NewServeMux()
	h := httpHandler{kvrw, currentNode}
	http.ListenAndServe(":8080", h)
}

// handler implements the http.Handler interface (ServeHTTP)
type httpHandler struct {
	kvrw storage.TwoPhaseReadWriter
	node *cluster.Node
}

func handleHealthz(w http.ResponseWriter, r *http.Request, currentNode *cluster.Node) {
	currentNode.Healthcheck()
	hostName, _ := os.Hostname()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hostName))
}

// handle requests to /
// body should contain {"key": key, "value": value}
func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("handling %s request to %s", r.Method, r.URL)

	switch r.URL.Path {
	case "/healthz":
		logrus.Infof("handling /healthz | %s", r.URL.Path)
		handleHealthz(w, r, h.node)

	default:
		logrus.Infof("handling /something | %s", r.URL.Path)
		switch r.Method {

		case http.MethodPost:
			t, err := parseBody(w, r)
			err = h.kvrw.CreateResource(*t)
			if err != nil {
				logrus.Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		case http.MethodPut:
			t, err := parseBody(w, r)
			err = h.kvrw.EditResource(*t)
			if err != nil {
				logrus.Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		case http.MethodDelete:
			t, err := parseBody(w, r)
			key := t.GetKey()

			err = h.kvrw.DeleteResource(key)
			if err != nil {
				logrus.Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case http.MethodGet:
			t, err := parseBody(w, r)
			key := t.GetKey()

			val, err := h.kvrw.ReadResource(key)
			if err != nil {
				logrus.Warn(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(val))
		}

	}
}

func parseBody(w http.ResponseWriter, r *http.Request) (*storage.Tuple, error) {
	var m map[string]string

	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	key := m["key"]
	value := m["value"]

	return storage.NewTuple(key, value), nil
}
