package utils

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

func UserHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("\nUserHomeDir: How'd you fucked this up\n")
	}

	return dir
}

func CheckForFFmpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	if err != nil {
		return errors.New("The ffmpeg package is not installed. Install it to use this program!")
	}
	return nil
}

func CheckForNetworkConn() error {
	_, err := net.Dial("tcp", "www.google.com:80")
	if err != nil {
		return errors.New("To use the downloading capabilities of this tool, please connect to a network!")
	}
	return nil
}

func CheckForDependencies() {
	if err := CheckForNetworkConn(); err != nil {
		fmt.Printf("\n%v\n", err)
		os.Exit(126)
	}

	if err := CheckForFFmpeg(); err != nil {
		fmt.Printf("\n%v\n", err)
		os.Exit(126)
	}
}

func FormatFFmpegSilentDuration(seconds float64) string {
	hours := int(seconds / 3600)
	minutes := int(seconds / 60)
	secondsRem := int(seconds) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secondsRem)
}

func FormatTimeDuration(d time.Duration) string {
	hours := d / time.Hour
	minutes := (d % time.Hour) / time.Minute
	seconds := (d % time.Minute) / time.Second

	timestamp := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	return timestamp
}
