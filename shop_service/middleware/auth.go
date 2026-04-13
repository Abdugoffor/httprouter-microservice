package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"shop_service/appctx"
	"shop_service/helper"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

// var jwtKey = []byte(helper.ENV("JWT_KEY"))

// type Claims struct {
// 	UserID int    `json:"user_id"`
// 	Role   string `json:"role"`
// 	jwt.RegisteredClaims
// }

// func ParseToken(tokenStr string) (*Claims, error) {
// 	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	claims := token.Claims.(*Claims)
// 	return claims, nil
// }

// func AdminOnly(next httprouter.Handle) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

// 		auth := r.Header.Get("Authorization")
// 		if auth == "" {
// 			http.Error(w, "no token", 401)
// 			return
// 		}

// 		tokenStr := strings.Replace(auth, "Bearer ", "", 1)

// 		claims, err := ParseToken(tokenStr)
// 		if err != nil {
// 			http.Error(w, "invalid token", 401)
// 			return
// 		}

// 		if claims.Role != "admin" {
// 			http.Error(w, "forbidden: admin only", 403)
// 			return
// 		}

// 		next(w, r, ps)
// 	}
// }

// func AuthRequired(next httprouter.Handle) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

// 		auth := r.Header.Get("Authorization")
// 		if auth == "" {
// 			http.Error(w, "token required", 401)
// 			return
// 		}

// 		tokenStr := strings.Replace(auth, "Bearer ", "", 1)

// 		claims, err := ParseToken(tokenStr)
// 		if err != nil {
// 			http.Error(w, "invalid token", 401)
// 			return
// 		}

// 		// contextga user saqlash (optional lekin juda muhim)
// 		r.Header.Set("user_id", fmt.Sprint(claims.UserID))
// 		r.Header.Set("role", claims.Role)

// 		next(w, r, ps)
// 	}
// }

var jwtKey = []byte(helper.ENV("JWT_KEY"))

type Claims struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
	jwt.RegisteredClaims
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(*Claims), nil
}

func AuthRequired(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "token required", 401)
			return
		}

		tokenStr := strings.Replace(auth, "Bearer ", "", 1)

		claims, err := ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", 401)
			return
		}

		// route pattern ni context dan olamiz: "/products/:id" (loop yo'q, O(1))
		routePath, _ := r.Context().Value(appctx.RoutePatternKey).(string)

		// 🔥 auth_service ga request
		body := map[string]any{
			"user_id": claims.UserID,
			"role_id": claims.RoleID,
			"method":  r.Method,
			"path":    routePath,
		}

		jsonBody, _ := json.Marshal(body)

		resp, err := http.Post(
			"http://localhost:8080/internal/check",
			"application/json",
			bytes.NewBuffer(jsonBody),
		)

		if err != nil {
			http.Error(w, "auth service error", 500)
			return
		}

		defer resp.Body.Close()

		var result map[string]bool
		json.NewDecoder(resp.Body).Decode(&result)

		if !result["allowed"] {
			http.Error(w, "forbidden", 403)
			return
		}

		// optional context
		r.Header.Set("user_id", string(rune(claims.UserID)))

		next(w, r, ps)
	}
}
