package admin

import (
	"net/http"
)

func AdminSignup(rw http.ResponseWriter, r *http.Request) {

}

// func AdminSignin(rw http.ResponseWriter, r *http.Request) {

// 	token := jwt.New(jwt.GetSigningMethod("HS256"))
// 	claims := make(jwt.MapClaims)
// 	claims["userName"] = loginRequest.UserName
// 	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
// 	token.Claims = claims
// 	tokenString, err := token.SignedString([]byte(SecretKey))
// 	tokenByte, err := json.Marshal(data)
// 	w.WriteHeader(201)
// 	w.Write(tokenByte)
// }
