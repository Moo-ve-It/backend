package main

import (
	"net/http"
	"time"
)

// Cow represents a cow with sensor data
type Cow struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Tag         string     `json:"tag"`
	Location    Location   `json:"location"`
	Health      Health     `json:"health"`
	Sensors     CowSensors `json:"sensors"`
	LastUpdated time.Time  `json:"last_updated"`
}

// Location represents GPS coordinates
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Zone      string  `json:"zone"`
}

// Health represents health status
type Health struct {
	Status      string  `json:"status"`      // healthy, sick, injured
	Temperature float64 `json:"temperature"` // in Celsius
	HeartRate   int     `json:"heart_rate"`  // beats per minute
	Activity    string  `json:"activity"`    // grazing, resting, moving
}

// CowSensors represents sensor data from cow
type CowSensors struct {
	Temperature  float64 `json:"temperature"`
	HeartRate    int     `json:"heart_rate"`
	Activity     string  `json:"activity"`
	BatteryLevel int     `json:"battery_level"` // percentage
}

// RoboDog represents the robo-dog with sensor data
type RoboDog struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Status       string         `json:"status"` // active, idle, charging, maintenance
	Location     Location       `json:"location"`
	Sensors      RoboDogSensors `json:"sensors"`
	BatteryLevel int            `json:"battery_level"` // percentage
	LastUpdated  time.Time      `json:"last_updated"`
}

// RoboDogSensors represents sensor data from robo-dog
type RoboDogSensors struct {
	Temperature    float64 `json:"temperature"`
	Humidity       float64 `json:"humidity"`
	MotionDetected bool    `json:"motion_detected"`
	CameraStatus   string  `json:"camera_status"` // active, inactive
	AudioLevel     float64 `json:"audio_level"`   // decibels
}

// Drone represents the drone with sensor data
type Drone struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	Status       string       `json:"status"` // flying, landed, charging, maintenance
	Location     Location     `json:"location"`
	Altitude     float64      `json:"altitude"` // meters
	Sensors      DroneSensors `json:"sensors"`
	BatteryLevel int          `json:"battery_level"` // percentage
	LastUpdated  time.Time    `json:"last_updated"`
}

// DroneSensors represents sensor data from drone
type DroneSensors struct {
	Temperature  float64 `json:"temperature"`
	Humidity     float64 `json:"humidity"`
	WindSpeed    float64 `json:"wind_speed"`    // km/h
	CameraStatus string  `json:"camera_status"` // active, inactive
	GPSAccuracy  float64 `json:"gps_accuracy"`  // meters
	AirQuality   float64 `json:"air_quality"`   // AQI
}

// FarmState represents the overall state of the farm
type FarmState struct {
	TotalCows     int       `json:"total_cows"`
	HealthyCows   int       `json:"healthy_cows"`
	SickCows      int       `json:"sick_cows"`
	RoboDogStatus string    `json:"robodog_status"`
	DroneStatus   string    `json:"drone_status"`
	LastUpdated   time.Time `json:"last_updated"`
}

