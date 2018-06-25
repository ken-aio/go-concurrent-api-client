package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/resty.v1"
	"gopkg.in/workanator/go-floc.v2"
	"gopkg.in/workanator/go-floc.v2/run"
)

const apiBase = "http://localhost:9998/api"

type (
	// Episode episode
	Episode struct {
		ID      int    `json:"id"`
		TitleID int    `json:"title_id"`
		Name    string `json:"name"`
		Desc    string `json:"desc"`
	}
	// Episodes episodes
	Episodes struct {
		Episodes []*Episode `json:"episodes"`
		TitleID  int        // custom value
	}
	// Title title
	Title struct {
		ID       int        `json:"id"`
		Name     string     `json:"name"`
		Desc     string     `json:"desc"`
		Episodes []*Episode `json:"episodes"`
	}
	// Titles titles
	Titles struct {
		Titles []*Title `json:"titles"`
	}
	// APIResult api result
	APIResult struct {
		*Titles
	}
)

const (
	keyTitles    = "titles"
	keyEpisodes  = "episodes"
	keyAPIResult = "result"
)

func main() {
	initApp()

	job := run.Sequence(
		createTitlesFunc(),
		run.Parallel(
			createTitleDetailFunc(),
			createEpisodesFunc(),
		),
		createDoMergeFunc(),
	)
	ctx := floc.NewContext()
	ctrl := floc.NewControl(ctx)
	ctx.AddValue(keyAPIResult, &APIResult{})

	_, _, err := floc.RunWith(ctx, ctrl, job)
	if err != nil {
		panic(err)
	}

	result := ctx.Value(keyAPIResult).(*APIResult)
	resp, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	log.Printf("resp = %+v\n", string(resp))
}

// createTitlesFunc is return func for request titles api
func createTitlesFunc() floc.Job {
	return func(ctx floc.Context, ctrl floc.Control) error {
		log.Println("doTitles")
		time.Sleep(500 * time.Millisecond)
		titlesResp, err := reqTitles()
		if err != nil {
			return err
		}
		ctx.AddValue(keyTitles, titlesResp)
		episodes := make([]*Episodes, len(titlesResp.Titles))
		for i, t := range titlesResp.Titles {
			episodes[i] = &Episodes{TitleID: t.ID}
		}
		ctx.AddValue(keyEpisodes, episodes)
		return nil
	}
}

// createTitleDetailFunc is return func for request title api
func createTitleDetailFunc() floc.Job {
	return func(ctx floc.Context, ctrl floc.Control) error {
		titles := ctx.Value(keyTitles).(*Titles)
		detailFuncs := make([]floc.Job, len(titles.Titles))
		for i, t := range titles.Titles {
			index, id := i, t.ID // *here[Point] same index and id if not set here
			detailFuncs[index] = func(ctx floc.Context, ctrl floc.Control) error {
				log.Println("doTitle detail", index, id)
				time.Sleep(1500 * time.Millisecond)
				var err error
				titles.Titles[index], err = reqTitle(strconv.Itoa(id))
				if err != nil {
					return err
				}
				return nil
			}
		}
		job := run.Parallel(detailFuncs...)
		_, _, err := floc.Run(job)
		return err
	}
}

// createEpisodesFunc is return func for request episodes api
func createEpisodesFunc() floc.Job {
	return func(ctx floc.Context, ctrl floc.Control) error {
		log.Println("doEpisodes")
		episodes := ctx.Value(keyEpisodes).([]*Episodes)
		detailFuncs := make([]floc.Job, len(episodes))
		for i, e := range episodes {
			index, titleID := i, e.TitleID // *here[Point] same index and id if not set here
			detailFuncs[index] = func(ctx floc.Context, ctrl floc.Control) error {
				log.Println("doEpisodes", index, titleID)
				time.Sleep(1500 * time.Millisecond)
				resp, err := reqEpisodes(strconv.Itoa(titleID))
				if err != nil {
					return err
				}
				episodes[index].Episodes = resp.Episodes
				err = runEpisodeDetails(titleID, episodes[index])
				return err
			}
		}
		job := run.Parallel(detailFuncs...)
		_, _, err := floc.Run(job)
		return err
	}
}

func createDoMergeFunc() floc.Job {
	return func(ctx floc.Context, ctrl floc.Control) error {
		titles := ctx.Value(keyTitles).(*Titles)
		episodes := ctx.Value(keyEpisodes).([]*Episodes)
		result := ctx.Value(keyAPIResult).(*APIResult)

		result.Titles = titles
		for i, t := range result.Titles.Titles {
			for _, e := range episodes {
				if t.ID == e.TitleID {
					result.Titles.Titles[i].Episodes = e.Episodes
				}
			}
		}
		fmt.Printf("result = %+v\n", result)
		return nil
	}
}

func runEpisodeDetails(titleID int, es *Episodes) error {
	detailFuncs := make([]floc.Job, len(es.Episodes))
	for i, e := range es.Episodes {
		index, episodeID := i, e.ID // *here[Point] same index and id if not set here
		detailFuncs[index] = func(ctx floc.Context, ctrl floc.Control) error {
			log.Println("doEpisodeDetail ", index, titleID, episodeID)
			time.Sleep(1500 * time.Millisecond)
			var err error
			es.Episodes[index], err = reqEpisode(strconv.Itoa(titleID), strconv.Itoa(episodeID))
			return err
		}
	}
	job := run.Parallel(detailFuncs...)
	_, _, err := floc.Run(job)
	return err
}

func initApp() {
	resty.SetDebug(false)
	resty.SetLogger(os.Stdout)
	resty.SetTimeout(1 * time.Second)
}

func doGet(url string, result interface{}) error {
	_, err := resty.R().SetResult(result).Get(url)
	return err
}

func reqTitles() (*Titles, error) {
	url := apiBase + "/titles"
	titles := &Titles{}
	err := doGet(url, titles)
	return titles, err
}

func reqTitle(id string) (*Title, error) {
	url := apiBase + "/titles/" + id
	title := &Title{}
	err := doGet(url, title)
	return title, err
}

func reqEpisodes(titleID string) (*Episodes, error) {
	url := apiBase + "/titles/" + titleID + "/episodes"
	episodes := &Episodes{}
	err := doGet(url, episodes)
	return episodes, err
}

func reqEpisode(titleID, episodeID string) (*Episode, error) {
	url := apiBase + "/titles/" + titleID + "/episodes/" + episodeID
	episode := &Episode{}
	err := doGet(url, episode)
	return episode, err
}
