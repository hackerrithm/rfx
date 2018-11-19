package web

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
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

const maxUploadSize = 20000 * 1024
const uploadPath = ".../../../assets"

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int, file multipart.File) string {
	h := sha1.New()
	io.Copy(h, file)
	b := make([]byte, len)
	rand.Read(b)
	ve := h.Sum(nil)
	val := append(ve[:], b[:]...)
	return fmt.Sprintf("%x", val)
}

// FileUpload ...
func FileUpload(w http.ResponseWriter, r *http.Request) (string, error) {

	// validate file size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
		return "", nil
	}

	// parse and validate file and post parameters
	fileType := r.PostFormValue("type")
	file, fh, err := r.FormFile("contentPhoto")
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return "", nil
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return "", nil
	}

	ext := strings.Split(fh.Filename, ".")[1]

	// check file type, detectcontenttype only needs the first 512 bytes
	filetype := http.DetectContentType(fileBytes)
	fmt.Println("filetype: ", ext)
	switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png":
	case "application/pdf":
		break
	default:
		renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		return "", nil
	}
	fileName := randToken(12, file)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fileName = fileName + "." + ext
	newPath := filepath.Join(wd+uploadPath, "photos", fileName)
	fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)

	// write file
	newFile, err := os.Create(newPath)
	if err != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return "", nil
	}
	defer newFile.Close() // idempotent, okay to call twice
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return "", nil
	}
	w.Write([]byte("SUCCESS"))

	return fileName, nil
}
