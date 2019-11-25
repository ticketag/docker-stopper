package restart_server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"

	"github.com/astaxie/beego/config/env"
)

func main() {
	home := env.Get("HOME", "/home/ubuntu")
	Run("", 30001, path.Join(home, "dockerimages/selenium/zalenium/run.sh"))
}
func Run(Host string, Port uint, scriptPath string) {
	serverAddr := fmt.Sprintf("%s:%d", Host, Port)
	srv := &http.Server{Addr: serverAddr}
	fmt.Println("Server running on: ", serverAddr)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("/bin/sh", scriptPath)
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
