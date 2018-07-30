// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package start

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/medtune/beta-platform/cmd/root"
	"github.com/medtune/beta-platform/pkg/config"
	"github.com/medtune/beta-platform/pkg/initpkg"
	"github.com/medtune/beta-platform/pkg/store"
	"github.com/medtune/beta-platform/server"
	"github.com/spf13/cobra"
)

var (
	static      string
	configFile  string
	port        int
	ginMode     int
	syncdb      bool
	wait        bool
	maxattempt  int
	timestamp   int
	createUsers bool
)

func init() {
	startCmd.Flags().StringVarP(&configFile, "file", "f", "./config.yml", "Configuration file name")
	startCmd.Flags().StringVarP(&static, "static", "s", "./static", "Static files directory")

	startCmd.Flags().IntVarP(&port, "port", "p", 8005, "port")
	startCmd.Flags().IntVarP(&ginMode, "gin-mode", "g", 0, "Gin server mode [0 OR 1]")
	startCmd.Flags().BoolVarP(&syncdb, "syncdb", "x", true, "Sync database before start")
	startCmd.Flags().BoolVarP(&createUsers, "create-users", "y", true, "Create default users before start")

	startCmd.Flags().BoolVarP(&wait, "wait", "w", false, "Wait all services to go up")
	startCmd.Flags().IntVarP(&maxattempt, "wait-attempts", "c", 30, "Wait max attempts")
	startCmd.Flags().IntVarP(&timestamp, "wait-timestamp", "t", 1, "Wait timestamp")

	root.Cmd.AddCommand(startCmd)
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"run", "run-server"},
	Short:   "Run Medtune beta server",
	Long:    `Run Medtune beta server`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func runServer() {
	if ginMode == 0 {
		gin.SetMode(gin.DebugMode)
	} else if ginMode == 1 {
		gin.SetMode(gin.ReleaseMode)
	}

	var configuration *config.StartupConfig

	if cfg, err := config.LoadConfigFromPath(configFile); err != nil {
		log.Panic(err)
	} else {
		configuration = cfg
	}

	if err := initpkg.InitFromConfig(configuration); err != nil {
		log.Panic(err)
	}

	if syncdb {
		var err error

		err = store.Agent.Sync()
		fmt.Printf("trying to connect database %s\n", configuration.Database.Prod)

		if err != nil && wait {
			attempt := 0
			for err != nil && attempt < maxattempt {
				time.Sleep(time.Duration(timestamp) * time.Second)
				fmt.Printf("waiting for database response (host: %s)\n", configuration.Database.Creds.Host)
				err = store.Agent.Sync()
				attempt++
			}
		}
		if err != nil {
			panic(err)
		}
		fmt.Println("connected to database...")
	}

	if createUsers && configuration.Create != nil {
		var err error
		for _, user := range configuration.Create.Users {
			err = store.Agent.CreateUser(user.Email, user.Username, user.Password)
			if err != nil {
				fmt.Printf("failed to create user: %s\n    error: %v\n", user.Username, err)
				continue
			}
			fmt.Printf("created user %s %s %s\n", user.Email, user.Username, user.Password)
		}
	}

	fmt.Println("starting server...")

	Server := server.New(
		static,
		port,
	)

	Server.Run()
}
