package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	srv := &http.Server{Addr: ":30001"}
	script_path := "/home/ubuntu/dockerimages/selenium/zalenium/run.sh"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("/bin/sh", script_path)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		output := out.String()
		if err != nil {

			w.WriteHeader(500)
			fmt.Fprint(w, output, err)
			return
		}

		w.WriteHeader(200)
		fmt.Fprint(w, output)
		return
	})
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}

}
