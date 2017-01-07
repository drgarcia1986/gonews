package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/drgarcia1986/gonews/client"
	"github.com/drgarcia1986/gonews/gui"
	"github.com/drgarcia1986/gonews/providers"
	"github.com/drgarcia1986/gonews/providers/hackernews"
	"github.com/drgarcia1986/gonews/story"
)

const version = "0.0.1"

var (
	limit            int
	storyType        int
	flagStoryType    string
	flagProviderType string = "hn"
)

func init() {
	flag.IntVar(&limit, "limit", 10, "Number of Stories to get")
	flag.StringVar(&flagStoryType, "type", "top", "Stories Type ('new' or 'top')")

	flag.Parse()

	if len(flag.Args()) > 0 && flag.Args()[0] == "version" {
		fmt.Printf("GoNews %s\n", version)
		os.Exit(0)
	}

	if flagStoryType == "new" {
		storyType = story.NewStories
	} else {
		storyType = story.TopStories
	}
}

func getProvider(providerType string) providers.Provider {
	return hackernews.New()
}

func main() {
	fmt.Println("Getting HackerNews stories")

	p := getProvider(flagProviderType)
	c := client.NewWithPB(p)
	stories, err := c.GetStories(storyType, limit)
	if err != nil {
		panic(err)
	}

	g := gui.New(stories)
	if err = g.Run(); err != nil {
		panic(err)
	}
}
