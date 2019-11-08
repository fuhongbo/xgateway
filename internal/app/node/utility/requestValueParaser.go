/**
* @Author: HongBo Fu
* @Date: 2019/10/24 13:52
 */

package utility

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

func ReqeustParaser(bind string, r *http.Request) (string, bool) {
	lep := strings.Split(bind, ".")

	if lep[0] != "request" {
		return Error_NotSupport_In_RequestParaser01, false
	}
	if lep[0] == "request" {
		request := reflect.ValueOf(r)
		requestValue := request.Elem()

		switch lep[1] {
		case "RequestURI", "Host", "RemoteAddr":
			return requestValue.FieldByName(lep[1]).Interface().(string), true
		case "Header":
			return requestValue.FieldByName(lep[1]).Interface().(http.Header).Get(lep[2]), true
		case "Form", "PostForm":
			_ = r.ParseForm()
			return requestValue.FieldByName(lep[1]).Interface().(url.Values).Get(lep[2]), true
		case "Body":
			rbody, _ := ioutil.ReadAll(r.Body)
			return string(rbody), true
		}

		return Error_NotSupport_In_RequestParaser02, false

	}

	return Error_NotSupport_In_RequestParaser02, false

}
