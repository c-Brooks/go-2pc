package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/c-Brooks/zookeeper-demo/persistence"
)

// serve listens and serves HTTP on port 8080
func serve(kvrw persistence.KeyValueReadWriter) {
	h := handler{kvrw}
	http.ListenAndServe(":8080", h)
}

// handler implements the http.Handler interface (ServeHTTP)
type handler struct {
	kvrw persistence.KeyValueReadWriter
}

// handle requests to /
// body should contain {"key": key, "value": value}
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("handling %s request to %s", r.Method, r.URL)

	var m map[string]string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &m)
	if err != nil {
		logrus.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	key := m["key"]
	value := m["value"]

	switch r.Method {

	case http.MethodPost:
		t := persistence.NewTuple(key, value)
		err := h.kvrw.CreateResource(t)
		if err != nil {
			logrus.Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case http.MethodPut:
		// edit k/v
	case http.MethodDelete:
		// delete k/v
	}
}
