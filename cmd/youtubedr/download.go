package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Downloads a video from youtube",
	Example: `youtubedr -o "Campaign Diary".mp4 https://www.youtube.com/watch\?v\=XbNghLqsVwU`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exitOnError(download(args[0]))
	},
}

var (
	ffmpegCheck        error
	outputFile         string
	outputDir          string
	downloadTranscript bool
)

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&outputFile, "filename", "o", "", "The output file, the default is genated by the video title.")
	downloadCmd.Flags().StringVarP(&outputDir, "directory", "d", ".", "The output directory.")
	downloadCmd.Flags().BoolVarP(&downloadTranscript, "transcript", "t", false, "Download the transcript only")
	addQualityFlag(downloadCmd.Flags())
	addMimeTypeFlag(downloadCmd.Flags())
}

func download(id string) error {
	video, format, err := getVideoWithFormat(id)
	if err != nil {
		return err
	}

	if strings.HasPrefix(outputQuality, "hd") {
		if err := checkFFMPEG(); err != nil {
			return err
		}
		return downloader.DownloadComposite(context.Background(), outputFile, video, outputQuality, mimetype)
	}

	if downloadTranscript {
		log.Println("download transcript to directory", outputDir)
		return downloader.DownloadTranscript(context.Background(), video, outputFile, "")
	}

	log.Println("download to directory", outputDir)
	return downloader.Download(context.Background(), video, format, outputFile)
}

func checkFFMPEG() error {
	fmt.Println("check ffmpeg is installed....")
	if err := exec.Command("ffmpeg", "-version").Run(); err != nil {
		ffmpegCheck = fmt.Errorf("please check ffmpegCheck is installed correctly")
	}

	return ffmpegCheck
}
