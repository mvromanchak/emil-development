package entities

type GPSData struct {
	DeviceId string
	GPS      []GPS
}

type GPS struct {
	Lat  string
	Lon  string
	Ele  string
	Time string
}
