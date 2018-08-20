package main

import (
	"encoding/json"
	"net/http"
	"github.com/bsm/openrtb"
	"github.com/labstack/echo"
)

type User struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/users", getUsers)
	e.GET("/users/:id", getUser)
	e.POST("/get_bid", getBid)
	e.Logger.Fatal(e.Start(":1323"))
}

func getUser(c echo.Context) error {
  	// User ID from path `users/:id`
  	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

/**
   Sample request:
   {
     "name": "aa",
     "email": "bb",
     "smth_ignored": "cc"
   }
 */
func getUsers (c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
	// or
	// return c.XML(http.StatusCreated, u)
}

func getBid (c echo.Context) error {
	var req *openrtb.BidRequest
	err = json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
	  return err
	}

	/**
	  * We are returning a bid request though in a production project we have
	  * to return a valid BidResponse
	  * @see RTB docs
	 */
	return c.JSON(http.StatusCreated, req)
}
