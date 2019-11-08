/**
* @Author: HongBo Fu
* @Date: 2019/10/23 14:21
 */

package utility

import "strings"

func SemanticDetection(value string) (t, v string) {
	if strings.Contains(value, "@run") {
		t = "js"
		v = strings.ReplaceAll(strings.ReplaceAll(value, "\n", ""), "@run", "")
	} else if strings.Contains(value, "@bind") {
		t = "bind"
		v = strings.ReplaceAll(strings.ReplaceAll(value, "\n", ""), "@bind", "")
	} else {
		t = "fix"
		v = strings.ReplaceAll(value, "\n", "")
	}

	return
}
