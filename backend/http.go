package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type attribute struct {
	ID    uuid.UUID `json:"-" db:"komp_id"`
	Name  string    `json:"name" db:"name"`
	Value string    `json:"value" db:"value"`
}

type komp struct {
	ID              uuid.UUID   `json:"-" db:"id"`
	SerialNumber    string      `json:"serial_number" db:"serial_number"`
	State           string      `json:"state" db:"state"`
	SoftwareVersion string      `json:"software_version" db:"software_version"`
	ProductCode     string      `json:"product_code" db:"product_code"`
	MacAddress      string      `json:"mac_address" db:"mac_address"`
	Comment         string      `json:"comment" db:"comment"`
	Attributes      []attribute `json:"attributes"`
}

func newHttpServer(addr string, conn *sqlx.DB) http.Server {
	router := &mux.Router{}
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	})
	router.HandleFunc("/api/komps", getKomps(conn)).Methods(http.MethodGet)
	router.HandleFunc("/api/komps/{serial_number}", getKomp(conn)).Methods(http.MethodGet)
	router.HandleFunc("/api/komps/{serial_number}", patchKomp(conn)).Methods(http.MethodPatch)

	return http.Server{Addr: addr, Handler: router}
}

func getKomps(conn *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		komps, err := findAllKomps(r.Context(), conn)
		if err != nil {
			writeTextResponse(http.StatusInternalServerError, err.Error(), w)
			return
		}

		writeJSONResponse(http.StatusOK, komps, w)
	}
}

func getKomp(conn *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serialNumber, err := parseSerialNumberFromRequest(r)
		if err != nil {
			writeTextResponse(http.StatusBadRequest, err.Error(), w)
			return
		}

		komp, err := findKompWithSerialNumber(r.Context(), serialNumber, conn)
		if err == sql.ErrNoRows {
			writeTextResponse(http.StatusNotFound, fmt.Sprintf("found no Komp with serial number %s", serialNumber), w)
			return
		}

		if err != nil {
			writeTextResponse(http.StatusInternalServerError, err.Error(), w)
			return
		}

		writeJSONResponse(200, komp, w)
	}
}

type kompPatch struct {
	State   *string `json:"state"`
	Comment *string `json:"comment"`
}

func patchKomp(conn *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serialNumber, err := parseSerialNumberFromRequest(r)
		if err != nil {
			writeTextResponse(http.StatusBadRequest, err.Error(), w)
			return
		}

		var patch kompPatch
		err = json.NewDecoder(r.Body).Decode(&patch)
		if err != nil {
			writeTextResponse(http.StatusBadRequest, "failed to decode request body", w)
			return
		}

		var (
			newState, newComment string
		)

		if patch.State != nil {
			newState = strings.ToLower(*patch.State)
			if newState != "available" && newState != "allocated" {
				writeTextResponse(http.StatusNotFound, fmt.Sprintf("invalid state \"%s\", accepted values are \"available\" and \"allocated\"", newState), w)
				return
			}
		}

		if patch.Comment != nil {
			newComment = *patch.Comment
		}

		komp, err := updateKompStateAndComment(r.Context(), serialNumber, newState, newComment, conn)
		if err == sql.ErrNoRows {
			writeTextResponse(http.StatusNotFound, fmt.Sprintf("found no Komp with serial number %s", serialNumber), w)
			return
		}

		if err != nil {
			writeTextResponse(http.StatusInternalServerError, err.Error(), w)
			return
		}

		writeJSONResponse(200, komp, w)
	}
}

var (
	isAlphaNumeric = regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
)

func parseSerialNumberFromRequest(r *http.Request) (string, error) {
	serialNumber, ok := mux.Vars(r)["serial_number"]
	if !ok {
		return "", errors.New("serial number is required")
	}

	if len(serialNumber) != 8 {
		return "", errors.New("serial number must be 8 characters")
	}

	if !isAlphaNumeric(serialNumber) {
		return "", errors.New("serial number must be alphanumeric")
	}

	return serialNumber, nil
}

func writeTextResponse(code int, response string, w http.ResponseWriter) {
	w.WriteHeader(code)
	w.Header().Add("content-type", "text/plain")
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Printf("failed to encode json response: %s", err)
	}
}

func writeJSONResponse(code int, response interface{}, w http.ResponseWriter) {
	w.WriteHeader(code)
	w.Header().Add("content-type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("failed to encode json response: %s", err)
	}
}
