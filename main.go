package main

import (
	"fmt"
	"io"

	"github.com/hugobcar/oto/server/app"
)

type K8sOperations interface {
	PodLogs(namespace, podName string, opts *app.LogOptions) (io.ReadCloser, error)
	Logs(appName string, opts *app.LogOptions) (io.ReadCloser, error)
}

type AppOperations struct {
	kops K8sOperations
}

// Logs -
func (ops *AppOperations) Logs(appName string, opts *app.LogOptions) (io.ReadCloser, error) {
	namespace := "default"
	podName := "counter"

	logs, err := ops.kops.PodLogs(namespace, podName, opts)
	if err != nil {
		fmt.Sprintf("streaming logs from pod %s", podName)
		return logs, err
	}
	fmt.Sprintf("streaming logs from pod %s", podName)
	defer logs.Close()

	return logs, err

	// scanner := bufio.NewScanner(logs)
	// for scanner.Scan() {
	// 	fmt.Fprintf(w, "[%s] - %s\n", podName, scanner.Text())
	// }
	// if err := scanner.Err(); err != nil {
	// 	fmt.Printf("streaming logs from pod %s", podName)
	// 	fmt.Println(err)
	// }
}

func main() {
	opts := &app.LogOptions{
		Lines:     2,
		Follow:    true,
		PodName:   "counter",
		Previous:  true,
		Container: "counter",
	}

	var ops *AppOperations

	rc, err := ops.kops.Logs("counter", opts)
	if err != nil {
		fmt.Println(err)
	}
	defer rc.Close()

	fmt.Println(rc)
	// chLogs, errCh := goutil.LineGenerator(rc)

	// fmt.Println(chLogs)
	// fmt.Println(errCh)
}
