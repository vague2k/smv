package cmd

import (
	"fmt"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/vague2k/smv/cmd/ui"
	"github.com/vague2k/smv/cmd/video"
	"github.com/vague2k/smv/utils"
)

var (
	finishedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	vd            = &video.Video{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "smv",
	Args:  cobra.MinimumNArgs(1),
	Short: "A youtube mp4 to mp3 downloader audio splitter",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var spinner *tea.Program
		wg := sync.WaitGroup{}

		threshold := cmd.Flag("threshold").Value.String()
		silenceDuration := cmd.Flag("duration").Value.String()

		// Init wait groups to synchronize spinners with main function
		wg.Add(3)
		go func() {
			defer wg.Done()
			spinner = tea.NewProgram(ui.SpinnerModel("Downloading youtube video for processing..."))
			if _, err := spinner.Run(); err != nil {
				cobra.CheckErr(err)
			}
		}()

		vd.DownloadMP4AsTemp(args[0])
		ui.VideoMetadata(vd.Metadata)
		spinner.ReleaseTerminal()

		go func() {
			defer wg.Done()
			spinner = tea.NewProgram(ui.SpinnerModel("Converting mp4 to mp3 using ffmpeg..."))
			if _, err := spinner.Run(); err != nil {
				cobra.CheckErr(err)
			}
		}()

		vd.CovertToMp3()
		spinner.ReleaseTerminal()

		// Parse out the silent parts and split the mp3 based on those segments
		go func() {
			defer wg.Done()
			spinner = tea.NewProgram(ui.SpinnerModel(fmt.Sprintf("Splitting audios based on silence. threshold: %sdb. duration: %ss", threshold, silenceDuration)))
			if _, err := spinner.Run(); err != nil {
				cobra.CheckErr(err)
			}
		}()

		vd.Split(threshold, silenceDuration)
		spinner.ReleaseTerminal()

		done := fmt.Sprintf("\n  %s", finishedStyle.Render("Finished splitting audios, check your downloads folder!"))
		fmt.Println(done)
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
	cobra.OnFinalize(vd.CleanupTempFile)

	rootCmd.Flags().Int8P("threshold", "t", -30, "Minimum sound level (decibel) to trigger a silence detection.")
	rootCmd.Flags().Float64P("duration", "s", 0.6, "Duration of silence to detect if threshold was met.")
}
