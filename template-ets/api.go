package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type sumberDayaInput struct {
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type peminjamanInput struct {
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
}

type errorBody struct {
	Error string `json:"error"`
}

type API struct {
	store *Store
}

func NewAPI(store *Store) *API {
	return &API{store: store}
}

func NewHandler() http.Handler {
	return NewAPI(NewStore()).Routes()
}

func (a *API) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /sumber-daya", a.createSumberDaya)
	mux.HandleFunc("GET /sumber-daya", a.listSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", a.getSumberDaya)
	mux.HandleFunc("POST /peminjaman", a.createPeminjaman)
	mux.HandleFunc("GET /peminjaman", a.listPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", a.getPeminjaman)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("gagal encode JSON: %v", err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrSumberDayaTidakDitemukan), errors.Is(err, ErrPeminjamanTidakDitemukan):
		writeJSON(w, http.StatusNotFound, errorBody{Error: err.Error()})
	case errors.Is(err, ErrInputTidakValid):
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
	default:
		// Detail asli hanya ke log server, klien menerima pesan generik.
		log.Printf("error internal: %v", err)
		writeJSON(w, http.StatusInternalServerError, errorBody{Error: "terjadi kesalahan internal"})
	}
}

// POST /sumber-daya
func (a *API) createSumberDaya(w http.ResponseWriter, r *http.Request) {
	var in sumberDayaInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "body bukan JSON yang valid"})
		return
	}
	in.Nama = strings.TrimSpace(in.Nama)
	in.Jenis = strings.TrimSpace(in.Jenis)
	if err := validasiSumberDaya(in); err != nil {
		writeError(w, err)
		return
	}
	sd := a.store.AddSumberDaya(SumberDaya{Nama: in.Nama, Jenis: in.Jenis, Stok: in.Stok})
	writeJSON(w, http.StatusCreated, sd)
}

// GET /sumber-daya
func (a *API) listSumberDaya(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.ListSumberDaya())
}

// GET /sumber-daya/{id}
func (a *API) getSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "id harus berupa angka"})
		return
	}
	sd, err := a.store.GetSumberDaya(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sd)
}

// POST /peminjaman
func (a *API) createPeminjaman(w http.ResponseWriter, r *http.Request) {
	var in peminjamanInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "body bukan JSON yang valid"})
		return
	}
	in.Peminjam = strings.TrimSpace(in.Peminjam)
	in.Tanggal = strings.TrimSpace(in.Tanggal)
	if err := validasiPeminjaman(in); err != nil {
		writeError(w, err)
		return
	}
	p, err := a.store.AddPeminjaman(Peminjaman{
		SumberDayaID: in.SumberDayaID,
		Peminjam:     in.Peminjam,
		Tanggal:      in.Tanggal,
		SlotMulai:    in.SlotMulai,
		SlotSelesai:  in.SlotSelesai,
		Jumlah:       in.Jumlah,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

// GET /peminjaman?sumber_daya_id=&tanggal=&status=  (semua filter opsional)
func (a *API) listPeminjaman(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var f filterPeminjaman
	if s := q.Get("sumber_daya_id"); s != "" {
		id, err := strconv.Atoi(s)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorBody{Error: "sumber_daya_id harus berupa angka"})
			return
		}
		f.SumberDayaID = &id
	}
	f.Tanggal = strings.TrimSpace(q.Get("tanggal"))
	f.Status = strings.TrimSpace(q.Get("status"))
	writeJSON(w, http.StatusOK, a.store.ListPeminjaman(f))
}

// GET /peminjaman/{id}
func (a *API) getPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "id harus berupa angka"})
		return
	}
	p, err := a.store.GetPeminjaman(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func validasiSumberDaya(in sumberDayaInput) error {
	if in.Nama == "" {
		return errInput("nama wajib diisi")
	}
	if utf8.RuneCountInString(in.Nama) > 50 { // per karakter, bukan per byte
		return errInput("nama maksimal 50 karakter")
	}
	if in.Jenis != "ruang" && in.Jenis != "alat" {
		return errInput(`jenis harus "ruang" atau "alat"`)
	}
	if in.Stok < 1 {
		return errInput("stok minimal 1")
	}
	return nil
}

func validasiPeminjaman(in peminjamanInput) error {
	if in.Peminjam == "" {
		return errInput("peminjam wajib diisi")
	}
	if utf8.RuneCountInString(in.Peminjam) > 50 {
		return errInput("peminjam maksimal 50 karakter")
	}
	// time.Parse menolak tanggal yang tidak ada (30 Februari), dan menerima 29 Feb 2028.
	if _, err := time.Parse("2006-01-02", in.Tanggal); err != nil {
		return errInput("tanggal harus berformat YYYY-MM-DD dan tanggal yang valid")
	}
	if in.SlotMulai < 1 || in.SlotMulai > 8 || in.SlotSelesai < 1 || in.SlotSelesai > 8 {
		return errInput("slot harus antara 1 sampai 8")
	}
	if in.SlotMulai > in.SlotSelesai {
		return errInput("slot_mulai tidak boleh lebih besar dari slot_selesai")
	}
	if in.Jumlah < 1 {
		return errInput("jumlah minimal 1")
	}
	return nil
}
