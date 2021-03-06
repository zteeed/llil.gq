package app

import (
	"encoding/json"
	"fmt"
	"llil.gq/go/database"
	"log"
	"net/http"
	"net/url"
)

type DataShortenHandler struct {
	database database.Database
	baseUrl  string
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

func (h *DataShortenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
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
	shortURL := h.generateShortUrl(data)

	w.WriteHeader(http.StatusCreated)
	jsonResponse := FormatResponse(h.baseUrl, shortURL)
	w.Write(jsonResponse)
}

func (h *DataShortenHandler) generateShortUrl(data longUrlPayload) string {
	shortURL := computeShortURL(data.LongURL)
	result, err := database.SelectShortURL(h.database, shortURL)
	if err != nil {
		database.AddShortUrl(h.database, data.LongURL, shortURL)
	} else {
		if result.LongURL != data.LongURL {
			for err == nil {
				shortURL = computeShortURL(h.baseUrl + data.LongURL)
				result, err = database.SelectShortURL(h.database, shortURL)
			}
			database.AddShortUrl(h.database, data.LongURL, shortURL)
		}
	}
	return shortURL
}
