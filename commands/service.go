package commands

import (
	"fmt"
	"os"
	"os/exec"
)

const svc = `[Unit]
Description=Selenium restarter service.

[Service]
Type=simple
ExecStart=/usr/bin/%s

[Install]
WantedBy=multi-user.target
`

func InstallService() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	cpCmd := exec.Command("cp", "-rf", ex, "/usr/bin/")
	err = cpCmd.Run()
	if err != nil {
		panic(err)
	}
	serviceStr := fmt.Sprintf(svc, ex)
	fmt.Println(serviceStr)
}
