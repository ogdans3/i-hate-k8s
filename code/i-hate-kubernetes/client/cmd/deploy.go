package cmd

import (
	"os"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/yaml"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runDeploy,
}

var file string

func init() {
	deployCmd.Flags().StringVarP(&file, "file", "f", "", "Specify the project file to use")
	rootCmd.AddCommand(deployCmd)
}

func runDeploy(cmd *cobra.Command, args []string) {
	pwd, _ := os.Getwd()
	console.Log(pwd)
	console.Log(pwd + "/../" + file)
	project := yaml.ReadFile(pwd + "/../" + file)
	console.Log(project)

	c := client.CreateClient()
	c.AddProject(project)
	c.Loop()
}
