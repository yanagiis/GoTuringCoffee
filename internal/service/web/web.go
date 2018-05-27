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

func (s *Service) Run(ctx context.Context, nc *nats.EncodedConn) (err error) {
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
	e.GET("/api/cookbooks/{id}", s.GetCookbook)
	e.PUT("/api/cookbooks/{id}", s.UpdateCookbook)
	e.DELETE("/api/cookbooks/{id}", s.DeleteCookbook)
	e.GET("/api/machine", s.GetMachineStatus)
	// e.PUT("/api/machine/tank/temperature", s.SetTargetTemperature)
	if err = e.Start(fmt.Sprintf(":%d", s.Web.Port)); err != nil {
		e.Logger.Fatal(err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			err = e.Shutdown(ctx)
		case <-time.After(time.Second):
		}
	}
}

func (s *Service) ListCookbook(c echo.Context) error {
	cc := c.(CustomContext)
	cookbooks, err := cc.cookbookModel.ListCookbooks()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cookbooks)
}

func (s *Service) GetCookbook(c echo.Context) error {
	cc := c.(CustomContext)
	id := cc.Param("id")
	cookbook, err := cc.cookbookModel.GetCookbook(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cookbook)
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
	return c.JSON(http.StatusOK, "")
}

func (s *Service) DeleteCookbook(c echo.Context) error {
	cc := c.(CustomContext)
	id := cc.Param("id")
	err := cc.cookbookModel.DeleteCookbook(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "")
}

func (s *Service) GetMachineStatus(c echo.Context) error {
	cc := c.(CustomContext)
	status, err := cc.machineModel.GetMachineStatus()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, status)
}

// func (s *Service) SetTargetTemperature(c echo.Context) error {
// }
