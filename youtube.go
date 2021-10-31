package main

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"google.golang.org/api/youtube/v3"
)

var (
	ytURLRe = regexp.MustCompile(ytURLRegex)
)

func ytSearch(args string, client *youtube.Service) (string, error) {
	searchClient := youtube.NewSearchService(client)

	search := searchClient.List([]string{"snippet"}).
		Q(args).
		MaxResults(1)

	results, err := search.Do()

	if err != nil {
		return "", err
	}

	if len(results.Items) == 0 {
		return fmt.Sprintf("No results found for %s", args), nil
	}

	result := results.Items[0]
	snippet := result.Snippet
	channel := snippet.ChannelTitle
	title := html.UnescapeString(snippet.Title)

	out := fmt.Sprintf("{red}YouTube video by %s:{clear} %s", channel, title)

	return out, nil
}

func ytTitle(msg string, client *youtube.Service) (string, error) {
	urls := ytURLRe.FindAllString(msg, -1)

	outList := []string{}

	for _, url := range urls {
		res, err := ytSearch(url, client)

		if err != nil {
			return "", err
		}

		outList = append(outList, res)
	}

	return strings.Join(outList, "\n"), nil
}
