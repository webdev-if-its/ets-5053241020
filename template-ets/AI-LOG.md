# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Deskripsi Tugas
- Prompt: Menjelaskan step-by-step cara pengerjaan
- Dipakai: Langkah-langkah pengerjaan
- Diubah: Tidak ada

## Koreksi Rancangan Struktur Data
- Prompt: Mengirimkan struktur yang sudah saya buat, lalu minta koreksi
- Dipakai: Struct Store, karena sebelumnya saya masih salah dalam membuat struct Store
- Diubah: Sisanya menggunakan rancangan yang sudah saya buat

## Koreksi Endpoint
- Prompt: Mengirimkan endpoint yang sudah saya buat lalu minta untuk dikoreksi
- Dipakai: Perbaikan kesalahan
- Diubah: tidak ada

## Koreksi Fungsi Service
- Prompt: Mengirimkan kode yang sudah saya buat
- Dipakai: Perbaikan kesalahan
- Diubah: tidak ada

## Implementasi Error Bentrok
- Prompt: Meminta bantuan mengimplementasikan deteksi bentrok peminjaman dan penanganan response HTTP 409 Conflict
- Dipakai: Penambahan error ErrBentrok, pengecekan bentrok pada proses peminjaman, dan mapping error ke status 409
- Diubah: Penyesuaian nama fungsi dan integrasi dengan struktur kode yang sudah ada

## Implementasi Persona A
- Prompt: Meminta bantuan menyelesaikan requirement Persona A 
- Dipakai: Perhitungan ketersediaan per slot dan rancangan endpoint
- Diubah: Format response dan implementasi disesuaikan dengan struktur API yang dibuat

## Implementasi Change Request Frontend B
- Prompt: Meminta bantuan mengerjakan Change Request B
- Dipakai: Penambahan field bentrok_dengan pada response 409 dan logika pencarian ID peminjaman yang beririsan
- Diubah: Struktur response dan implementasi fungsi disesuaikan dengan kode saya

## Penyusunan PROSES.md bagian Desain, Bug, dan Change Request
- Prompt: Meminta bantuan mengisi bagian keputusan desain, bug atau kendala, dan change request pada PROSES.md.
- Dipakai: Menyusun penjelasan keputusan desain pengecekan bentrok per slot, mendokumentasikan bug yang ditemukan selama pengembangan, serta menjelaskan perubahan yang dilakukan pada change request daftar bentrok.
- Diubah: Detail isi disesuaikan dengan bahasa saya.

