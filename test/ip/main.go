/**
* @Author: HongBo Fu
* @Date: 2019/10/24 15:30
 */

package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(isBelong(`10.187.102.222`, `10.187.102.223-10.187.102.221`))

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
