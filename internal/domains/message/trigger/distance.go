package trigger

import "math"

const EarthRadiusMeters = 6371000 // в метрах

func DistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {

	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)

	lat1Rad := toRad(lat1)
	lat2Rad := toRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadiusMeters * c
}

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}
