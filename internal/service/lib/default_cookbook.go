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
	circle := Process{
		BaseItem: BaseItem{
			ID:          "DefaultProcessCircle",
			Name:        "Circle",
			Description: "Default Circle",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Tags:        []string{},
			Notes:       []string{},
		},
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

	return circle
}

func defaultSpiral() Process {
	return Process{
		BaseItem: BaseItem{
			ID:          "DefaultProcessSpiral",
			Name:        "Spiral",
			Description: "Default Spiral",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Tags:        []string{},
			Notes:       []string{},
		},
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
		BaseItem: BaseItem{
			ID:          "DefaultProcessPolygon",
			Name:        "Polygon",
			Description: "Default Polygon",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Tags:        []string{},
			Notes:       []string{},
		},
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
		BaseItem: BaseItem{
			ID:          "DefaultProcessFixed",
			Name:        "Fixed",
			Description: "Default Fixed",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Tags:        []string{},
			Notes:       []string{},
		},
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
		BaseItem: BaseItem{
			ID:          "DefaultProcessHome",
			Name:        "Home",
			Description: "Default Home",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Tags:        []string{},
			Notes:       []string{},
		},
		Impl: &Home{},
	}
}

func defaultMove() Process {
	return Process{
		BaseItem: BaseItem{
			ID:          "DefaultProcessMove",
			Name:        "Move",
			Description: "Default Move",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Tags:        []string{},
			Notes:       []string{},
		},
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
		BaseItem: BaseItem{
			ID:          "DefaultProcessWait",
			Name:        "Wait",
			Description: "Default Wait",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Tags:        []string{},
			Notes:       []string{},
		},
		Impl: &Wait{
			Time: 10,
		},
	}
}

func defaultMix() Process {
	return Process{
		BaseItem: BaseItem{
			ID:          "DefaultProcessMix",
			Name:        "Mix",
			Description: "Default Mix",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Tags:        []string{},
			Notes:       []string{},
		},
		Impl: &Mix{
			Temperature: 80,
		},
	}
}

// GenerateDefaultCookbook Generate a cookbook with default processes
func GenerateDefaultCookbook() Cookbook {
	return Cookbook{
		BaseItem: BaseItem{
			ID:          DefaultCookbookID,
			Name:        "Default Cookbook",
			Description: "Cookbook",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Tags:        []string{},
			Notes:       []string{},
		},
		Processes: GetDefaultProcesses(),
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
