package main

import (
	"bufio"
	"fmt"

	// "github.com/hugobcar/oto/pkg/protobuf/app"
	"github.com/hugobcar/oto/server/app"
)

// Logs -
func Logs(appName string, opts *app.LogOptions) (map[string]string, error) {
	namespace := "monitoring"
	podName := []string{"tiamat-67556f7c5b-rpjvh", "tiamat-67556f7c5b-rpjvh"}

	for _, pName := range podName {
		// logChan := make(chan map[string]string, 1)
		fmt.Println(pName)

		go getLogs(pName, namespace, opts)

		// select {
		// case result := <-logChan:
		// 	return result, nil
		// 	// case <-time.After(2 * time.Minute):
		// 	// 	return nil, fmt.Errorf("Timeout when reading the logs of all pds created by Oto job '%s'", namespace)
		// }
	}
	return nil, nil
}

func main() {
	opts := &app.LogOptions{
		TailLines: 10,
		Follow:    true,
		PodName:   "tiamat-67556f7c5b-rpjvh",
		Previous:  false,
		//Container: "absolute",
	}

	// namespace := "default"
	// fmt.Println(app.PodList(namespace))

	logList, err := Logs("tiamat-67556f7c5b-rpjvh", opts)
	if err != nil {
		fmt.Printf("Could not retrieve the logs of the pods belonging to Oto job '%s'", err.Error())
		logList = map[string]string{}
	}
	for podName, podLogs := range logList {
		fmt.Printf("Logs of Pod '%s' belonging to Oto job:\n%s", podName, podLogs)
	}
}

func getLogs(podName, namespace string, opts *app.LogOptions) {
	// var logList = map[string]string{}
	logsBuffer, err := app.PodLogs(namespace, podName, opts)

	sc := bufio.NewScanner(logsBuffer)
	for sc.Scan() {
		fmt.Println(podName)
		fmt.Println(sc.Text())
	}

	if err != nil {
		fmt.Printf("Could not retrieve the logs of Oto job pod %s: '%v'", podName, err)
	}
	// logList[podName] = logsBuffer.String()
	//logChan <- logList
}
