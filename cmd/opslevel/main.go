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
	servicesAlias = serviceGetCmd.Arg("alias", "service alias").Required().String()

	version = "dev"
)

func main() {
	app.Author("Matt Morrison (sl1pm4t)")
	app.Version("opslevel version " + version)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	var cmd = kingpin.MustParse(app.Parse(os.Args[1:]))

	var authToken = os.Getenv("OPSLEVEL_TOKEN")
	client := opslevel.NewClient(authToken)

	switch cmd {
	case serviceGetCmd.FullCommand():
		alias := "mgob"
		svc, err := client.GetService(context.Background(), alias)
		if err != nil {
			log.Fatalln(err)
		}

		pretty.Println(svc)
	}
}
