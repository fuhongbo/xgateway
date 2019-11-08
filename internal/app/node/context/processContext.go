/**
* @Author: HongBo Fu
* @Date: 2019/10/23 13:26
 */

package context

import "xgateway/internal/pkg/syncmap"

type ProcessContext struct {
	Variable syncmap.Map
	Header   syncmap.Map
}
