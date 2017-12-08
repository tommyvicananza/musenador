package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	discogs "github.com/irlndts/go-discogs"
	"github.com/robdimsdale/wl"
	"github.com/robdimsdale/wl/logger"
	"github.com/robdimsdale/wl/oauth"
)

type song struct {
	artist string
	title  string
}

func main() {
	accessToken := flag.String("access-token", os.Getenv("WL_ACCESS_TOKEN"), "access token of your WunderList Account.")
	clientID := flag.String("client-id", os.Getenv("WL_CLIENT_ID"), "client ID of your WunderList account.")
	discogsToken := flag.String("discogs-token", os.Getenv("DGS_TOKEN"), "Discogs token.")
	flag.Parse()
	if *accessToken == "" || *clientID == "" || *discogsToken == "" {
		exit("Missing arguments", fmt.Errorf("error"))
	}
	var titles []string
	client := oauth.NewClient(
		*accessToken,
		*clientID,
		wl.APIURL,
		logger.NewLogger(logger.INFO),
	)
	c := discogs.NewClient("MolinasTest/0.1", *discogsToken)

	// Ignore error
	lists, err := client.Lists()
	if err != nil {
		exit("Error getting list", err)
	}
	for _, lMusic := range lists {
		if lMusic.Title == "Música" {
			tasks, err := client.TasksForListID(lMusic.ID)
			if err != nil {
				exit("Error getting tasks", err)
			}
			for _, task := range tasks {
				titles = append(titles, task.Title)
			}
		}
	}
	sort.Strings(titles)
	for _, title := range titles {
		//if title == "exium - subtoned" {
		if strings.Contains(title, " - ") {
			result := strings.Split(title, " - ")
			t := &song{
				artist: result[0],
				title:  result[1],
			}
			time.Sleep(2 * time.Second)
			request := &discogs.SearchRequest{Q: title, Artist: t.artist, Page: 0, Per_page: 1}
			s, _, err := c.Search.Search(request)
			if err != nil {
				fmt.Println(err)
			}
			for _, n := range s.Results {
				fmt.Println("-------")
				fmt.Printf("Canción: %s\nDisco: %s\n", title, n.Title)
				for _, style := range n.Style {
					fmt.Printf("Estilo: %s\n", style)
				}
				fmt.Println("-------")
			}
		}
	}
}

func exit(msg string, e error) {
	fmt.Printf("%s: %s\n", msg, e)
	os.Exit(1)
}
