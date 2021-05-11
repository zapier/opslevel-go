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

	serviceCmd            = app.Command("service", "")
	serviceGetCmd         = serviceCmd.Command("get", "get a service by name / alias")
	serviceGetCmdAliasArg = serviceGetCmd.Arg("alias", "service alias").Required().String()

	serviceTagCmd            = serviceCmd.Command("tag", "manipulate service tags")
	serviceTagAddCmd         = serviceTagCmd.Command("add", "Add service tag")
	serviceTagAddCmdAliasArg = serviceTagAddCmd.Arg("alias", "service alias").Required().String()
	serviceTagAddCmdTagsArg  = serviceTagAddCmd.Arg("k=v", "key=value tag").Required().StringMap()

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
	case serviceTagAddCmd.FullCommand():
		for key, value := range *serviceTagAddCmdTagsArg {
			tag, err := client.CreateTag(context.Background(), key, value, *serviceTagAddCmdAliasArg, "Service")
			if err != nil {
				log.Fatalln(err)
			}
			pretty.Println(tag)
		}

	case serviceGetCmd.FullCommand():
		svc, err := client.GetService(context.Background(), *serviceGetCmdAliasArg)
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
