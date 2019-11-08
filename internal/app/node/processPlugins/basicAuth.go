/**
* @Author: HongBo Fu
* @Date: 2019/10/23 08:29
 */

package processPlugins

import (
	"encoding/base64"
	"github.com/robertkrimen/otto"
	"net/http"
	"strings"
)

type BasicAuth struct {
	UserName string
	Password string
}

func (p *BasicAuth) Exec(w http.ResponseWriter, r *http.Request, vm *otto.Otto) bool {

	if !checkAuth(w, r, p.UserName, p.Password) {
		w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
		return false
	}

	return true

}

func checkAuth(w http.ResponseWriter, r *http.Request, userName, password string) bool {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}

	return pair[0] == userName && pair[1] == password
}
