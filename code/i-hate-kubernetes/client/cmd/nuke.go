package cmd

import (
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runNuke,
}

func init() {
	rootCmd.AddCommand(nukeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runNuke(cmd *cobra.Command, args []string) {
	c := client.CreateClient()
	c.Nuke()
}
