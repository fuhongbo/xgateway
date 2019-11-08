/**
* @Author: HongBo Fu
* @Date: 2019/10/28 09:59
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	//e.Debug = true
	//e.Use(middleware.Logger())

	e.POST("/test", func(c echo.Context) error {
		s, _ := ioutil.ReadAll(c.Request().Body) //把  body 内容读入字符串 s

		//println(string(s))
		//
		//println(c.Request().Header.Get("Authorization"))
		//println(c.Request().Header.Get("X-Service-ID"))
		var x interface{}
		json.Unmarshal(s, &x)
		return c.JSON(200, x)
	})

	e.Logger.Fatal(e.Start(":" + os.Args[1]))
}
