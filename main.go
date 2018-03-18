package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/CyCoreSystems/audimance/agenda"
	"github.com/CyCoreSystems/audimance/showtime"
	"github.com/boltdb/bolt"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"gopkg.in/inconshreveable/log15.v2"
)

var dbFile = "/var/db/ringfree/ipc.db"
var db *bolt.DB

// addr is the listen address
var addr string

// agiaddr is the listen address for the FastAGI service
var agiaddr string

// debug enables debug mode, which uses local files
// instead of bundled ones
var debug bool

// logger is the top-level logger
var logger log15.Logger

// ErrNilTarget indicates that the row/day/date has no target
// specified.
var ErrNilTarget = errors.New("Empty Target")

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// CustomerContext extends the Echo context to allow for custom data
type CustomContext struct {
	echo.Context

	Agenda *agenda.Agenda

	ShowTime *showtime.Service
}

func init() {
	flag.StringVar(&addr, "addr", ":9000", "Address binding")
	flag.StringVar(&agiaddr, "agiaddr", ":9001", "Address binding for FastAGI service")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode, which uses separate files for web development")
}

func main() {
	flag.Parse()

	// Create a logger
	logger = log15.New()

	// Read the Agenda
	a, err := agenda.New("agenda.yaml")
	if err != nil {
		fmt.Printf("failed to read agenda: %s", err.Error())
		os.Exit(1)
	}
	spew.Dump(a)

	// Create the showtime service
	svc := new(showtime.Service)
	go func() {
		err := svc.Run()
		if err != nil {
			fmt.Printf("showtime service died: %s", err.Error())
			os.Exit(1)
		}
	}()

	// Create web server
	e := echo.New()

	// Attach middleware
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			c := &CustomContext{
				Context:  ctx,
				Agenda:   a,
				ShowTime: svc,
			}
			return h(c)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Compile and attach templates
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	// Attach handlers

	// Handle the index
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", a)
	})

	e.Static("/app", "app")
	e.Static("/media", "media")

	e.GET("/room/:id", enterRoom)

	e.GET("/ws/performanceTime", performanceTime)

	// Listen to OS kill signals
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		e.Logger.Info("exiting on signal")
		os.Exit(100)
	}()

	// Listen for connections
	e.Logger.Debugf("listening on %s\n", addr)
	e.Logger.Fatal(e.Start(addr))
}

func enterRoom(c echo.Context) error {
	ctx := c.(*CustomContext)

	// Find our room
	id := ctx.Param("id")
	ctx.Logger().Errorf("received ID is %s", id)

	var r *agenda.Room
	for _, room := range ctx.Agenda.Rooms {
		if room.ID() == id {
			r = &room
			break
		}
	}
	if r == nil {
		return ctx.String(http.StatusNotFound, "no such room")
	}

	data := struct {
		Announcements []agenda.Announcement `json:"announcements"`
		Room          *agenda.Room          `json:"room"`
	}{
		Announcements: ctx.Agenda.Announcements,
		Room:          r,
	}

	return ctx.Render(200, "room.html", data) // TODO: template linking itself
}

func performanceTime(c echo.Context) error {
	ctx := c.(*CustomContext)

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// Create a subscription to the showtime service
		sub := ctx.ShowTime.Subscribe()
		defer sub.Cancel()

		for {
			// Process announcements
			ann := <-sub.C

			err := websocket.JSON.Send(ws, ann)
			if err != nil {
				ctx.Logger().Error(errors.Wrap(err, "failed to send announcement"))
				break
			}
		}
	}).ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}
