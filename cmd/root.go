// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/BrianHicks/habitifocus/habitica"
	"github.com/BrianHicks/habitifocus/omnifocus"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "habitifocus",
	Short: "sync OmniFocus tasks to Habitica",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := omnifocus.GetTasks()
		if err != nil {
			log.Fatal(err)
		}
		logrus.WithField("tasks", len(tasks)).Info("got OmniFocus task list")

		client := habitica.Client{
			UserID: viper.GetString("userid"),
			APIKey: viper.GetString("apikey"),
		}

		remoteTasks, err := client.List()
		if err != nil {
			log.Fatal(err)
		}
		logrus.WithField("tasks", len(remoteTasks)).Info("got Habitica task list")

		touched := map[string]struct{}{}
		for id, task := range tasks {
			remoteTask, exists := remoteTasks[id]

			if !task.Done && !exists {
				err := client.Create(&habitica.HabiticaTODO{
					Alias:     id,
					Text:      task.Name,
					Type:      "todo",
					Completed: task.Done,
				})
				if err != nil {
					logrus.WithError(err).WithField("task", task).Fatal("could not create task")
				} else {
					logrus.WithField("task", task).Info("created")
				}
				touched[id] = struct{}{}
			}

			if task.Done && exists && !remoteTask.Completed {
				err := client.Complete(remoteTask)
				if err != nil {
					logrus.WithError(err).WithField("task", task).Fatal("could not complete task")
				} else {
					logrus.WithField("task", task).Info("completed")
				}
				touched[id] = struct{}{}
			}
		}

		for id, task := range remoteTasks {
			_, exists := tasks[id]

			if !exists {
				err := client.Delete(task)
				if err != nil {
					logrus.WithError(err).WithField("task", task).Fatal("could not delete task")
				} else {
					logrus.WithField("task", task).Info("deleted")
				}
			}
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.habitifocus.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().String("userid", "", "user ID for Habitica")
	RootCmd.Flags().String("apikey", "", "API key for Habitica")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.BindPFlags(RootCmd.Flags()); err != nil {
		log.Fatal(err)
	}

	viper.SetConfigName(".habitifocus") // name of config file (without extension)
	viper.AddConfigPath("$HOME")        // adding home directory as first search path
	viper.AutomaticEnv()                // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.WithField("file", viper.ConfigFileUsed()).Info("using values from config file")
	}
}
