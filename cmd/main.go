package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NoteInfo struct {
	Title    string `json:"title"`
	Context  string `json:"context"`
	Author   string `json:"author"`
	IsPublic bool   `json:"is_public"`
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		// handle returned error here.
		w.WriteHeader(503)
		w.Write([]byte("bad"))
	}
}

var (
	ErrnoData = errors.New("No Data")
)

type CreateResponse struct {
	UUID uuid.UUID `json:"uuid"`
	Data NoteInfo  `json:"data"`
}

type DataMap struct {
	mu   sync.RWMutex
	data map[uuid.UUID]NoteInfo
}

func newData() *DataMap {
	var dataMap = map[uuid.UUID]NoteInfo{}
	return &DataMap{
		data: dataMap,
	}
}

func (d *DataMap) Set(value NoteInfo) (uuid.UUID, NoteInfo, error) {
	uuid := uuid.New()
	fmt.Println(uuid)
	d.data[uuid] = value
	return uuid, d.data[uuid], nil
}

func (d *DataMap) Get(key uuid.UUID) (NoteInfo, error) {
	if _, ok := d.data[key]; !ok {
		return NoteInfo{}, ErrnoData
	}
	return d.data[key], nil
}

func main() {
	dataMap := newData()
	//dataMap.Set("AHAHA", "VALUUEEE")

	r := chi.NewRouter()
	r.Method("GET", "/api/note/{id}", noteGetHandler(dataMap))
	r.Method("POST", "/api/note", noteStoreHandler(dataMap))
	http.ListenAndServe(":3333", r)
}

func noteStoreHandler(d *DataMap) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {

		var note NoteInfo
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			return err
		}

		uuid, data, err := d.Set(note)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(CreateResponse{
			UUID: uuid,
			Data: data,
		})
	}

}

func noteGetHandler(d *DataMap) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {

		idStr := chi.URLParam(r, "id")

		// 2. Парсим в uuid.UUID
		uuid, err := uuid.Parse(idStr)
		if err != nil {
			return ErrnoData
		}

		data, err := d.Get(uuid)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return nil
		}

		//w.Write([]byte(json.Marshal(data)))
		//fmt.Println(d)
		//return nil

		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(data)
	}

}

/*
func customHandler(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query().Get("err")

	if q != "" {
		return errors.New(q)
	}

	w.Write([]byte("foo"))
	return nil
}
*/
