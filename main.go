package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/drgarcia1986/gonews/client"
	"github.com/drgarcia1986/gonews/gui"
	"github.com/drgarcia1986/gonews/providers"
	"github.com/drgarcia1986/gonews/providers/hackernews"
	"github.com/drgarcia1986/gonews/providers/reddit"
	"github.com/drgarcia1986/gonews/story"
	"github.com/drgarcia1986/gonews/utils"
)

var (
	limit            int
	storyType        int
	flagStoryType    string
	flagProviderType string
)

func init() {
	flag.IntVar(&limit, "limit", 10, "Number of Stories to get")
	flag.StringVar(&flagStoryType, "type", "top", "Stories Type ('new' or 'top')")
	flag.StringVar(
		&flagProviderType, "provider", "hackernews",
		"Stories Provider (hackernews, reddit, reddit-<subreddit>)")

	flag.Parse()

	if len(flag.Args()) > 0 && flag.Args()[0] == "version" {
		fmt.Printf("GoNews %s\n", utils.Version)
		os.Exit(0)
	}

	if flagStoryType == "new" {
		storyType = story.NewStories
	} else {
		storyType = story.TopStories
	}
}

func getProvider(providerType string) providers.Provider {
	switch providerType {
	case "hackernews":
		return hackernews.New()
	case "reddit":
		return reddit.New()
	default:
		providerArgs := strings.Split(providerType, "-")
		if len(providerArgs) != 2 || providerArgs[0] != "reddit" {
			return nil
		}
		return reddit.NewSubReddit(providerArgs[1])
	}
}

func main() {
	p := getProvider(flagProviderType)
	if p == nil {
		fmt.Println("Invalid Provider ", flagProviderType)
		os.Exit(0)
	}

	fmt.Printf("Getting %s stories\n", p.Name())
	c := client.NewWithPB(p)
	stories, err := c.GetStories(storyType, limit)
	if err != nil {
		panic(err)
	}

	g := gui.New(stories, p.Name())
	if err = g.Run(); err != nil {
		panic(err)
	}
}
