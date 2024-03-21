package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/kkdai/youtube/v2"
	"github.com/vague2k/smv/utils"
)

var labelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("245"))

type MetadataItem struct {
	Label string
	Value string
}

func VideoMetadata(video *youtube.Video) {
	metadata := []MetadataItem{
		{Label: "Video ID", Value: video.ID},
		{Label: "Title", Value: video.Title},
		{Label: "Posted by", Value: video.Author},
		{Label: "Views", Value: fmt.Sprintf("%d", video.Views)},
		{Label: "Duration", Value: utils.FormatTimeDuration(video.Duration)},
	}

	fmt.Println() // adds a line of padding on the top of the metadata info block
	for _, item := range metadata {
		fmt.Println(fmt.Sprintf("     %s %s", labelStyle.Render(item.Label+":"), item.Value))
	}
}
