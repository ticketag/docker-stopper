package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/urfave/cli"

	"os"
	"sort"
	"strconv"
	"time"

	"github.com/ticketag/docker-stopper/commands"
	"github.com/ticketag/docker-stopper/restart_server"
)

func main() {
	runCliApp()
}

type ByCreated []types.Container

func (a ByCreated) Len() int           { return len(a) }
func (a ByCreated) Less(i, j int) bool { return a[i].Created < a[j].Created }
func (a ByCreated) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func runCliApp() {
	cliapp := cli.NewApp()
	cliapp.Name = "Docker restarter"
	cliapp.Usage = "Watcher per container selenium"
	cliapp.Version = "0.1"
	cliapp.Compiled = time.Now()
	cliapp.Copyright = "(c) Ticketag"
	cliapp.HelpName = "docker-restarter"
	cliapp.Flags = []cli.Flag{}
	serverArgs := commands.ServerArgs{}
	cliapp.Commands = []cli.Command{
		{
			Name:      "server",
			Aliases:   []string{"s"},
			Flags:     serverArgs.CliFlags(),
			Usage:     "server --host localhost --port 30001",
			ArgsUsage: "",
			Action: func(c *cli.Context) error {
				restart_server.Run(serverArgs.Host, serverArgs.Port, serverArgs.ScriptPath)
				return nil
			},
		},
		{
			Name:      "restarter",
			Aliases:   []string{"r"},
			Flags:     serverArgs.CliFlags(),
			Usage:     "restarter",
			ArgsUsage: "",
			Action: func(c *cli.Context) error {
				run()
				return nil
			},
		},
		{
			Name:      "install",
			Aliases:   []string{"i"},
			Usage:     "install",
			ArgsUsage: "",
			Action: func(c *cli.Context) error {
				fmt.Println("Installing...")
				commands.InstallService()
				return nil
			},
		},
	}

	//cliapp.Action = cliapp.Command("restarter").Action
	_ = cliapp.Run(os.Args)
}

func run() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	seleniumContainers := []types.Container{}
	for _, i := range containers {
		if i.Image == "elgalu/selenium:latest" {
			seleniumContainers = append(seleniumContainers, i)
		}
	}
	max := os.Getenv("MAX_SELENIUM_DOCKER")
	var max_num int
	if max == "" {
		max_num = 2
	} else {
		max_num, err = strconv.Atoi(max)
		if err != nil {
			max_num = 2
		}
	}
	to_delete := os.Getenv("TO_DELETE_DOCKER")
	var to_delete_num int
	if to_delete == "" {
		to_delete_num = 1
	} else {
		to_delete_num, err = strconv.Atoi(max)
		if err != nil {
			to_delete_num = 1
		}
	}
	if len(seleniumContainers) > max_num {
		sort.Sort(ByCreated(seleniumContainers))
		timeout := time.Duration(-1)
		for i := range seleniumContainers {
			if i < to_delete_num {
				err = cli.ContainerStop(ctx, seleniumContainers[i].ID, &timeout)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(fmt.Sprintf("DOCKER STOPPED: %v", seleniumContainers[i].ID))
			}
		}
	} else {
		fmt.Println("OK ", len(seleniumContainers))
	} /*
		httpClient := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := httpClient.Get("http://localhost:4444/wd/hub/status")
		if err != nil {
			resp, err = httpClient.Get("http://localhost:4444/grid/sessions?action=doCleanupActiveSessions")
			if err != nil {
				fmt.Println("ERROR IN ELIMINATING CONTAINERS")
				return
			}
			resp, _ := ioutil.ReadAll(resp.Body)
			if string(resp) == "SUCCESS" {
				fmt.Println("SUCCESS IN ELIMINATING ALL CONTAINERS")
			}
			return
		}*/
	return
}
