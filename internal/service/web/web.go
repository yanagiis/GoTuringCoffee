package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"GoTuringCoffee/internal/service/barista"
	"GoTuringCoffee/internal/service/lib"
	"GoTuringCoffee/internal/service/web/model"
	dbmodel "GoTuringCoffee/internal/service/web/model"
	"GoTuringCoffee/internal/service/web/model/repository"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	nats "github.com/nats-io/go-nats"
)

type WebConfig struct {
	StaticFilePath string
	Port           int
}

type CustomContext struct {
	echo.Context
	repoManager  *repository.RepositoryManager
	machineModel *dbmodel.Machine
	context      context.Context
	nc           *nats.EncodedConn
}

type Service struct {
	DBConfig  repository.MongoDBConfig
	WebConfig WebConfig
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
	Tags        []string      `json:"tags"`
	Notes       []string      `json:"notes"`
	Processes   []ProcessJson `json:"processes"`
	CreatedAt   int64         `json:"created_at`
	UpdatedAt   int64         `json:"updated_at`
	TotalTime   float64       `json:"time"`
	TotalWater  float64       `json:"water"`
}

func LibCookbookToJson(cookbook lib.Cookbook) (cj CookbookJson) {
	cj.ID = cookbook.ID
	cj.Name = cookbook.Name
	cj.Description = cookbook.Description

	cj.Tags = cookbook.Tags
	cj.Notes = cookbook.Notes

	cj.CreatedAt = cookbook.CreatedAt.Unix()
	cj.UpdatedAt = cookbook.UpdatedAt.Unix()

	cj.TotalTime = cookbook.GetTotalTime()
	cj.TotalWater = cookbook.GetTotalWater()

	for _, p := range cookbook.Processes {
		cj.Processes = append(cj.Processes, LibProcessToJson(p))
	}
	return
}

type ProcessJson struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	CreatedAt int64           `json:"created_at"`
	UpdatedAt int64           `json:"updated_at"`
	Impl      json.RawMessage `json:"impl"`
}

// LibProcessToJson Convert lib.Process to json
func LibProcessToJson(process lib.Process) (pj ProcessJson) {
	pj.ID = process.ID
	pj.Name = process.Name
	pj.CreatedAt = process.CreatedAt.Unix()
	pj.UpdatedAt = process.UpdatedAt.Unix()
	pj.Impl, _ = json.Marshal(process.Impl)
	return
}

// JsonProcessToLib Convert lib.process to json
func JsonProcessToLib(pj ProcessJson) (p lib.Process) {
	p.ID = pj.ID
	p.Name = pj.Name
	p.CreatedAt = time.Unix(pj.CreatedAt, 0)
	p.UpdatedAt = time.Unix(pj.UpdatedAt, 0)

	impl, err := lib.NewProcessImpl(pj.Name)
	if err != nil {
		return
	}
	json.Unmarshal(pj.Impl, impl)

	p.Impl = impl
	return
}

func (s *Service) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	repoManager, err := repository.NewRepositoryManager(ctx, &s.DBConfig)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	machineModel := model.NewMachine(ctx, nc)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1234"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := CustomContext{c, repoManager, machineModel, ctx, nc}
			return h(cc)
		}
	})

	e.GET("/api/cookbooks", s.ListCookbook)
	e.POST("/api/cookbooks", s.CreateCookbook)
	e.GET("/api/cookbooks/:id", s.GetCookbook)
	e.GET("/api/cookbooks/:id/points", s.GetCookbookPoints)
	e.PUT("/api/cookbooks/:id", s.UpdateCookbook)
	e.POST("/api/cookbooks/:id/cover", s.UploadCookbookCover)
	e.DELETE("/api/cookbooks/:id", s.DeleteCookbook)
	e.GET("/api/machine", s.GetMachineStatus)
	e.POST("/api/barista/:id/brew", s.BrewCookbook)
	e.GET("/api/processes", s.ListProcesses)
	e.GET("/api/processes/all", s.GetAllDefaultProcesses)
	e.GET("/api/processes/units", s.GetAllFieldUnits)
	e.Static("/", s.WebConfig.StaticFilePath)
	e.PUT("/api/machine/tank/temperature", s.SetTargetTemperature)

	go func() {
		if err = e.Start(fmt.Sprintf(":%d", s.WebConfig.Port)); err != nil {
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
	cookbooks, err = cc.repoManager.Cookbook.List(cc.context)
	if err != nil {
		return
	}
	for _, c := range cookbooks {
		cookbookJsons = append(cookbookJsons, LibCookbookToJson(c))
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
		cookbook.Processes = append(cookbook.Processes, JsonProcessToLib(pj))
	}
	if _, err := cc.repoManager.Cookbook.Create(cc.context, cookbook); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Response{
		Status: 200,
	})
}

func (s *Service) GetCookbook(c echo.Context) (err error) {
	var cookbook lib.Cookbook

	cc := c.(CustomContext)
	id := cc.Param("id")
	if cookbook, err = cc.repoManager.Cookbook.Get(cc.context, id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Payload: LibCookbookToJson(cookbook),
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
		cookbook.Processes = append(cookbook.Processes, JsonProcessToLib(pj))
	}

	_, err := cc.repoManager.Cookbook.Update(cc.context, cookbook)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Response{
		Status: 200,
	})
}

func (s *Service) DeleteCookbook(c echo.Context) error {
	cc := c.(CustomContext)
	id := cc.Param("id")
	err := cc.repoManager.Cookbook.Delete(cc.context, id)
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
	if cookbook, err = cc.repoManager.Cookbook.Get(cc.context, id); err != nil {
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
	if cookbook, err = cc.repoManager.Cookbook.Get(cc.context, id); err != nil {
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
		Payload: cc.repoManager.Cookbook.GetProcessNameList(),
	})
}

func (s *Service) GetAllDefaultProcesses(c echo.Context) error {
	cc := c.(CustomContext)
	processes, err := cc.repoManager.Cookbook.GetAllDefaultProcesses(cc.context)
	if err != nil {
		return err
	}

	result := map[string]ProcessJson{}
	for _, process := range processes {
		result[process.Name] = LibProcessToJson(process)
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
		Payload: cc.repoManager.Cookbook.GetAllFieldUnits(),
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

func (s *Service) UploadCookbookCover(c echo.Context) error {
	fmt.Printf("Uploading cookbook cover")
	id := c.Param("id")
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Description
	dst, err := os.Create(id + "_" + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Response{
		Status: 200,
	})
}
