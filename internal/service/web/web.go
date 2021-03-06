package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"GoTuringCoffee/internal/service/barista"
	"GoTuringCoffee/internal/service/lib"
	"GoTuringCoffee/internal/service/web/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	nats "github.com/nats-io/nats.go"
)

type WebConfig struct {
	StaticFilePath string
	Port           int
}

type CustomContext struct {
	echo.Context
	cookbookModel *model.CookbookModel
	machineModel  *model.Machine
	context       context.Context
	nc            *nats.EncodedConn
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

type CookbookJson struct {
	ID          string        `json:"id,omitempty"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Processes   []ProcessJson `json:"processes"`
	TotalTime   float64       `json:"time"`
	TotalWater  float64       `json:"water"`
}

func NewCookbookJson(cookbook lib.Cookbook) (cj CookbookJson) {
	cj.ID = cookbook.ID
	cj.Name = cookbook.Name
	cj.Description = cookbook.Description
	cj.TotalTime = cookbook.GetTotalTime()
	cj.TotalWater = cookbook.GetTotalWater()
	for _, p := range cookbook.Processes {
		cj.Processes = append(cj.Processes, NewProcessJson(&p))
	}
	return
}

type ProcessJson struct {
	Name   string          `json:"name"`
	Params json.RawMessage `json:"params"`
}

func NewProcessJson(process *lib.Process) (pj ProcessJson) {
	switch (*process).(type) {
	case *lib.Circle:
		pj.Name = "Circle"
	case *lib.Spiral:
		pj.Name = "Spiral"
	case *lib.Polygon:
		pj.Name = "Polygon"
	case *lib.Fixed:
		pj.Name = "FixedPoint"
	case *lib.Move:
		pj.Name = "Move"
	case *lib.Wait:
		pj.Name = "Wait"
	case *lib.Mix:
		pj.Name = "Mix"
	case *lib.Home:
		pj.Name = "Home"
	}
	pj.Params, _ = json.Marshal(*process)
	return
}

func DecodeProcess(pj *ProcessJson) (p lib.Process) {
	switch pj.Name {
	case "Circle":
		p = new(lib.Circle)
	case "Spiral":
		p = new(lib.Spiral)
	case "FixedPoint":
		p = new(lib.Fixed)
	case "Move":
		p = new(lib.Move)
	case "Wait":
		p = new(lib.Wait)
	case "Mix":
		p = new(lib.Mix)
	case "Home":
		p = new(lib.Home)
	}
	json.Unmarshal(pj.Params, p)
	return
}

func (s *Service) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	cookbookModel := model.NewCookbookModel(&s.DB)
	machineModel := model.NewMachine(ctx, nc)
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1234"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := CustomContext{c, cookbookModel, machineModel, ctx, nc}
			return h(cc)
		}
	})

	e.GET("/api/cookbooks", s.ListCookbook)
	e.POST("/api/cookbooks", s.CreateCookbook)
	e.GET("/api/cookbooks/:id", s.GetCookbook)
	e.GET("/api/cookbooks/:id/points", s.GetCookbookPoints)
	e.PUT("/api/cookbooks/:id", s.UpdateCookbook)
	e.DELETE("/api/cookbooks/:id", s.DeleteCookbook)
	e.GET("/api/machine", s.GetMachineStatus)
	e.POST("/api/barista/:id/brew", s.BrewCookbook)
	e.GET("/api/processes", s.ListProcesses)
	e.GET("/api/processes/:name", s.GetDefaultProcess)
	e.GET("/api/processes/all", s.GetAllDefaultProcesses)
	e.GET("/api/processes/units", s.GetAllFieldUnits)
	e.Static("/", s.Web.StaticFilePath)
	e.PUT("/api/machine/tank/temperature", s.SetTargetTemperature)

	go func() {
		if err = e.Start(fmt.Sprintf(":%d", s.Web.Port)); err != nil {
			e.Logger.Info(err)
		}
	}()

	timer := time.NewTimer(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			e.Logger.Info("stoping web service")
			err = e.Shutdown(ctx)
			defer func() { fin <- struct{}{} }()
			e.Logger.Info("stop web service")
			return
		case <-timer.C:
			timer = time.NewTimer(1 * time.Second)
		}
	}
}

func (s *Service) ListCookbook(c echo.Context) (err error) {
	var cookbookJsons []CookbookJson
	var cookbooks []lib.Cookbook

	cc := c.(CustomContext)
	cookbooks, err = cc.cookbookModel.ListCookbooks()
	if err != nil {
		return err
	}
	for _, c := range cookbooks {
		cookbookJsons = append(cookbookJsons, NewCookbookJson(c))
	}
	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Payload: cookbookJsons,
	})
}

func (s *Service) CreateCookbook(c echo.Context) (err error) {
	cc := c.(CustomContext)
	cookbookJson := new(CookbookJson)
	var cookbook lib.Cookbook

	if err := cc.Bind(cookbookJson); err != nil {
		return err
	}
	cookbook.Name = cookbookJson.Name
	cookbook.Description = cookbookJson.Description
	for _, pj := range cookbookJson.Processes {
		cookbook.Processes = append(cookbook.Processes, DecodeProcess(&pj))
	}
	if err := cc.cookbookModel.CreateCookbook(&cookbook); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Response{
		Status: 200,
	})
}

func (s *Service) GetCookbook(c echo.Context) (err error) {
	cc := c.(CustomContext)
	id := cc.Param("id")
	cookbook, err := cc.cookbookModel.GetCookbook(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Payload: NewCookbookJson(cookbook),
	})
}

func (s *Service) UpdateCookbook(c echo.Context) error {
	cc := c.(CustomContext)
	cookbookJson := new(CookbookJson)
	var cookbook lib.Cookbook

	if err := cc.Bind(cookbookJson); err != nil {
		return err
	}
	cookbook.ID = cookbookJson.ID
	cookbook.Name = cookbookJson.Name
	cookbook.Description = cookbookJson.Description
	for _, pj := range cookbookJson.Processes {
		cookbook.Processes = append(cookbook.Processes, DecodeProcess(&pj))
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

func (s *Service) GetCookbookPoints(c echo.Context) error {
	var cookbook lib.Cookbook
	var err error

	cc := c.(CustomContext)
	id := cc.Param("id")
	if cookbook, err = cc.cookbookModel.GetCookbook(id); err != nil {
		return err
	}

	points := cookbook.ToPoints()

	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Payload: points,
	})
}

func (s *Service) BrewCookbook(c echo.Context) error {
	var cookbook lib.Cookbook
	var err error

	cc := c.(CustomContext)
	id := cc.Param("id")
	if cookbook, err = cc.cookbookModel.GetCookbook(id); err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(cc.context, 2*time.Second)
	resp, err := barista.Brew(ctx, cc.nc, cookbook.ToPoints())
	if resp.IsFailure() {
		return c.JSON(http.StatusInternalServerError, Response{
			Status: 500,
		})
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

func (s *Service) Stop() error {
	return nil
}

func (s *Service) ListProcesses(c echo.Context) error {
	cc := c.(CustomContext)

	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Payload: cc.cookbookModel.GetProcessNameList(),
	})
}

func (s *Service) GetDefaultProcess(c echo.Context) error {
	cc := c.(CustomContext)
	name := cc.Param("name")
	process := cc.cookbookModel.GetDefaultProcess(name)
	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Payload: NewProcessJson(&process),
	})
}

func (s *Service) GetAllDefaultProcesses(c echo.Context) error {
	cc := c.(CustomContext)
	processes := cc.cookbookModel.GetAllDefaultProcesses()

	result := map[string]ProcessJson{}
	for name, process := range processes {
		result[name] = NewProcessJson(&process)
	}

	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Payload: result,
	})
}

func (s *Service) GetAllFieldUnits(c echo.Context) error {
	cc := c.(CustomContext)

	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Payload: cc.cookbookModel.GetAllFieldUnits(),
	})
}

type SetTemperaturePayload struct {
	Temperature float64 `json:"temperature"`
}

func (s *Service) SetTargetTemperature(c echo.Context) error {
	cc := c.(CustomContext)

	var payload SetTemperaturePayload
	if err := cc.Bind(&payload); err != nil {
		return err
	}

	if err := cc.machineModel.SetTargetTemperature(payload.Temperature); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Response{
		Status: 200,
	})
}
