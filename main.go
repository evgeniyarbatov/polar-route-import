package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func main() {
	gpxInputFile := os.Args[1]
	gpxOutputFileName := os.Args[2]
	gpxOutputDir := os.Args[3]

	waypoints := GetWaypoints(gpxInputFile)

	sort.Slice(waypoints, func(i, j int) bool {
		return waypoints[j].Time.Before(waypoints[i].Time)
	})

	chunks := SplitWaypoints(waypoints, 500)

	for _, chunk := range chunks {
		firstWaypoint := chunk[0]
		lastWaypoint := chunk[len(chunk)-1]

		filename := fmt.Sprintf(
			gpxOutputDir+"/"+gpxOutputFileName+"%d_%dkm.gpx",
			firstWaypoint.Distance,
			lastWaypoint.Distance,
		)

		CreateGPXFile(filename, chunk)
	}
}

func GetWaypoints(
	gpxInputFile string,
) []Waypoint {
	file, _ := os.Open(gpxInputFile)
	xmlData, _ := io.ReadAll(file)

	var gpxData GPX
	xml.Unmarshal(xmlData, &gpxData)

	var waypoints []Waypoint
	var prevLat, prevLong float64

	for _, point := range gpxData.Trk.Trkseg.Trkpt {
		if prevLat == 0 || prevLong == 0 {
			prevLat = point.Lat
			prevLong = point.Lon
		}

		distance := GetHaversineDistance(
			prevLat,
			prevLong,
			point.Lat,
			point.Lon,
		)

		timestamp, _ := time.Parse(time.RFC3339Nano, point.Time)

		waypoints = append(
			waypoints,
			Waypoint{
				Time:      timestamp,
				Latitude:  point.Lat,
				Longitude: point.Lon,
				Distance:  distance,
			},
		)
	}

	return waypoints
}

func GetHaversineDistance(lat1, lon1, lat2, lon2 float64) int {
	const R = 6371000 // Earth radius in meters

	if lat1 == 0.0 || lon1 == 0.0 || lat2 == 0.0 || lon2 == 0.0 {
		return 0.0
	}

	var φ1 = lat1 * math.Pi / 180
	var φ2 = lat2 * math.Pi / 180
	var Δφ = (lat2 - lat1) * math.Pi / 180
	var Δλ = (lon2 - lon1) * math.Pi / 180

	var a = math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distanceMeters := R * c
	distanceKilometers := distanceMeters / 1000

	return int(math.Round(distanceKilometers))
}

func SplitWaypoints(
	waypoints []Waypoint,
	chunkSize int,
) [][]Waypoint {
	var chunks [][]Waypoint
	for i := 0; i < len(waypoints); i += chunkSize {
		end := i + chunkSize
		if end > len(waypoints) {
			end = len(waypoints)
		}
		chunks = append(chunks, waypoints[i:end])
	}
	return chunks
}

func CreateGPXFile(filename string, waypoints []Waypoint) {
	outputDir := filepath.Dir(filename)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Println("Error creating output directory:", err)
			return
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	gpx := GPX{
		Trk: Trk{
			Trkseg: Trkseg{
				Trkpt: GetTrackpoints(waypoints),
			},
		},
	}
	err = encoder.Encode(gpx)
	if err != nil {
		fmt.Println("Error encoding XML:", err)
		return
	}

	fmt.Println("Created file:", filename)
}

func GetTrackpoints(waypoints []Waypoint) []Trkpt {
	var trackpoints []Trkpt
	for _, wp := range waypoints {
		trackpoints = append(trackpoints, Trkpt{
			Lat: wp.Latitude,
			Lon: wp.Longitude,
		})
	}
	return trackpoints
}

type GPX struct {
	Trk Trk `xml:"trk"`
}

type Trk struct {
	Trkseg Trkseg `xml:"trkseg"`
}

type Trkseg struct {
	Trkpt []Trkpt `xml:"trkpt"`
}

type Trkpt struct {
	Time string  `xml:"time,omitempty"`
	Lat  float64 `xml:"lat,attr"`
	Lon  float64 `xml:"lon,attr"`
}

type Waypoint struct {
	Time      time.Time
	Latitude  float64
	Longitude float64
	Distance  int
}
