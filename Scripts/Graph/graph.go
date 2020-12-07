package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	pusher "github.com/pusher/pusher-http-go"
)

// We register the Pusher client
var client = pusher.Client{
	AppID:   "1118979",
	Key:     "7befe6ab035a03a2ada9",
	Secret:  "1605eea01f08a987d955",
	Cluster: "us2",
	Secure:  true,
}

// visitsData is a struct
type visitsData struct {
	Pages int
	Count int
}

func setInterval(ourFunc func(), milliseconds int, async bool) chan bool {

	// How often to fire the passed in function in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticker and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	// Put the selection in a go routine so that the for loop is none blocking
	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					// This won't block
					go ourFunc()
				} else {
					// This will block
					ourFunc()
				}
			case <-clear:
				ticker.Stop()
				return
			}
		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clear
}
func simulate(c echo.Context) error {
	setInterval(func() {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		newVisitsData := visitsData{
			Pages: r1.Intn(100),
			Count: r1.Intn(100),
		}
		client.Trigger("visitorsCount", "addNumber", newVisitsData)
	}, 5, true)

	return c.String(http.StatusOK, "Simulation begun")
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define the HTTP routes
	e.File("/", "Frontend/index.html")
	e.File("/style.css", "Frontend/style.css")
	e.File("/app.js", "Frontend/app.js")
	e.GET("/simulate", simulate)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
