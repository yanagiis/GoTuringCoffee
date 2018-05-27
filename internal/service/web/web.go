package web

***REMOVED***
	"context"
***REMOVED***
	"net/http"
	"time"

	"github.com/labstack/echo"
	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
***REMOVED***
***REMOVED***

type WebConfig struct {
	StaticFilePath string
	Port           int
***REMOVED***

type CustomContext struct {
	echo.Context
	cookbookModel *model.Cookbook
	machineModel  *model.Machine
***REMOVED***

type Service struct {
	DB  model.MongoDBConfig
	Web WebConfig
***REMOVED***

func (s *Service***REMOVED*** Run(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (err error***REMOVED*** {
	cookbookModel := model.NewCookbook(&s.DB***REMOVED***
	machineModel := model.NewMachine(ctx, nc***REMOVED***
	e := echo.New(***REMOVED***
	e.Use(func(h echo.HandlerFunc***REMOVED*** echo.HandlerFunc {
		return func(c echo.Context***REMOVED*** error {
			cc := CustomContext{c, cookbookModel, machineModel***REMOVED***
			return h(cc***REMOVED***
	***REMOVED***
***REMOVED******REMOVED***
	e.Static("/", s.Web.StaticFilePath***REMOVED***
	e.GET("/api/cookbooks", s.ListCookbook***REMOVED***
	e.GET("/api/cookbooks/{id***REMOVED***", s.GetCookbook***REMOVED***
	e.PUT("/api/cookbooks/{id***REMOVED***", s.UpdateCookbook***REMOVED***
	e.DELETE("/api/cookbooks/{id***REMOVED***", s.DeleteCookbook***REMOVED***
	e.GET("/api/machine", s.GetMachineStatus***REMOVED***
	// e.PUT("/api/machine/tank/temperature", s.SetTargetTemperature***REMOVED***
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

func (s *Service***REMOVED*** ListCookbook(c echo.Context***REMOVED*** error {
	cc := c.(*CustomContext***REMOVED***
	cookbooks, err := cc.cookbookModel.ListCookbooks(***REMOVED***
***REMOVED***
		return err
***REMOVED***
	return c.JSON(http.StatusOK, cookbooks***REMOVED***
***REMOVED***

func (s *Service***REMOVED*** GetCookbook(c echo.Context***REMOVED*** error {
	cc := c.(*CustomContext***REMOVED***
	id := cc.Param("id"***REMOVED***
	cookbook, err := cc.cookbookModel.GetCookbook(id***REMOVED***
***REMOVED***
		return err
***REMOVED***
	return c.JSON(http.StatusOK, cookbook***REMOVED***
***REMOVED***

func (s *Service***REMOVED*** UpdateCookbook(c echo.Context***REMOVED*** error {
	cc := c.(*CustomContext***REMOVED***
	var cookbook lib.Cookbook
	if err := cc.Bind(cookbook***REMOVED***; err != nil {
		return err
***REMOVED***

	id := cc.Param("id"***REMOVED***
	if err := cc.cookbookModel.UpdateCookbook(id, &cookbook***REMOVED***; err != nil {
		return err
***REMOVED***
	return c.JSON(http.StatusOK, ""***REMOVED***
***REMOVED***

func (s *Service***REMOVED*** DeleteCookbook(c echo.Context***REMOVED*** error {
	cc := c.(*CustomContext***REMOVED***
	id := cc.Param("id"***REMOVED***
	err := cc.cookbookModel.DeleteCookbook(id***REMOVED***
***REMOVED***
		return err
***REMOVED***
	return c.JSON(http.StatusOK, ""***REMOVED***
***REMOVED***

func (s *Service***REMOVED*** GetMachineStatus(c echo.Context***REMOVED*** error {
	cc := c.(*CustomContext***REMOVED***
	status, err := cc.machineModel.GetMachineStatus(***REMOVED***
***REMOVED***
		return err
***REMOVED***
	return c.JSON(http.StatusOK, status***REMOVED***
***REMOVED***

// func (s *Service***REMOVED*** SetTargetTemperature(c echo.Context***REMOVED*** error {
// ***REMOVED***
