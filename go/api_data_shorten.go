/*
 * URL Shortener API
 *
 * This is a URL Shortener API
 *
 * API version: 1.0.0
 * Contact: aurelien@duboc.xyz
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"log"
	"net/http"
	"net/url"
)

type CreateNewShortURL struct {
	db      *pg.DB
	baseUrl string
}

func FormatResponse(baseUrl string, shortURL string) []byte {
	response := make(map[string]string)
	shortURLResponse := fmt.Sprintf("%s/%s", baseUrl, shortURL)
	response["shortURL"] = shortURLResponse
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	return jsonResponse
}

func (h *CreateNewShortURL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var data longUrlPayload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = url.ParseRequestURI(data.LongURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURL := computeShortURL(data.LongURL)
	shortUrlExist, shortUrlMap := selectShortURL(h.db, shortURL)
	if !shortUrlExist {
		addShortUrl(h.db, data.LongURL, shortURL)
	} else {
		if shortUrlMap.LongURL != data.LongURL {
			// TODO: Manage hash collision
			fmt.Println("COLLISION")
		}
	}

	w.WriteHeader(http.StatusCreated)
	jsonResponse := FormatResponse(h.baseUrl, shortURL)
	w.Write(jsonResponse)
}
