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

func ytGetResults(args string, client *youtube.Service) (*youtube.SearchListResponse, error) {
	searchClient := youtube.NewSearchService(client)

	search := searchClient.List([]string{"snippet"}).
		Q(args).
		MaxResults(1)

	return search.Do()
}

func formatResult(res *youtube.SearchResult) string {
	snippet := res.Snippet
	channel := snippet.ChannelTitle
	title := html.UnescapeString(snippet.Title)

	return fmt.Sprintf("{red}YouTube video by %s:{clear} %s", channel, title)
}

func ytSearch(args string, client *youtube.Service) (string, error) {
	results, err := ytGetResults(args, client)

	if err != nil {
		return "", err
	}

	if len(results.Items) == 0 {
		return fmt.Sprintf("No results found for %s", args), nil
	}

	res := results.Items[0]

	return fmt.Sprintf("%s - https://youtu.be/%s", formatResult(res), res.Id.VideoId), nil
}

func ytTitle(msg string, client *youtube.Service) (string, error) {
	urls := ytURLRe.FindAllString(msg, -1)

	outList := []string{}

	for _, url := range urls {
		results, err := ytGetResults(url, client)

		if err != nil {
			return "", err
		}

		if len(results.Items) == 0 {
			continue
		}

		outList = append(outList, formatResult(results.Items[0]))
	}

	return strings.Join(outList, "\n"), nil
}
