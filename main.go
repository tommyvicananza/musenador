package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	discogs "github.com/irlndts/go-discogs"
	"github.com/robdimsdale/wl"
	"github.com/robdimsdale/wl/logger"
	"github.com/robdimsdale/wl/oauth"
)

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
		fmt.Println(title)
	}

	c := discogs.NewClient()
	c.UserAgent("MolinasTest/0.1 +http://discogs.vicananza.com")
	c.Token(*discogsToken)
	searchRequest := c.Search
	request := &discogs.SearchRequest{Artist: "reggaenauts", Release_title: "river rock", Page: 0, Per_page: 1}
	s, _, err := searchRequest.Search(request)
	if err != nil {
		fmt.Println(err)
	}
	for i, n := range s.Results {
		fmt.Println(i, n.Title)
	}
}

func exit(msg string, e error) {
	fmt.Printf("%s: %s\n", msg, e)
	os.Exit(1)
}
