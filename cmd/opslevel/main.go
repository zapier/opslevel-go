package main

import (
	"context"
	"log"
	"os"

	"github.com/kr/pretty"

	"github.com/zapier/opslevel-go"

	"github.com/alecthomas/kingpin"
)

var (
	app = kingpin.New("opslevel", "OpsLevel.com command line utility")

	serviceCmd    = app.Command("service", "")
	serviceGetCmd = serviceCmd.Command("get", "get a service by name / alias")
	serviceAlias  = serviceGetCmd.Arg("alias", "service alias").Required().String()

	teamCmd    = app.Command("team", "")
	teamGetCmd = teamCmd.Command("get", "get a team by name / alias")
	teamAlias  = teamGetCmd.Arg("alias", "team alias").Required().String()

	version = "dev"
)

func main() {
	app.Version("opslevel version: " + version)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	var cmd = kingpin.MustParse(app.Parse(os.Args[1:]))

	var authToken = os.Getenv("OPSLEVEL_TOKEN")
	client := opslevel.NewClient(authToken)

	switch cmd {
	case serviceGetCmd.FullCommand():
		svc, err := client.GetService(context.Background(), *serviceAlias)
		if err != nil {
			log.Fatalln(err)
		}

		pretty.Println(svc)

	case teamGetCmd.FullCommand():
		team, err := client.GetTeam(context.Background(), *teamAlias)
		if err != nil {
			log.Fatalln(err)
		}

		pretty.Println(team)
	}
}
