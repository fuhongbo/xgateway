/**
* @Author: HongBo Fu
* @Date: 2019/11/1 13:51
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type JWTClaims struct { // token里面添加用户信息，验证token后可能会用到用户信息
	jwt.StandardClaims
	UserID      uint     `json:"user_id"`
	Username    string   `json:"username"`
	FullName    string   `json:"full_name"`
	Permissions []string `json:"permissions"`
}

func main() {
	//hmacSampleSecret := []byte("use")

	//https://login.partner.microsoftonline.cn/c86a7f25-a468-49aa-9e5c-d0cd1d3a3eb8/v2.0/.well-known/openid-configuration
	// sample token string taken from the New example
	tokenString := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IndSZDZ5UmVBSmg5bWJocEtKT1BNdEE3NnhwUSIsImtpZCI6IndSZDZ5UmVBSmg5bWJocEtKT1BNdEE3NnhwUSJ9.eyJhdWQiOiI5NTIyNjA3Mi02MDU3LTQ1ZjAtOGNjNy1hODY3NjU2ZTc2ZDEiLCJpc3MiOiJodHRwczovL3N0cy5jaGluYWNsb3VkYXBpLmNuLzI5YjY5ZDI5LWVhNmMtNDNlNi04Y2RhLWVkYmMwYTViYmVlMC8iLCJpYXQiOjE1NzI1OTc3MjMsIm5iZiI6MTU3MjU5NzcyMywiZXhwIjoxNTcyNjAxNjIzLCJhaW8iOiJZMlZnWU1pZEZuNndYbXduNjlZTkJZSDF0OFRLQUE9PSIsImFwcGlkIjoiOTUyMjYwNzItNjA1Ny00NWYwLThjYzctYTg2NzY1NmU3NmQxIiwiYXBwaWRhY3IiOiIxIiwiaWRwIjoiaHR0cHM6Ly9zdHMuY2hpbmFjbG91ZGFwaS5jbi8yOWI2OWQyOS1lYTZjLTQzZTYtOGNkYS1lZGJjMGE1YmJlZTAvIiwib2lkIjoiNGI0YzA1MzUtMmQ1NS00MmI4LTkzNGUtOTc0MjIzOWFiMDE1Iiwic3ViIjoiNGI0YzA1MzUtMmQ1NS00MmI4LTkzNGUtOTc0MjIzOWFiMDE1IiwidGlkIjoiMjliNjlkMjktZWE2Yy00M2U2LThjZGEtZWRiYzBhNWJiZWUwIiwidXRpIjoiTjYtLUVtRmNIa0tqT1dSQmhlOFVBQSIsInZlciI6IjEuMCJ9.n-LT7HYBjVezZZGpLVeDTuppL1qH8V-BktNF5xWPv0reGawCxLrKJ_3kVB2eypu7PmfgzdS8CtnJqsAKsYaXovxirGXsjob-ENgwfgOBe3SfhZDuqEGPpu5rzaWKISlKdR0-e2EZVltVRKyMKWIwDWnHG48oNL9YOF1FbHNkYD5_sd68XyemOEfPu88nNJlI42qNEpeN3Lo-XGi2lrkvHv59ocYNLFJPg_84y2NYQr4wyOq9hSCPAcrOvEP3FgWe14duizKGdXblNjOAP1ijNPpI76Q90DueQUuCirC2QkcmzdD5lQ4r_1-u5wO5ql_Hn2TVWskhnth-9fKdYBjeLg"

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		//这里主要是检查是否是规定的加密算法
		if token.Method.Alg() != "RS256" {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		k5c := "MIIDDzCCAfegAwIBAgIQT65UBoMqsYROZFcZWgTZADANBgkqhkiG9w0BAQsFADAyMTAwLgYDVQQDEydhY2NvdW50cy5hY2Nlc3Njb250cm9sLmNoaW5hY2xvdWRhcGkuY24wHhcNMTkwOTE1MDAwMDAwWhcNMjQwOTE0MDAwMDAwWjAyMTAwLgYDVQQDEydhY2NvdW50cy5hY2Nlc3Njb250cm9sLmNoaW5hY2xvdWRhcGkuY24wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC670YLZ193tU0+MbF46n7aTGg8zxCjfww3aUow36MYgk76snCBLwEn7K61cALCVry8G1Cy2ijhvg/yEx0QElPccOFv3C/vY72m+s52BIEZcpq+PHIhy0mt0OO28wPZ2Q9BkQmPKtJo7YuFnDXO5LWuTuVDu7XnjsCsB3EHqeJXGmQ1F9IleVZ4SMio/J0+rErqiaQlIig3qro11/FIcku0jLZ27fO3N3ZPEOT7cp3m93Z20ssqDdmmtJRBCtE/aNfVTZOp44UGA10Xs3sp2/cBfQMseb1LItgnM4gKDbjx5GsMJ55hzm2yDDZ05wY5eleyfIBwf+NJq3J9nM2laZ5/AgMBAAGjITAfMB0GA1UdDgQWBBS9WHjN/s7ZlUocwvxD0lJ+elTo5DANBgkqhkiG9w0BAQsFAAOCAQEAup6F+snhwAqr8HnPp0B4SA8pyH63Veh5cca8zQ3z11EwooG8/ChGBJEfHALQzsxfHBemvWHbGKZPZygBQb0iyEFe4z7+pIFHzYMA2VaBLsV9Fh/fKdrnTCFcDtq4vByJdm1Mc13RIN9Sb63qT0BbgFF9hgJyOrFf2SrAhO7IdFPVXp8EFZimzLLC0fZwT50hQje9zJ0xMpYT+yd/Rwa67BJ64lXByDDKyJ1ahgHT2cz9UUThE/Joj20E5+RChgtYpuIhHzp0hyvvBqjD9byzEwMXIY+0dsSReuLtmYvcXeaQ0fBfEsqdNqyebOZ2QDbOPVu+daRLDjjKBNAR/mnA6A=="
		cert := "-----BEGIN CERTIFICATE-----\n" + k5c + "\n-----END CERTIFICATE-----"

		var result interface{}
		//if strings.Contains(strings.ToLower(token.Method.Alg()), "rs") {
		result, _ = jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		//}
		//
		//if strings.Contains(strings.ToLower(token.Method.Alg()), "ec") {
		//	result, _ = jwt.ParseECPublicKeyFromPEM([]byte(cert))
		//}
		nN, e := json.Marshal(token.Claims)
		if e != nil {
			println(e.Error())
		}
		println(string(nN))
		return result, nil
	})

	if err != nil {
		println("err1", err.Error())
	}
	if !token.Valid {
		println("invalid")
		return
	}

	x := token.Claims.Valid()

	if x != nil {
		println("invalid")
		return
	}

	println("valid")

}
