package lib

import (
	"time"
)

type ProcessIDType string

// Generate from
// src := []byte("DefaultCB012")
// encodedStr := hex.EncodeToString(src)
const DefaultCookbookID = "44656661756c744342303132"

const (
	DefaultCircle  ProcessIDType = "Default Circle"
	DefaultSpiral  ProcessIDType = "Default Spiral"
	DefaultPolygon ProcessIDType = "Default Polygon"
	DefaultFixed   ProcessIDType = "Default Fixed"
	DefaultHome    ProcessIDType = "Default Home"
	DefaultMove    ProcessIDType = "Default Move"
	DefaultWait    ProcessIDType = "Default Wait"
	DefaultMix     ProcessIDType = "Default Mix"
)

func defaultCircle() Process {
	return Process{
		ID:        "DefaultProcessCircle",
		Name:      "Circle",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Impl: &Circle{
			Coords: Coordinate{
				X: 0,
				Y: 0,
				Z: 200,
			},
			ToZ: 0,
			Radius: Range{
				From: 0,
				To:   2,
			},
			Cylinder:    5,
			Time:        10,
			Water:       150,
			Temperature: 80,
		},
	}
}

func defaultSpiral() Process {
	return Process{
		ID:        "DefaultProcessSpiral",
		Name:      "Spiral",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Impl: &Spiral{
			Coords: Coordinate{
				X: 0,
				Y: 0,
				Z: 200,
			},
			ToZ: 0,
			Radius: Range{
				From: 0,
				To:   2,
			},
			Cylinder:    1,
			Time:        10,
			Water:       150,
			Temperature: 80,
		},
	}
}

func defaultPolygon() Process {
	return Process{
		ID:        "DefaultProcessPolygon",
		Name:      "Polygon",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Impl: &Polygon{
			Coords: Coordinate{
				X: 0,
				Y: 0,
				Z: 200,
			},
			ToZ: 0,
			Radius: Range{
				From: 0,
				To:   2,
			},
			Polygon:     2,
			Cylinder:    5,
			Time:        10,
			Water:       150,
			Temperature: 80,
		},
	}
}

func defaultFixed() Process {
	return Process{
		ID:        "DefaultProcessFixed",
		Name:      "Fixed",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Impl: &Fixed{
			Coords: Coordinate{
				X: 0,
				Y: 0,
				Z: 200,
			},
			Time:        10,
			Water:       150,
			Temperature: 80,
		},
	}
}

func defaultHome() Process {
	return Process{
		ID:        "DefaultProcessHome",
		Name:      "Home",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Impl:      &Home{},
	}
}

func defaultMove() Process {
	return Process{
		ID:        "DefaultProcessMove",
		Name:      "Move",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Impl: &Move{
			Coords: Coordinate{
				X: 0,
				Y: 0,
				Z: 200,
			},
		},
	}
}

func defaultWait() Process {
	return Process{
		ID:        "DefaultProcessWait",
		Name:      "Wait",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Impl: &Wait{
			Time: 10,
		},
	}
}

func defaultMix() Process {
	return Process{
		ID:        "DefaultProcessMix",
		Name:      "Mix",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Impl: &Mix{
			Temperature: 80,
		},
	}
}

// GenerateDefaultCookbook Generate a cookbook with default processes
func GenerateDefaultCookbook() Cookbook {
	return Cookbook{
		ID:          DefaultCookbookID,
		Name:        "Default Cookbook",
		Description: "Cookbook",
		Tags:        []string{},
		Notes:       []string{},
		Processes:   GetDefaultProcesses(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

// GetDefaultProcesses Return default processes
func GetDefaultProcesses() []Process {
	return []Process{
		defaultCircle(),
		defaultSpiral(),
		defaultPolygon(),
		defaultFixed(),
		defaultHome(),
		defaultMove(),
		defaultWait(),
		defaultMix(),
	}
}
