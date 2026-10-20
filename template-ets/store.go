package main

import (
	"errors"
	"sync"
)

const (
	StatusAktif = "aktif"
	JenisRuang  = "ruang"
	JenisAlat   = "alat"
	MaxSlot     = 8
)

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"`
}

type Store struct {
	mu               sync.Mutex
	sumberDaya       []SumberDaya
	peminjaman       []Peminjaman
	nextSumberDayaID int
	nextPeminjamanID int
}

func NewStore() *Store {
	return &Store{
		sumberDaya:       []SumberDaya{},
		peminjaman:       []Peminjaman{},
		nextSumberDayaID: 1,
		nextPeminjamanID: 1,
	}
}

type filterPeminjaman struct {
	SumberDayaID *int
	Tanggal      string
	Status       string
}

var (
	ErrSumberDayaTidakDitemukan = errors.New("sumber daya tidak ditemukan")
	ErrPeminjamanTidakDitemukan = errors.New("peminjaman tidak ditemukan")
	ErrInputTidakValid          = errors.New("input tidak valid")
	ErrBentrok                  = errors.New("bentrok")
)