// Mock data storage
var mockCows = []Cow{
	{
		ID:   1,
		Name: "Bessie",
		Tag:  "COW-001",
		Location: Location{
			Latitude:  40.7128,
			Longitude: -74.0060,
			Zone:      "Pasture A",
		},
		Health: Health{
			Status:      "healthy",
			Temperature: 38.5,
			HeartRate:   65,
			Activity:    "grazing",
		},
		Sensors: CowSensors{
			Temperature:  38.5,
			HeartRate:    65,
			Activity:     "grazing",
			BatteryLevel: 85,
		},
		LastUpdated: time.Now(),
	},
	{
		ID:   2,
		Name: "Daisy",
		Tag:  "COW-002",
		Location: Location{
			Latitude:  40.7130,
			Longitude: -74.0062,
			Zone:      "Pasture A",
		},
		Health: Health{
			Status:      "healthy",
			Temperature: 38.7,
			HeartRate:   70,
			Activity:    "resting",
		},
		Sensors: CowSensors{
			Temperature:  38.7,
			HeartRate:    70,
			Activity:     "resting",
			BatteryLevel: 92,
		},
		LastUpdated: time.Now(),
	},
	{
		ID:   3,
		Name: "Moo",
		Tag:  "COW-003",
		Location: Location{
			Latitude:  40.7125,
			Longitude: -74.0058,
			Zone:      "Pasture B",
		},
		Health: Health{
			Status:      "sick",
			Temperature: 39.8,
			HeartRate:   85,
			Activity:    "resting",
		},
		Sensors: CowSensors{
			Temperature:  39.8,
			HeartRate:    85,
			Activity:     "resting",
			BatteryLevel: 78,
		},
		LastUpdated: time.Now(),
	},
	{
		ID:   4,
		Name: "Clover",
		Tag:  "COW-004",
		Location: Location{
			Latitude:  40.7135,
			Longitude: -74.0065,
			Zone:      "Pasture B",
		},
		Health: Health{
			Status:      "healthy",
			Temperature: 38.4,
			HeartRate:   62,
			Activity:    "moving",
		},
		Sensors: CowSensors{
			Temperature:  38.4,
			HeartRate:    62,
			Activity:     "moving",
			BatteryLevel: 88,
		},
		LastUpdated: time.Now(),
	},
	{
		ID:   5,
		Name: "Buttercup",
		Tag:  "COW-005",
		Location: Location{
			Latitude:  40.7120,
			Longitude: -74.0063,
			Zone:      "Pasture A",
		},
		Health: Health{
			Status:      "healthy",
			Temperature: 38.6,
			HeartRate:   68,
			Activity:    "grazing",
		},
		Sensors: CowSensors{
			Temperature:  38.6,
			HeartRate:    68,
			Activity:     "grazing",
			BatteryLevel: 90,
		},
		LastUpdated: time.Now(),
	},
}

var mockRoboDog = RoboDog{
	ID:     1,
	Name:   "Rex",
	Status: "active",
	Location: Location{
		Latitude:  40.7129,
		Longitude: -74.0061,
		Zone:      "Central Area",
	},
	Sensors: RoboDogSensors{
		Temperature:    22.5,
		Humidity:       65.0,
		MotionDetected: true,
		CameraStatus:   "active",
		AudioLevel:     45.2,
	},
	BatteryLevel: 72,
	LastUpdated:  time.Now(),
}

var mockDrone = Drone{
	ID:     1,
	Name:   "SkyEye",
	Status: "flying",
	Location: Location{
		Latitude:  40.7132,
		Longitude: -74.0059,
		Zone:      "Airspace",
	},
	Altitude: 150.0,
	Sensors: DroneSensors{
		Temperature:  18.3,
		Humidity:     58.0,
		WindSpeed:    12.5,
		CameraStatus: "active",
		GPSAccuracy:  2.5,
		AirQuality:   45.0,
	},
	BatteryLevel: 68,
	LastUpdated:  time.Now(),
}

// listCowsHandler returns a list of all cows with their sensor data
func (app *application) listCowsHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"cows":  mockCows,
		"total": len(mockCows),
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getCowHandler returns a specific cow by ID
func (app *application) getCowHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	for _, cow := range mockCows {
		if cow.ID == int(id) {
			env := envelope{"cow": cow}
			err := app.writeJSON(w, http.StatusOK, env, nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
			return
		}
	}

	app.notFoundResponse(w, r)
}

// getRoboDogHandler returns the robo-dog state and sensor data
func (app *application) getRoboDogHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{"robodog": mockRoboDog}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getDroneHandler returns the drone state and sensor data
func (app *application) getDroneHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{"drone": mockDrone}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getFarmStateHandler returns the overall farm state
func (app *application) getFarmStateHandler(w http.ResponseWriter, r *http.Request) {
	healthyCount := 0
	sickCount := 0
	for _, cow := range mockCows {
		if cow.Health.Status == "healthy" {
			healthyCount++
		} else if cow.Health.Status == "sick" {
			sickCount++
		}
	}

	farmState := FarmState{
		TotalCows:     len(mockCows),
		HealthyCows:   healthyCount,
		SickCows:      sickCount,
		RoboDogStatus: mockRoboDog.Status,
		DroneStatus:   mockDrone.Status,
		LastUpdated:   time.Now(),
	}

	env := envelope{"farm_state": farmState}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
