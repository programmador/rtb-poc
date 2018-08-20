package main

import (
	//"regexp"
	//"log"
	"net/url"
	"encoding/json"
	"net/http"
	"github.com/bsm/openrtb"
	"github.com/labstack/echo"
)

type Telemetry struct {
	OS		string	`json:"os,omitempty"`
	Device	string	`json:"device,omitempty"`
	Client	string	`json:"client,omitempty"`
	State	string	`json:"state,omitempty"`
	Domain	string	`json:"domain,omitempty"`
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
	fillDomain(t, req)
	t.State = "UA"

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
	/*re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)
	submatchall := re.FindAllString(*url,-1)
	for _, element := range submatchall {
		log.Printf("%+v\n", element)
		//return element
	}*/
	u, err := url.Parse(*urltext)
	if err != nil {
		return ""
	}
	return u.Hostname()
}
