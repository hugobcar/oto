package main

import (
	"bufio"
	"fmt"
	"time"

	"github.com/hugobcar/oto/server/app"
)

// Logs -
func Logs(appName string, opts *app.LogOptions) (map[string]string, error) {
	namespace := "default"
	podName := "counter"

	logChan := make(chan map[string]string, 1)
	go func() {
		var logList = map[string]string{}
		logsBuffer, err := app.PodLogs(namespace, podName, opts)

		sc := bufio.NewScanner(logsBuffer)
		for sc.Scan() {
			fmt.Println(sc.Text())
		}

		if err != nil {
			fmt.Printf("Could not retrieve the logs of Terraform job pod %s: '%v'", podName, err)
		}
		// logList[podName] = logsBuffer.String()
		logChan <- logList
	}()

	select {
	case result := <-logChan:
		return result, nil
	case <-time.After(2 * time.Minute):
		return nil, fmt.Errorf("Timeout when reading the logs of all pds created by Terraform job '%s'", namespace)
	}
}

func main() {
	opts := &app.LogOptions{
		TailLines: 10,
		Follow:    true,
		PodName:   "counter",
		Previous:  false,
		Container: "count",
	}

	logList, err := Logs("counter", opts)
	if err != nil {
		fmt.Printf("Could not retrieve the logs of the pods belonging to Terraform job '%s'", err.Error())
		logList = map[string]string{}
	}
	for podName, podLogs := range logList {
		fmt.Printf("Logs of Pod '%s' belonging to Terraform job:\n%s", podName, podLogs)
	}
}
