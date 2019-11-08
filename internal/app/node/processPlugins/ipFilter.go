/**
* @Author: HongBo Fu
* @Date: 2019/10/24 15:49
 */

package processPlugins

import (
	"github.com/robertkrimen/otto"
	"net/http"
	"strconv"
	"strings"
	"xgateway/internal/app/node/utility"
)

type IPFilter struct {
	Action  string
	Address string
}

func (p *IPFilter) Exec(w http.ResponseWriter, r *http.Request, vm *otto.Otto) bool {

	remoteHost := strings.Split(r.RemoteAddr, ":")

	if len(remoteHost) > 2 || remoteHost[0] == "127.0.0.1" {
		return true
	}

	if isBelong(remoteHost[0], p.Address) {
		if p.Action == "allow" {
			return true
		}
	}

	w.WriteHeader(510)
	w.Write([]byte(utility.Error_IPBlock))
	return false
}

func isBelong(ip, cidr string) bool {
	ipAddr := strings.Split(ip, `.`)
	if len(ipAddr) < 4 {
		return false
	}
	cidrArr := strings.Split(cidr, `/`)
	if len(cidrArr) < 2 {
		if strings.Contains(cidr, "-") {

			ips := strings.Split(cidr, "-")
			result := false
			for _, i := range ips {
				if i == ip {
					result = true
					break
				}
			}
			return result

		} else {
			if ip == cidr {
				return true
			} else {
				return false
			}
		}

	}
	var tmp = make([]string, 0)
	for key, value := range strings.Split(`255.255.255.0`, `.`) {
		iint, _ := strconv.Atoi(value)

		iint2, _ := strconv.Atoi(ipAddr[key])

		tmp = append(tmp, strconv.Itoa(iint&iint2))
	}
	return strings.Join(tmp, `.`) == cidrArr[0]
}
