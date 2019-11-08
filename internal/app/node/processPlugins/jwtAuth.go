/**
* @Author: HongBo Fu
* @Date: 2019/11/4 13:09
 */

package processPlugins

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/robertkrimen/otto"
	"net/http"
	"strings"
	"xgateway/internal/app/node/utility"
)

type JwtAuth struct {
	Type      string
	Method    string
	PublicKey map[interface{}]interface{}
	Kid       string
	Verify    map[interface{}]interface{}
}

func (p *JwtAuth) Exec(w http.ResponseWriter, r *http.Request, vm *otto.Otto) bool {

	tok := r.Header.Get("Authorization")
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		tok = tok[7:]
	} else {
		w.WriteHeader(510)
		w.Write([]byte(utility.Error_Empty_AccessToekn))
		return false
	}

	token, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {

		if token.Method.Alg() != p.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		pkey := ""
		if len(p.PublicKey) == 1 {
			pkey = p.PublicKey["un"].(string)
		} else {
			if utility.InMap(p.Kid, p.PublicKey) {
				pkey = p.PublicKey[token.Header["kid"].(string)].(string)
			} else {
				return nil, fmt.Errorf("Unexpected token format")
			}

		}

		switch p.Type {
		case "ECDS":
			result, _ := jwt.ParseECPublicKeyFromPEM([]byte(pkey))
			return result, nil
		case "RS":
			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(pkey))
			return result, nil
		case "HMAC":
			return pkey, nil

		}

		mp := token.Claims.(jwt.MapClaims)
		for k, v := range p.Verify {
			vV := mp[k.(string)].(string)
			if !utility.InArray(vV, v.([]string)) {
				return nil, errors.New(k.(string) + "验证不匹配")
				break
			}
		}

		return nil, errors.New("加密方法暂不支持")

	})

	if err != nil {
		w.WriteHeader(510)
		w.Write([]byte(utility.Error_Jwt_Auth + err.Error()))
		return false
	}

	if !token.Valid {
		w.WriteHeader(510)
		w.Write([]byte(utility.Error_Jwt_Auth))
		return false
	}

	vm.Set("jwt", token)

	return true
}
