/**
* @Author: HongBo Fu
* @Date: 2019/10/16 08:38
 */

package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	//e.Debug = true
	//e.Use(middleware.Logger())

	e.GET("/test1/hi", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!"+os.Args[1])
	})

	e.POST("test1/haha", func(c echo.Context) error {
		s, _ := ioutil.ReadAll(c.Request().Body)
		//var x interface{}
		//json.Unmarshal(s, &x)
		return c.String(http.StatusOK, string(s))
	})

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})

	e.Logger.Fatal(e.Start(":" + os.Args[1]))
}
