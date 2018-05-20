package web

***REMOVED***
	"context"
***REMOVED***
	"time"

	"github.com/labstack/echo"
	nats "github.com/nats-io/go-nats"
***REMOVED***

type MongoDBConfig struct {
	Url string
***REMOVED***

type WebConfig struct {
	StaticFilePath string
	Port           int
***REMOVED***

type Service struct {
	DB  MongoDBConfig
	Web WebConfig
***REMOVED***

func (s *Service***REMOVED*** Run(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (err error***REMOVED*** {
	model := NewModel(&s.DB***REMOVED***

	e := echo.New(***REMOVED***
	e.Static("/", s.Web.StaticFilePath***REMOVED***
	if err = e.Start(fmt.Sprintf(":%d", s.Web.Port***REMOVED******REMOVED***; err != nil {
		e.Logger.Fatal(err***REMOVED***
		return
***REMOVED***

	for {
		select {
		case <-ctx.Done(***REMOVED***:
			err = e.Shutdown(ctx***REMOVED***
		case <-time.After(time.Second***REMOVED***:
	***REMOVED***
***REMOVED***
***REMOVED***
