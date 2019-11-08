/**
* @Author: HongBo Fu
* @Date: 2019/11/4 14:09
 */

package utility

func InArray(in string, ins []string) bool {
	include := false
	for _, v := range ins {
		if in == v {
			include = true
		}
	}
	return include
}

func InMap(key string, ins map[interface{}]interface{}) bool {
	include := false
	for k, _ := range ins {
		if key == k.(string) {
			include = true
		}
	}
	return include
}
