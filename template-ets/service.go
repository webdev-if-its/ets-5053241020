package main

import (
	"fmt"
	"sort"
)

func errInput(pesan string) error {
	return fmt.Errorf("%w: %s", ErrInputTidakValid, pesan)
}

func menahanSlot(p Peminjaman) bool {
	return p.Status == "aktif"
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
	sd.ID = s.nextSumberDayaID
	s.nextSumberDayaID++
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

func (s *Store) adaBentrok(p Peminjaman, stok int, abaikanID int) bool {
	for slot := p.SlotMulai; slot <= p.SlotSelesai; slot++ {
		terpakai := 0
		for _, q := range s.peminjaman {
			if q.ID == abaikanID || !menahanSlot(q) {
				continue
			}
			if q.SumberDayaID != p.SumberDayaID || q.Tanggal != p.Tanggal {
				continue
			}
			if slot >= q.SlotMulai && slot <= q.SlotSelesai {
				terpakai += q.Jumlah
			}
		}
		if terpakai+p.Jumlah > stok {
			return true
		}
	}
	return false
}

func (s *Store) periksaPeminjaman(p Peminjaman, abaikanID int) error {
	sd, ok := s.cariSumberDaya(p.SumberDayaID)
	if !ok {
		return errInput("sumber_daya_id tidak menunjuk sumber daya yang ada")
	}
	if p.Jumlah > sd.Stok {
		return errInput("jumlah melebihi stok sumber daya")
	}
	if s.adaBentrok(p, sd.Stok, abaikanID) {
		return ErrBentrok
	}
	return nil
}

func (s *Store) AddPeminjaman(p Peminjaman) (Peminjaman, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.periksaPeminjaman(p, 0); err != nil {
		return Peminjaman{}, err
	}

	p.ID = s.nextPeminjamanID
	s.nextPeminjamanID++
	p.Status = "aktif"
	s.peminjaman = append(s.peminjaman, p)
	return p, nil
}

func (s *Store) ListPeminjaman(f filterPeminjaman) []Peminjaman {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []Peminjaman{}
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

func (s *Store) UpdatePeminjaman(id int, p Peminjaman) (Peminjaman, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := -1
	for i, q := range s.peminjaman {
		if q.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return Peminjaman{}, ErrPeminjamanTidakDitemukan
	}

	p.ID = id
	p.Status = s.peminjaman[idx].Status
	if err := s.periksaPeminjaman(p, id); err != nil {
		return Peminjaman{}, err
	}

	s.peminjaman[idx] = p
	return p, nil
}

func (s *Store) DeletePeminjaman(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, q := range s.peminjaman {
		if q.ID == id {
			s.peminjaman = append(s.peminjaman[:i], s.peminjaman[i+1:]...)
			return nil
		}
	}
	return ErrPeminjamanTidakDitemukan
}

func (s *Store) Ketersediaan(
	sumberDayaID int,
	tanggal string,
) ([]slotKetersediaan, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	sd, ok := s.cariSumberDaya(sumberDayaID)
	if !ok {
		return nil, ErrSumberDayaTidakDitemukan
	}

	hasil := make([]slotKetersediaan, 0, 8)

	for slot := 1; slot <= 8; slot++ {

		terpakai := 0

		for _, p := range s.peminjaman {

			if !menahanSlot(p) {
				continue
			}

			if p.SumberDayaID != sumberDayaID {
				continue
			}

			if p.Tanggal != tanggal {
				continue
			}

			if slot >= p.SlotMulai &&
				slot <= p.SlotSelesai {

				terpakai += p.Jumlah
			}
		}

		sisa := sd.Stok - terpakai

		hasil = append(hasil, slotKetersediaan{
			Slot:     slot,
			Sisa:     sisa,
			Tersedia: sisa > 0,
		})
	}

	return hasil, nil
}

func (s *Store) daftarBentrok(p Peminjaman) []int {
	ids := map[int]bool{}

	for _, q := range s.peminjaman {
		if !menahanSlot(q) {
			continue
		}

		if q.SumberDayaID != p.SumberDayaID {
			continue
		}

		if q.Tanggal != p.Tanggal {
			continue
		}

		// cek apakah rentang slot beririsan
		if p.SlotMulai <= q.SlotSelesai &&
			p.SlotSelesai >= q.SlotMulai {

			ids[q.ID] = true
		}
	}

	hasil := make([]int, 0, len(ids))

	for id := range ids {
		hasil = append(hasil, id)
	}

	sort.Ints(hasil)

	return hasil
}
