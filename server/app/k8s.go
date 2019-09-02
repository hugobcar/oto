package app

import (
	"flag"
	"fmt"
	"io"
	"time"

	// "github.com/hugobcar/oto/server/app"

	"github.com/pkg/errors"
	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

// Shamelessly adapted from Kubernetes
func k8sPodToAppPod(pod *k8sv1.Pod) *Pod {
	p := &Pod{Name: pod.Name}
	if pod.Status.StartTime != nil {
		p.Age = int64(time.Since(pod.Status.StartTime.Time))
	}
	p.State = string(pod.Status.Phase)
	for _, status := range pod.Status.InitContainerStatuses {
		if status.State.Terminated == nil || status.State.Terminated.ExitCode != 0 {
			p.State = "Initializing"
			return p
		}
	}
	nready := 0
	for _, status := range pod.Status.ContainerStatuses {
		if status.State.Waiting != nil {
			p.State = status.State.Waiting.Reason
		} else if status.State.Terminated != nil {
			p.State = status.State.Terminated.Reason
		} else if status.State.Running != nil && status.Ready {
			nready++
		}
		p.Restarts += status.RestartCount
	}
	if len(pod.Status.ContainerStatuses) == nready && nready != 0 {
		p.Ready = true
	}
	if pod.DeletionTimestamp != nil {
		p.State = "Terminating"
	}
	return p
}

func appPodListOptsToK8s(opts *PodListOptions) *metav1.ListOptions {
	var k8sOpts metav1.ListOptions

	if opts.PodName != "" {
		k8sOpts.FieldSelector = fmt.Sprintf("metadata.name=%s", opts.PodName)
	}

	return &k8sOpts
}

func buildClient() (kubernetes.Interface, error) {
	kubeconfig := flag.String("kubeconfig", "/home/hugo.carvalho/.kube/config", "absolute path to the kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "create k8s client failed")
	}
	return c, nil
}

// PodList -
func PodList(namespace string, opts *PodListOptions) ([]*Pod, error) {
	kc, err := buildClient()
	if err != nil {
		return nil, err
	}
	k8sOpts := appPodListOptsToK8s(opts)
	podList, err := kc.CoreV1().Pods(namespace).List(*k8sOpts)
	if err != nil {
		return nil, err
	}

	pods := make([]*Pod, len(podList.Items))
	for i, pod := range podList.Items {
		pods[i] = k8sPodToAppPod(&pod)
	}
	return pods, nil
}

// PodLogs -
func PodLogs(namespace string, podName string, opts *LogOptions) (io.ReadCloser, error) {
	kc, err := buildClient()
	if err != nil {
		return nil, err
	}
	request := kc.CoreV1().Pods(namespace).GetLogs(
		podName,
		&k8sv1.PodLogOptions{
			Follow:    opts.Follow,
			TailLines: &opts.TailLines,
			Previous:  opts.Previous,
			Container: opts.Container,
		},
	)

	return request.Stream()
}
