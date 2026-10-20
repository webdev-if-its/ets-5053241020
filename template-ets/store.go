package main

import (
	"errors"
	"fmt"
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
	mu           sync.Mutex
	sumberDaya   []SumberDaya
	peminjaman   []Peminjaman
	nextSumberID int
	nextPinjamID int
}

func NewStore() *Store {
	return &Store{
		sumberDaya:   []SumberDaya{},
		peminjaman:   []Peminjaman{},
		nextSumberID: 1,
		nextPinjamID: 1,
	}
}

type filterPeminjaman struct {
	SumberDayaID *int
	Tanggal      string
	Status       string
}

func (s *Store) cariSumberDaya(id int) (SumberDaya, bool) {
	for _, sd := range s.sumberDaya {
		if sd.ID == id {
			return sd, true
		}
	}
	return SumberDaya{}, false
}

func (s *Store) AddSumberDaya(sd SumberDaya) SumberDaya {
	s.mu.Lock()
	defer s.mu.Unlock()
	sd.ID = s.nextSumberID
	s.nextSumberID++
	s.sumberDaya = append(s.sumberDaya, sd)
	return sd
}

func (s *Store) ListSumberDaya() []SumberDaya {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]SumberDaya, len(s.sumberDaya))
	copy(out, s.sumberDaya)
	return out
}

func (s *Store) GetSumberDaya(id int) (SumberDaya, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sd, ok := s.cariSumberDaya(id)
	if !ok {
		return SumberDaya{}, ErrSumberDayaTidakDitemukan
	}
	return sd, nil
}

func (s *Store) AddPeminjaman(p Peminjaman) (Peminjaman, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sd, ok := s.cariSumberDaya(p.SumberDayaID)
	if !ok {
		return Peminjaman{}, errInput("sumber_daya_id tidak menunjuk sumber daya yang ada")
	}
	if p.Jumlah > sd.Stok {
		return Peminjaman{}, errInput("jumlah melebihi stok sumber daya")
	}

	p.ID = s.nextPinjamID
	s.nextPinjamID++
	p.Status = "aktif"
	s.peminjaman = append(s.peminjaman, p)
	return p, nil
}

func (s *Store) ListPeminjaman(f filterPeminjaman) []Peminjaman {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []Peminjaman{} // bukan nil -> JSON [] bukan null
	for _, p := range s.peminjaman {
		if f.SumberDayaID != nil && p.SumberDayaID != *f.SumberDayaID {
			continue
		}
		if f.Tanggal != "" && p.Tanggal != f.Tanggal {
			continue
		}
		if f.Status != "" && p.Status != f.Status {
			continue
		}
		out = append(out, p)
	}
	return out
}

func (s *Store) GetPeminjaman(id int) (Peminjaman, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.peminjaman {
		if p.ID == id {
			return p, nil
		}
	}
	return Peminjaman{}, ErrPeminjamanTidakDitemukan
}

var (
	ErrSumberDayaTidakDitemukan = errors.New("sumber daya tidak ditemukan")
	ErrPeminjamanTidakDitemukan = errors.New("peminjaman tidak ditemukan")
	ErrInputTidakValid          = errors.New("input tidak valid")
)

func errInput(pesan string) error {
	return fmt.Errorf("%w: %s", ErrInputTidakValid, pesan)
}
