package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"
	"github.com/vague2k/smv/utils"
)

type Video struct {
	Metadata *youtube.Video
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "smv",
	Run: func(cmd *cobra.Command, args []string) {
		client := youtube.Client{}
		url := args[0]

		video, err := client.GetVideo(url)
		if err != nil {
			log.Fatalf("\nCould not get video from url/id: %s\nThe following error was given:\n%s", url, err)
		}

		fmt.Println(video.ID)
		fmt.Println(video.Author)
		fmt.Println(video.Title)
		fmt.Println(video.Views)
		fmt.Println(video.Duration)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.splitter.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	cobra.OnInitialize(utils.CheckForDependencies)
}
