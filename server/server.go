package server

import (
	"encoding/json"

	"fmt"
	"github.com/YangKeao/reserve/config"
	"net/http"

	"github.com/juju/errors"
	"github.com/ngaut/log"

	k8s "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Server struct {
	client *kubernetes.Clientset
	config *config.Config
}

type ExtenderArgs struct {
	Pod k8s.Pod
	Nodes *k8s.NodeList
	NodeNames *[]string
}

type FailedNodesMap map[string]string

type FilterResult struct {
	Nodes *k8s.NodeList `json:"nodes,omitempty"`
	NodeNames *[]string `json:"nodenames,omitempty"`
	FailedNodes FailedNodesMap `json:"failedNodes,omitempty"`
	Error string `json:"error,omitempty"`
}

func New(config *config.Config) (*Server, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", config.Kube.ConfigPath)
	if err != nil {
		return nil, errors.Trace(err)
	}

	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, errors.Trace(err)
	}

	server := Server {
		client: client,
		config: config,
	}
	return &server, nil
}

func (server Server) ListenAndServe() {
	host := fmt.Sprintf("%s:%d", server.config.Server.Host, server.config.Server.Port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Receive Request: %s", r.URL)


		var args ExtenderArgs
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&args); err != nil {
			http.Error(w, "decode request error", http.StatusBadRequest)
		}

		failedNodes := make(map[string]string);
		for _, node := range args.Nodes.Items {
			failedNodes[node.Name] = "REFUSE!"
		}
		resp := FilterResult{
			Nodes: nil,
			NodeNames: nil,
			FailedNodes: failedNodes,
		}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(resp); err != nil {
			http.Error(w, "encode response error", http.StatusInternalServerError)
		}
	})

	log.Infof("Serve on %s", host)

	http.ListenAndServe(host, nil)
}
