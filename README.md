# polar-route-import

Polar Vantage V is limited to [500 waypoints](https://support.polar.com/sg-en/how-to-import-route):

> If the route is long, it will be compressed into 500 waypoints when it is transferred to Grit X/Grit X Pro/Pacer Pro/V650/Vantage/Vantage V2

Let's split long GPX route into 500-waypoint chunks. This is especially useful when you are running ultras.

## Run

```
go run main.go \
~/Downloads/SG200Miles2024.gpx \
"SG200" \
~/Downloads/SG200Miles2024 
```