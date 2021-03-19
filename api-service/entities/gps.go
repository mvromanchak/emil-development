package entities

import (
	"encoding/xml"
	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/mvromanchak/emil-development/api-service/sdk"
)

type GPSData struct {
	DeviceID string
	Trkpt    []GPSTrack
}
type GPSTrack struct {
	Lat  string
	Lon  string
	Ele  string
	Time string
}

type GPSDataRequest struct {
	sdk.JWTSecured
	DeviceID string
	Trkpt    []GPSTrack
}

// Authorize check jwt claims.
func (r GPSDataRequest) Authorize(deviceID string, claims jwt.MapClaims) error {
	if claims["devise_id"] == "" {
		return errors.New("missing devise_uuid")
	}
	if deviceID != claims["devise_id"] {
		return errors.New("auth error: devise_id")
	}
	return nil
}

type GPXRequestBody struct {
	XMLName        xml.Name `xml:"gpx"`
	Text           string   `xml:",chardata"`
	Xmlns          string   `xml:"xmlns,attr"`
	Gpxx           string   `xml:"gpxx,attr"`
	Gpxtpx         string   `xml:"gpxtpx,attr"`
	Creator        string   `xml:"creator,attr"`
	Version        string   `xml:"version,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Metadata       struct {
		Text string `xml:",chardata"`
		Link struct {
			Chardata string `xml:",chardata"`
			Href     string `xml:"href,attr"`
			Text     string `xml:"text"`
		} `xml:"link"`
		Time string `xml:"time"`
	} `xml:"metadata"`
	Trk struct {
		Text   string `xml:",chardata"`
		Name   string `xml:"name"`
		Trkseg struct {
			Text  string `xml:",chardata"`
			Trkpt []struct {
				Text string `xml:",chardata"`
				Lat  string `xml:"lat,attr"`
				Lon  string `xml:"lon,attr"`
				Ele  string `xml:"ele"`
				Time string `xml:"time"`
			} `xml:"trkpt"`
		} `xml:"trkseg"`
	} `xml:"trk"`
}

type GPXResponse struct {
	OK bool
}
