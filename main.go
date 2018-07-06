package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

/*
TODO: Include stderr in output
FIX: single command exec "ls" (no args)
*/

type result struct {
	Command string `json:"command"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
}

func commandCall(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query()

	//	cmd := strings.Fields(message.Get("cmd"))
	cmdln := message.Get("cmd")
	log.Println(cmdln)

	// Get full path of command
	cmd, lookErr := exec.LookPath(strings.Fields(cmdln)[0])
	if lookErr != nil {
		panic(lookErr)
	}

	//	log.Println(args)
	executing := exec.Command(cmd, strings.Fields(cmdln)[1:]...)

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
