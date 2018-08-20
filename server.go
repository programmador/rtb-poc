package main

import (
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
	// @TODO 1) if "site" exists 2) if "domain not exists" - parse "page"
	t.Domain = req.Site.Domain
	t.State = "UA"

	/**
	  * We are returning a Telemetry struct JSON-encoded response
	  * though in a production project we have to return a valid BidResponse
	  * @see RTB docs
	 */
	return c.JSON(http.StatusCreated, t)
}
