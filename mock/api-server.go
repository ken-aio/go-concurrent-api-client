package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/api/titles", func(c echo.Context) error {
		f := "titles.json"
		log(c, f)
		return c.JSONBlob(http.StatusOK, read(f))
	})
	e.GET("/api/titles/:id", func(c echo.Context) error {
		id := c.Param("id")
		f := "titles/" + id + ".json"
		log(c, f)
		return c.JSONBlob(http.StatusOK, read(f))
	})
	e.GET("/api/titles/:title_id/episodes", func(c echo.Context) error {
		titleID := c.Param("title_id")
		f := "titles/" + titleID + "/episodes.json"
		log(c, f)
		return c.JSONBlob(http.StatusOK, read(f))
	})
	e.GET("/api/titles/:title_id/episodes/:id", func(c echo.Context) error {
		titleID := c.Param("title_id")
		episodeID := c.Param("id")
		f := "titles/" + titleID + "/episodes/" + episodeID + ".json"
		log(c, f)
		return c.JSONBlob(http.StatusOK, read(f))
	})

	fmt.Printf("server start in port: 9998\n")
	if err := e.Start(":9998"); err != nil {
		fmt.Printf("err = %+v\n", err)
	}
}

func read(s string) []byte {
	data, err := ioutil.ReadFile("./" + s)
	if err != nil {
		panic(err)
	}
	return (data)
}

func log(c echo.Context, fn string) {
	fmt.Printf("%+v %+v read : %+v\n", time.Now(), c.Path(), fn)
}
