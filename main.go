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
		Titles []Title `json:"titles"`
	}
	// APIResult api result
	APIResult struct {
		*Titles
	}
)

func main() {
	initApp()

	const keyTitles = "titles"
	const keyEpisodes = "episodes"
	const keyAPIResult = "result"
	doTitles := func(ctx floc.Context, ctrl floc.Control) error {
		log.Println("doTitles")
		titlesResp := reqTitles()
		ctx.AddValue(keyTitles, titlesResp)
		episodes := make([]*Episodes, len(titlesResp.Titles))
		for i, t := range titlesResp.Titles {
			episodes[i] = &Episodes{TitleID: t.ID}
		}
		ctx.AddValue(keyEpisodes, episodes)
		return nil
	}
	doTitle := func(ctx floc.Context, ctrl floc.Control) error {
		titles := ctx.Value(keyTitles).(*Titles)
		detailFuncs := make([]floc.Job, len(titles.Titles))
		for i, t := range titles.Titles {
			index, id := i, t.ID // *here[Point] same index and id if not set here
			detailFuncs[index] = func(ctx floc.Context, ctrl floc.Control) error {
				log.Println("doTitle detail", index, id)
				titles.Titles[index] = *reqTitle(strconv.Itoa(id))
				return nil
			}
		}
		job := run.Parallel(detailFuncs...)
		_, _, err := floc.Run(job)
		return err
	}
	doEpisodes := func(ctx floc.Context, ctrl floc.Control) error {
		log.Println("doEpisodes")
		episodes := ctx.Value(keyEpisodes).([]*Episodes)
		detailFuncs := make([]floc.Job, len(episodes))
		for i, e := range episodes {
			index, titleID := i, e.TitleID // *here[Point] same index and id if not set here
			detailFuncs[index] = func(ctx floc.Context, ctrl floc.Control) error {
				log.Println("doEpisodes", index, titleID)
				episodes[index].Episodes = reqEpisodes(strconv.Itoa(titleID)).Episodes
				err := runEpisodeDetails(titleID, episodes[index])
				return err
			}
		}
		job := run.Parallel(detailFuncs...)
		_, _, err := floc.Run(job)
		return err
	}
	doMerge := func(ctx floc.Context, ctrl floc.Control) error {
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

	job := run.Sequence(
		doTitles,
		run.Parallel(
			doTitle,
			doEpisodes,
		),
		doMerge,
	)
	ctx := floc.NewContext()
	ctrl := floc.NewControl(ctx)
	ctx.AddValue(keyAPIResult, &APIResult{})

	_, _, err := floc.RunWith(ctx, ctrl, job)
	if err != nil {
		panic(err)
	}
	titles := ctx.Value(keyTitles).(*Titles)
	episodes := ctx.Value(keyEpisodes).([]*Episodes)
	log.Printf("titles = %+v / %+v\n", titles, episodes)
	for _, es := range episodes {
		log.Printf("es = %+v\n", es)
		for _, e := range es.Episodes {
			log.Printf("e = %+v\n", e)
		}
	}

	result := ctx.Value(keyAPIResult).(*APIResult)
	resp, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	log.Printf("resp = %+v\n", string(resp))
}

func runEpisodeDetails(titleID int, es *Episodes) error {
	detailFuncs := make([]floc.Job, len(es.Episodes))
	for i, e := range es.Episodes {
		index, episodeID := i, e.ID // *here[Point] same index and id if not set here
		detailFuncs[index] = func(ctx floc.Context, ctrl floc.Control) error {
			log.Println("doEpisodeDetail ", index, titleID, episodeID)
			es.Episodes[index] = reqEpisode(strconv.Itoa(titleID), strconv.Itoa(episodeID))
			return nil
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

func doGet(url string, result interface{}) {
	_, err := resty.R().SetResult(result).Get(url)
	if err != nil {
		panic(err)
	}
}

func reqTitles() *Titles {
	url := apiBase + "/titles"
	titles := &Titles{}
	doGet(url, titles)
	return titles
}

func reqTitle(id string) *Title {
	url := apiBase + "/titles/" + id
	title := &Title{}
	doGet(url, title)
	return title
}

func reqEpisodes(titleID string) *Episodes {
	url := apiBase + "/titles/" + titleID + "/episodes"
	episodes := &Episodes{}
	doGet(url, episodes)
	return episodes
}

func reqEpisode(titleID, episodeID string) *Episode {
	url := apiBase + "/titles/" + titleID + "/episodes/" + episodeID
	episode := &Episode{}
	doGet(url, episode)
	return episode
}
