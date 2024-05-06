package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	gpxInputFile := os.Args[1]
	gpxOutputFileName := os.Args[2]
	gpxOutputDir := os.Args[3]

	// Extract waypoints from source GPX file
	waypoints := GetWaypoints(gpxInputFile)

	// Split waypoints into 500-waypoint chunks
	chunks := SplitWaypoints(waypoints, 500)

	for i, chunk := range chunks {
		filename := fmt.Sprintf(
			gpxOutputDir+"/"+"%d_"+gpxOutputFileName+".gpx",
			i,
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

	for _, point := range gpxData.Trk.Trkseg.Trkpt {
		waypoints = append(
			waypoints,
			Waypoint{
				Latitude:  point.Lat,
				Longitude: point.Lon,
			},
		)
	}

	return waypoints
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
			return
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	gpx := GPX{
		Trk: Trk{
			Trkseg: Trkseg{
				Trkpt: GetTrackpoints(waypoints),
			},
		},
	}

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	err = encoder.Encode(gpx)
	if err != nil {
		return
	}
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
	Latitude  float64
	Longitude float64
}
