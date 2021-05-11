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
	app = kingpin.New("opslevel", "OpsLevel.com test client")

	// Services
	serviceCmd = app.Command("service", "")
	// Create Service
	serviceCreateCmd        = serviceCmd.Command("create", "create a service in OpsLevel")
	serviceCreateCmdNameArg = serviceCreateCmd.Arg("name", "service name").Required().String()
	// Get Service
	serviceGetCmd         = serviceCmd.Command("get", "get a service by name / alias")
	serviceGetCmdAliasArg = serviceGetCmd.Arg("alias", "service alias").Required().String()
	// Tag Service
	serviceTagCmd            = serviceCmd.Command("tag", "manipulate service tags")
	serviceTagAddCmd         = serviceTagCmd.Command("add", "Add service tag")
	serviceTagAddCmdAliasArg = serviceTagAddCmd.Arg("alias", "service alias").Required().String()
	serviceTagAddCmdTagsArg  = serviceTagAddCmd.Arg("k=v", "key=value tag").Required().StringMap()
	// Delete Service
	serviceDeleteCmd          = serviceCmd.Command("delete", "delete a service in OpsLevel")
	serviceDeleteCmdIdFlag    = serviceDeleteCmd.Flag("id", "Id of service to delete").String()
	serviceDeleteCmdAliasFlag = serviceDeleteCmd.Flag("alias", "Alias of service to delete").String()

	// Teams
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
	case serviceCreateCmd.FullCommand():
		svc := opslevel.Service{
			Name: *serviceCreateCmdNameArg,
		}
		resp, err := client.CreateService(context.Background(), svc)
		if err != nil {
			log.Fatalln(err)
		}

		pretty.Println(resp)

	case serviceDeleteCmd.FullCommand():
		var resp *opslevel.DeleteServiceResponse
		var err error

		if serviceDeleteCmdIdFlag != nil && *serviceDeleteCmdIdFlag != "" {
			resp, err = client.DeleteServiceById(context.Background(), *serviceDeleteCmdIdFlag)
		} else if serviceDeleteCmdAliasFlag != nil && *serviceDeleteCmdAliasFlag != "" {
			resp, err = client.DeleteServiceByAlias(context.Background(), *serviceDeleteCmdAliasFlag)
		} else {
			app.FatalUsage("--alias or --id flag must be provided")
		}
		if err != nil {
			log.Fatalln(err)
		}

		pretty.Println(resp)

	case serviceTagAddCmd.FullCommand():
		for key, value := range *serviceTagAddCmdTagsArg {
			tag, err := client.CreateTag(context.Background(), key, value, *serviceTagAddCmdAliasArg, "Service")
			if err != nil {
				log.Fatalln(err)
			}
			pretty.Println(tag)
		}

	case serviceGetCmd.FullCommand():
		svc, err := client.GetServiceByAlias(context.Background(), *serviceGetCmdAliasArg)
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
