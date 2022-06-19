/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"
)

var wg sync.WaitGroup

func DownloadVideo(client youtube.Client, playlist *youtube.Playlist, i int, v *youtube.PlaylistEntry) {
	entry := playlist.Videos[i]
	video, err := client.VideoFromPlaylistEntry(entry)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Downloading %s by '%s'!\n", video.Title, video.Author)

	formats := video.Formats.Quality("hd720").WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}

	file, err := os.Create(strconv.Itoa(i) + " " + v.Title + ".mp4")

	if err != nil {
		panic(err)
	}

	defer wg.Done()
	defer file.Close()
	_, err = io.Copy(file, stream)

	if err != nil {
		panic(err)
	}
}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads the specified youtube playlist link",
	Long:  "Provided the link, this command downloads the whole youtube playlist for you!",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide the URL of the playlist you want to download")
			return
		}
		client := youtube.Client{}

		playlist, err := client.GetPlaylist(args[0])
		if err != nil {
			panic(err)
		}

		/* ----- Enumerating playlist videos ----- */
		header := fmt.Sprintf("Playlist %s by %s", playlist.Title, playlist.Author)
		println(header)
		println(strings.Repeat("=", len(header)) + "\n")

		if err := os.Mkdir(header, os.ModePerm); err != nil {
			panic(err)
		}

		os.Chdir(header)

		for i, v := range playlist.Videos {
			wg.Add(1)
			go DownloadVideo(client, playlist, i, v)
		}

		wg.Wait()

		fmt.Println("Download Complete!")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
