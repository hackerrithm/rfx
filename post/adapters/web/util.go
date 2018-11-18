package web

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
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

func getCookie(w http.ResponseWriter, req *http.Request) *http.Cookie {
	c, err := req.Cookie("session")
	if err != nil {
		sID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
	}
	return c
}

// takes in a file name now also
func appendValue(w http.ResponseWriter, c *http.Cookie, fname string) *http.Cookie {
	s := c.Value
	if !strings.Contains(s, fname) {
		s += "|" + fname
	}
	c.Value = s
	http.SetCookie(w, c)
	return c
}

// FileUpload ...
func FileUpload(w http.ResponseWriter, r *http.Request) (string, error) {
	// cookie := getCookie(w, r)

	mf, fh, err := r.FormFile("contentPhoto")
	if err != nil {
		fmt.Println(err)
	}
	defer mf.Close()
	// create sha for file name
	ext := strings.Split(fh.Filename, ".")[1]
	h := sha1.New()
	io.Copy(h, mf)
	fname := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
	// create new file
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path := filepath.Join(wd+"../../../", "assets", "photos", fname)
	nf, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	defer nf.Close()
	// copy
	mf.Seek(0, 0)
	io.Copy(nf, mf)
	// add filename to this user's cookie
	// cookie = appendValue(w, cookie, fname)
	// xs := strings.Split(cookie.Value, "|")
	fmt.Println(fname, " :: fname")

	return fname, nil
}
