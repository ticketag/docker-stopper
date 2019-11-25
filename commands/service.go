package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

const svc = `[Unit]
Description=Selenium restarter service.

[Service]
Type=simple
WorkingDirectory=%s
ExecStart=%s
Restart=on-failure

[Install]
WantedBy=multi-user.target
`
const shFile = `#!/bin/bash
echo "Starting server..."
%s server --path %s
`
const systemDPath = `/etc/systemd/system/`

func InstallService(args *ServerArgs) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exName := filepath.Base(ex)
	home := os.Getenv("HOME")
	/*fmt.Println("Copy: ",ex)
	cpCmd := exec.Command("cp", ex, "/usr/bin/")
	err = cpCmd.Run()
	if err != nil {
		log.Fatal(err)
	}*/

	shScript := fmt.Sprintf(shFile, ex, args.ScriptPath)

	scriptPath := path.Join(home, exName+".sh")
	svcManifest := fmt.Sprintf(svc, home, scriptPath)
	installPath := path.Join(systemDPath, exName+".service")
	fmt.Println("Install service in: ", installPath)
	if err := ioutil.WriteFile(scriptPath, []byte(fmt.Sprintf(shScript)), 0644); err != nil {
		log.Fatal(err)
	}
	if err := os.Chmod(scriptPath, 0755); err != nil {
		log.Fatal(err)
	}
	fmt.Println("installed script: ", scriptPath)
	if err := ioutil.WriteFile(installPath, []byte(svcManifest), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Service installed")
	fmt.Println("Start service with: systemctl start " + exName + ".service")
}
