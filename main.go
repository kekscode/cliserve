package main

import (
	"net/http"
	"os/exec"
	"strings"
	"encoding/json"
	"io/ioutil"
	"log"
)

/*
TODO: Include stderr in output
*/

type result struct {
	Command string `json:"command"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
}

func commandCall(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query()

	cmd := strings.Fields(message.Get("cmd"))[0]
	params := strings.Fields(message.Get("cmd"))[1:]

	// Get full path of command
	cmd, lookErr := exec.LookPath(cmd)
	if lookErr != nil {
		panic(lookErr)
	}

	// Prepare full command
	executing := exec.Command(cmd, strings.Join(params, " "))

	// Prepare pipes
	cmdIn, _ := executing.StdinPipe()
	cmdOut, _ := executing.StdoutPipe()

	// Start and read output
	executing.Start()
	executingBytes, _ := ioutil.ReadAll(cmdOut)
	executing.Wait()

	// Close input pipe
	cmdIn.Close()

	log.Println(string(executingBytes))

	res := result{
		"stdout",
		string(executingBytes),
		"err",
	}

	re, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(re)
}

func main() {
	http.HandleFunc("/exec", commandCall)
	if err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil); err != nil {
		panic(err)
	}
}
