package commands

import (
	"github.com/urfave/cli"
)

type ServerArgs struct {
	Host       string
	Port       uint
	ScriptPath string
}

func (s *ServerArgs) CliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "host,k", Destination: &s.Host, EnvVar: "STOPPER_HOST"},
		cli.UintFlag{Name: "port,p", Destination: &s.Port, Value: 30001, EnvVar: "STOPPER_PORT"},
		cli.StringFlag{Name: "path", Destination: &s.ScriptPath, Value: "/home/ubuntu/dockerimages/selenium/selenoid/run.sh", EnvVar: "SELENIUM_SCRIPT_PATH"},
	}
}
