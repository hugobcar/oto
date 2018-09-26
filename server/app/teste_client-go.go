package app

import (
	"flag"
	"io"
	"time"

	"github.com/pkg/errors"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client - Struct
type Client struct {
	conf          *restclient.Config
	podRunTimeout time.Duration
	ingress       bool
	fake          kubernetes.Interface
	testing       bool
}

func (k *Client) buildClient() (kubernetes.Interface, error) {
	kubeconfig := flag.String("kubeconfig", "/home/hugo/.kube/config", "absolute path to the kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	if k.testing {
		k.fake = fake.NewSimpleClientset()
		return k.fake, nil
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "create k8s client failed")
	}
	return c, nil
}

// PodLogs - Get logs by PodLogs k8s
func (k *Client) PodLogs(namespace string, podName string, opts *LogOptions) (io.ReadCloser, error) {
	kc, err := k.buildClient()
	if err != nil {
		return nil, err
	}
	req := kc.CoreV1().Pods(namespace).GetLogs(
		podName,
		&k8sv1.PodLogOptions{
			Follow:    opts.Follow,
			TailLines: &opts.Lines,
			Previous:  opts.Previous,
			Container: opts.Container,
		},
	)

	return req.Stream()
}
