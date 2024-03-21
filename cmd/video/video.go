package video

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/kkdai/youtube/v2"
	"github.com/vague2k/smv/utils"
)

type Video struct {
	Metadata    *youtube.Video
	MP4FilePath string
	MP3FilePath string
	cmd         *exec.Cmd
}

func (v *Video) CleanupTempFile() {
	tempFiles := []string{v.MP4FilePath, v.MP3FilePath}

	for _, file := range tempFiles {
		_, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				// if file doesn't exist, no need to clean up
				continue
			}
			fmt.Printf("\nCould not stat temp file for cleanup.\nThe Following error was given:\n%v\n", err)
		}

		if err := os.Remove(file); err != nil {
			fmt.Printf("\nCould not clean up temp files.\nThe Following error was given:\n%v\n", err)
		}
	}
}

// Returns a temp file to use for proccessing. You can access the temp file via the Video struct
func (v *Video) DownloadMP4AsTemp(videoUrl string) {
	client := youtube.Client{}

	video, err := client.GetVideo(videoUrl)
	if err != nil {
		log.Fatalf("\nCould not get video from url/id: %s\nThe following error was given:\n%s", videoUrl, err)
	}
	v.Metadata = video

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		log.Fatalf("\nVideo metadata was fetched, but could not fetch byte stream.\nThe following error was given: %s", err)
	}

	defer stream.Close()

	file, err := os.CreateTemp("", fmt.Sprintf("%s", "smv_temp"))
	if err != nil {
		log.Fatalf("\nCould not create the video temp file.\nThe following error was given: %s", err)
	}
	v.MP4FilePath = file.Name()
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		log.Fatalf("\nFetching neccessary video data was successful, but could not copy byte stream to temp file.\nThe following error was given: %s", err)
	}
}

// Converts a mp4 file to mp3 using ffmpeg
func (v *Video) CovertToMp3() {
	splitFilePath := strings.Split(v.MP4FilePath, ".")
	v.MP3FilePath = splitFilePath[0] + ".mp3"

	v.cmd = exec.Command("ffmpeg", "-i", v.MP4FilePath, v.MP3FilePath)
	if err := v.cmd.Run(); err != nil {
		log.Fatalf("\nFFmpeg could convert the temp mp4 file to mp3.\nThe following error was given:\n%s", err)
	}
}

// HACK: Surely, there's a better way of doing this.
func (v *Video) Split(threshold, duration string) {
	v.cmd = exec.Command("ffmpeg", "-i", v.MP3FilePath, "-af", "silencedetect=n="+threshold+"dB:d="+duration, "-f", "null", "-")
	output, err := v.cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("\nThere was a problem detecting silent segments from the given video\nThe following error was given:\n%v", err)
	}

	str := string(output)
	pattern := `\bsilence_(start|end):\s*([\d.]+)`

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(str, -1)

	var startTimes []float64
	var endTimes []float64

	for _, match := range matches {
		event := match[1]      // "start" or "end"
		secondsStr := match[2] // timestamp of detected silence in seconds

		seconds, err := strconv.ParseFloat(secondsStr, 64)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if event == "start" {
			startTimes = append(startTimes, seconds)
		}
		if event == "end" {
			endTimes = append(endTimes, seconds)
		}
	}

	for i := 0; i < len(startTimes); i++ {

		start := startTimes[i]

		fileName := filepath.Join(utils.UserHomeDir(), "Downloads", fmt.Sprintf("smv_output_%d.mp3", i+1))
		if i != 0 {
			v.cmd = exec.Command("ffmpeg", "-i", v.MP3FilePath,
				"-ss", utils.FormatFFmpegSilentDuration(endTimes[i-1]-1),
				"-to", utils.FormatFFmpegSilentDuration(start+1),
				"-c", "copy",
				fmt.Sprint(fileName),
			)
			err := v.cmd.Run()
			if err != nil {
				fmt.Printf("\nThere was a problem splitting your videos audios.\nThe following error was given:\n%v", err)
			}
		}

		if i == 0 {
			// This is edge case if start of a segment matches end of a segment on the first iteration
			if utils.FormatFFmpegSilentDuration(start) == "00:00:00" {
				continue
			}

			v.cmd = exec.Command("ffmpeg", "-i", v.MP3FilePath,
				"-ss", "00:00:00",
				"-to", utils.FormatFFmpegSilentDuration(start),
				"-c", "copy",
				fmt.Sprint(fileName),
			)
			err := v.cmd.Run()
			if err != nil {
				fmt.Printf("\nThere was a problem splitting your videos audios.\nThe following error was given:\n%v", err)
			}
		}
	}
}
