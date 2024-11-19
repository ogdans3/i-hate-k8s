package cmd

import (
	"os/exec"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
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

func runNuke(c *cobra.Command, args []string) {

	stop()
	rm()
	rmVols()
	prune()

	//c := client.CreateClient()
	//c.Nuke()
}

func stop() {
	cmd := exec.Command("bash", "-c", "docker stop $(docker ps -q)")
	if err := cmd.Run(); err != nil {
		console.InfoLog.Fatal("Error stopping containers: %v", err)
	} else {
		console.InfoLog.Info("All running containers stopped successfully.")
	}
}

func rm() {
	cmd := exec.Command("bash", "-c", "docker rm -f $(docker ps -a -q)")
	if err := cmd.Run(); err != nil {
		console.InfoLog.Fatal("Error removing containers: %v", err)
	} else {
		console.InfoLog.Info("All stopped containers removed successfully.")
	}
}

func prune() {
	cmd := exec.Command("bash", "-c", "docker system prune -f --volumes")

	if err := cmd.Run(); err != nil {
		console.InfoLog.Fatal("Error running prune: %v", err)
	}

	console.InfoLog.Info("Docker system prune completed successfully.")
}

func rmVols() {
	removeVolumesCmd := exec.Command("bash", "-c", "docker volume rm $(docker volume ls -q)")
	if err := removeVolumesCmd.Run(); err != nil {
		console.InfoLog.Fatal("Error deleting volumes: %v", err)
	} else {
		console.InfoLog.Info("All Docker volumes deleted successfully.")
	}
}
