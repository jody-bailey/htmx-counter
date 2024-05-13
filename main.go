package main

import (
	"fmt"
	"html/template"
	"log"
	"time"

	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/r3labs/sse/v2"
)

var Count = 0

type Data struct {
	Name  string //exported field since it begins with a capital letter
	Count int
}

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func StartTimer() {
}

func main() {
	// Sse server
	server := sse.New()
	server.AutoReplay = false // copy paste because we got that dog in us
	_ = server.CreateStream("counter")

	// Echo instance
	e := echo.New()
	e.Renderer = NewTemplates()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/sse", func(c echo.Context) error { // longer variant with disconnect logic
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			start := time.Now()
			defer ticker.Stop()
			for {
				select {
				case <-c.Request().Context().Done():
					log.Printf("SSE client disconnected, ip: %v", c.RealIP())
					return
				case <-ticker.C:
					now := time.Now()
					duration := int(now.Sub(start).Seconds())
					server.Publish("counter", &sse.Event{
						Event: []byte("timer"),
						Data:  []byte(fmt.Sprintf("%d seconds elapsed", duration)),
					})
				}
			}
		}()

		server.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// Routes
	e.GET("/", func(c echo.Context) error {
		Count = 0
		template := "index"
		return c.Render(http.StatusOK, template, Data{
			Name:  "Daryl",
			Count: Count,
		})
	})

	e.GET("/counter", func(c echo.Context) error {
		template := "counter"
		return c.Render(http.StatusOK, template, Data{
			Count: Count,
		})
	})

	e.GET("/increment", func(c echo.Context) error {
		Count++
		// return c.HTML(http.StatusOK, fmt.Sprint(Count))
		server.Publish("counter", &sse.Event{
			Event: []byte("counter"),
			Data:  []byte(fmt.Sprint(Count)),
		})
		return nil
	})

	e.GET("/decrement", func(c echo.Context) error {
		Count--
		server.Publish("counter", &sse.Event{
			Event: []byte("counter"),
			Data:  []byte(fmt.Sprint(Count)),
		})
		return nil
	})

	e.GET("/reset", func(c echo.Context) error {
		Count = 0
		server.Publish("counter", &sse.Event{
			Event: []byte("counter"),
			Data:  []byte(fmt.Sprint(Count)),
		})
		return nil
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
