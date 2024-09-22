package cmd

import (
	"fmt"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/container-interface/docker"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/yaml"
	"os"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runStop(cmd *cobra.Command, args []string) {
	pwd, _ := os.Getwd()
	project := yaml.ReadFile(pwd + "/examples/hello-world.yml")

	fmt.Printf("%v\n", project.Project)
	fmt.Printf("%v\n", project)

	docker.ListAllContainers()
	docker.StopProjectContainers(&project)
	fmt.Printf("\n")
	docker.ListAllContainers()
}
