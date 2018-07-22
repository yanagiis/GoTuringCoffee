package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/web/model"
)

type WebConfig struct {
	StaticFilePath string
	Port           int
}

type CustomContext struct {
	echo.Context
	cookbookModel *model.Cookbook
	machineModel  *model.Machine
}

type Service struct {
	DB  model.MongoDBConfig
	Web WebConfig
}

type Response struct {
	Status  int64       `json:"status"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func (s *Service) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	cookbookModel := model.NewCookbook(&s.DB)
	machineModel := model.NewMachine(ctx, nc)
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := CustomContext{c, cookbookModel, machineModel}
			return h(cc)
		}
	})
	e.Static("/", s.Web.StaticFilePath)
	e.GET("/api/cookbooks", s.ListCookbook)
	e.GET("/api/cookbooks/:id", s.GetCookbook)
	e.PUT("/api/cookbooks/:id", s.UpdateCookbook)
	e.DELETE("/api/cookbooks/:id", s.DeleteCookbook)
	e.GET("/api/machine", s.GetMachineStatus)
	e.POST("/api/barista/:id", s.BrewCookbook)
	// e.PUT("/api/machine/tank/temperature", s.SetTargetTemperature)

	go func() {
		if err = e.Start(fmt.Sprintf(":%d", s.Web.Port)); err != nil {
			e.Logger.Info(err)
		}
	}()

	timer := time.NewTimer(1)
	for {
		select {
		case <-ctx.Done():
			e.Logger.Info("stoping web service")
			err = e.Shutdown(ctx)
			defer func() { fin <- struct{}{} }()
			e.Logger.Info("stop web service")
			return
		case <-timer.C:
			timer = time.NewTimer(1)
		}
	}
}

func (s *Service) ListCookbook(c echo.Context) (err error) {
	var cookbooks []lib.Cookbook

	cc := c.(CustomContext)
	cookbooks, err = cc.cookbookModel.ListCookbooks()
	if err != nil {
		return
	}
	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Payload: cookbooks,
	})
}

func (s *Service) GetCookbook(c echo.Context) (err error) {
	var cookbook *lib.Cookbook

	cc := c.(CustomContext)
	id := cc.Param("id")
	if cookbook, err = cc.cookbookModel.GetCookbook(id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Payload: cookbook,
	})
}

func (s *Service) UpdateCookbook(c echo.Context) error {
	cc := c.(CustomContext)
	var cookbook lib.Cookbook
	if err := cc.Bind(cookbook); err != nil {
		return err
	}

	id := cc.Param("id")
	if err := cc.cookbookModel.UpdateCookbook(id, &cookbook); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Status: 200,
	})
}

func (s *Service) DeleteCookbook(c echo.Context) error {
	cc := c.(CustomContext)
	id := cc.Param("id")
	err := cc.cookbookModel.DeleteCookbook(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Status: 200,
	})
}

func (s *Service) BrewCookbook(c echo.Context) error {
	// var cookbook *lib.Cookbook
	var err error

	cc := c.(CustomContext)
	id := cc.Param("id")
	if _, err = cc.cookbookModel.GetCookbook(id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Status: 200,
	})
}

func (s *Service) GetMachineStatus(c echo.Context) error {
	cc := c.(CustomContext)
	payload, err := cc.machineModel.GetMachineStatus()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Payload: payload,
	})
}

// func (s *Service) SetTargetTemperature(c echo.Context) error {
// }
