package main

import (
	"encoding/json"
	"github.com/bsm/openrtb"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"xojoc.pw/useragent"
)

var GEO_SERVICE = "https://ipinfo.io/"

type Telemetry struct {
	OS		string	`json:"os,omitempty"`
	Device	string	`json:"device,omitempty"`
	Client	string	`json:"client,omitempty"`
	State	string	`json:"state,omitempty"`
	Domain	string	`json:"domain,omitempty"`
}

type GeoInfo struct {
	IP  string `json:"ip" xml:"ip" form:"ip" query:"ip"`
	Hostname  string `json:"hostname" xml:"hostname" form:"hostname" query:"hostname"`
	City  string `json:"city" xml:"city" form:"city" query:"city"`
	Region  string `json:"region" xml:"region" form:"region" query:"region"`
	Country  string `json:"country" xml:"country" form:"country" query:"country"`
	Location  string `json:"loc" xml:"loc" form:"loc" query:"loc"`
	Postal  string `json:"postal" xml:"postal" form:"postal" query:"postal"`
	Org  string `json:"org" xml:"org" form:"org" query:"org"`
}

func main() {
	e := echo.New()
	e.POST("/", getBidAction)
	e.Logger.Fatal(e.Start(":1323"))
}

func getBidAction (c echo.Context) error {
	var req *openrtb.BidRequest
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
	  return err
	}

	t := new(Telemetry)

	fillFromUA(t, req)
	fillDomain(t, req)
	fillState(t, req)

	/**
	  * We are returning a Telemetry struct JSON-encoded response
	  * though in a production project we have to return a valid BidResponse
	  * @see RTB docs
	 */
	return c.JSON(http.StatusCreated, t)
}

func fillDomain(t *Telemetry, req *openrtb.BidRequest) {
	if req.Site == nil {
		return
	}
	if len(req.Site.Domain) <= 0 {
		t.Domain = getDomainFromUrl(&req.Site.Page)
		return
	}
	t.Domain = req.Site.Domain
}

func getDomainFromUrl(urltext *string) string {
	u, err := url.Parse(*urltext)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func fillFromUA(t *Telemetry, req *openrtb.BidRequest) {
	ua := getUserAgent(req)
	if ua == nil {
		return
	}
	fillOperatingSystem(t, ua)
	fillClient(t, ua)
	fillClientType(t, ua)
}

func getUserAgent(req *openrtb.BidRequest) *useragent.UserAgent {
	return useragent.Parse(req.Device.UA)
}

func fillOperatingSystem(t *Telemetry, ua *useragent.UserAgent) {
	t.OS = ua.OS;
}

func fillClient(t *Telemetry, ua *useragent.UserAgent) {
	t.Client = ua.Name;
}

func fillClientType(t *Telemetry, ua *useragent.UserAgent) {
    if ua.Mobile {
        t.Device = "Mobile"
    } else if ua.Tablet {
        t.Device = "Tablet"
    } else {
        t.Device = "Desktop"
    }
}

func fillState(t *Telemetry, req *openrtb.BidRequest) {
	geoInfo := getGeoIpInfo(req)
	if geoInfo != nil {
		t.State = geoInfo.Country
	}
}

func getGeoIpInfo(req *openrtb.BidRequest) *GeoInfo {
	serviceUrl := getGeoServiceUrl(req.Device)
	if len(serviceUrl) <= 0 {
		return nil
	}

	body := getGeoResponse(&serviceUrl)
	if body == nil {
		return nil
	}

	geoInfo := parseGeoInfo(body)
	if geoInfo == nil {
		return nil
	}

	return geoInfo
}

func getGeoResponse(url *string) []byte {
	resp, err := http.Get(*url)
	if err != nil {
		log.Printf("%+v\n", "ERROR REQUESTING GEOIP INFO")
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}

func parseGeoInfo(body []byte) *GeoInfo {
	geoInfo := new (GeoInfo)
	if err := json.Unmarshal(body, &geoInfo); err != nil {
		log.Printf("%+v\n", "ERROR PARSING GEOIP INFO")
		return nil
	}
	return geoInfo
}

func getGeoServiceUrl(d *openrtb.Device) string {
	switch {
		case len(d.IP) > 0:
			return getGeoServiceBaseUrl(&d.IP)
		case len(d.IPv6) > 0:
			return getGeoServiceBaseUrl(&d.IPv6)
		default:
			return ""
	}
}

func getGeoServiceBaseUrl(IP *string) string {
	return GEO_SERVICE + *IP
}
