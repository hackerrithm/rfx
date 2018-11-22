package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type response struct {
	Result interface{} `json:"result"`
}

// decodeReq decodes request's body to given interface
func decodeReq(r *http.Request, to interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(to); err != nil {
		if err != io.EOF {
			return err
		}
	}
	return nil
}

func queryValue(k string, r *http.Request) string {
	values := r.URL.Query()[k]

	if len(values) != 0 {
		return values[0]
	}

	return ""
}

func queryValueInt(k string, r *http.Request) (int, error) {
	qv := queryValue(k, r)
	if qv == "" {
		return 0, nil
	}
	return strconv.Atoi(qv)
}

func urlParamMustInt(k string, r *http.Request) int {
	i, err := strconv.Atoi(mux.Vars(r)[k])
	if err != nil {
		panic(fmt.Sprintf("url param can't convert to int: %v", err))
	}
	return i
}

func urlParamMust(k string, r *http.Request) string {
	v, ok := mux.Vars(r)[k]
	if !ok || v == "" {
		panic(fmt.Sprintf("there isn't url param with this key: %s", k))
	}
	return v
}

// JWTData is a struct with the structure of the jwt data
// type JWTData struct {
// 	// Standard claims are the standard jwt claims from the IETF standard
// 	// https://tools.ietf.org/html/rfc7519
// 	jwt.StandardClaims
// 	CustomClaims map[string]string `json:"custom,omitempty"`
// }
