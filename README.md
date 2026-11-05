# Tarna Minecraft Bedrock Server

Server Minecraft Bedrock Edition berkinerja tinggi yang dioptimalkan khusus untuk **Armbian Server H680P** dan dikembangkan menggunakan bahasa pemrograman Go.

---

## Cara Pasang & Menjalankan (Sangat Mudah)

### Langkah 1: Persiapan (Memasang Go)
Pastikan server Armbian Anda sudah terpasang compiler Go (versi 1.26 atau lebih tinggi).
```shell
# Untuk pengguna Debian/Ubuntu (Armbian)
sudo apt update && sudo apt install golang -y
```

### Langkah 2: Jalankan Server
Masuk ke direktori server dan jalankan perintah berikut untuk memulai server secara langsung:
```shell
go run main.go
```

Atau jika ingin melakukan kompilasi terlebih dahulu agar server berjalan lebih cepat:
```shell
# Kompilasi menjadi file binary
go build -o mc_server main.go

# Menjalankan server
./mc_server
```
Tekan tombol `Ctrl+C` di terminal untuk mematikan server dengan aman.

---

## Fitur Khusus Server Ini:
1. **Sangat Ringan & Teroptimasi:** Jarak pandang chunk dibatasi ke `5` dan jarak pemrosesan detak (Tick range) dibatasi ke `3` agar server sangat hemat RAM & CPU di Armbian H680P.
2. **Sistem Ban Multi-faktor:** Melalui perintah `/ban <pemain> [durasi] [alasan]`, server memblokir akses berdasarkan **Username, XUID, UUID, dan IP** agar pemain yang diban tidak bisa masuk kembali meskipun mengganti namanya.
3. **Tombol Unban Instan:** Admin bisa menghapus ban secara instan menggunakan `/uban <pemain>` atau `/unban <pemain>`.
4. **Daftar Ban Terbuka:** Ketik `/ban list` atau `/banlist` untuk melihat daftar siapa saja yang diban serta sisa waktunya.
5. **Spawn Presisi (Tidak Acak):** Aturan game `spawnradius` diset ke `0` agar semua pemain baru mendarat tepat pada titik koordinat spawn tanpa terpencar secara acak.
6. **Sistem Cheat & Administrasi Khusus:** Hak istimewa untuk mengeksekusi cheat dan perintah admin (seperti `/tp`, `/give`, `/heal`, `/time`, `/ban`, `/uban`, `/banlist`, `/gamemode`) hanya dapat dijalankan oleh dua akun administrator:
   * **XTarnaWijaya**
   * **TarnaWijaya**
7. **Sistem Pemunculan Mob Berkala (Mob Spawner):** Memunculkan mob pasif & agresif (Babi, Sapi, Domba, Zombie, Skeleton) secara acak di sekitar pemain setiap 20 detik, dengan batasan maksimal 25 mob demi menjaga performa server.
8. **Pembersihan Item & Optimasi Memori Otomatis:** Menghapus item jatuh di tanah secara otomatis setiap 15 menit dengan notifikasi peringatan 5 menit sebelumnya. Server juga melakukan Garbage Collection (`runtime.GC()`) dan pengembalian memori OS (`debug.FreeOSMemory()`) untuk menjaga penggunaan RAM tetap minimal.
9. **Peningkatan Generator Dunia (Enhanced Vanilla Generator):**
   * **Sistem Goa (Cave Carving):** Membuat goa bawah tanah secara prosedural menggunakan perhitungan trigonometri cepat.
   * **Pembangkitan Bijih (Ore Generation):** Menghasilkan blok bijih (Batu Bara, Besi, Tembaga, Emas, Redstone, Lapis Lazuli, Berlian, Zamrud) sesuai ketinggian/kedalaman.
   * **Vegetasi & Pepohonan:** Menumbuhkan rumput pendek, pakis, dandelion, serta pohon-pohon dengan bentuk batang/daun yang lebih alami.
   * **Rumah Desa (Village Houses):** Menghasilkan struktur rumah sederhana di permukaan tanah dengan peluang 3% per chunk.
10. **Teleportasi Acak saat Masuk (Cooperative Random Spawn):** Pemain baru otomatis masuk dengan Mode Survival dan diteleportasi secara acak dalam radius 40 blok di dekat pemain aktif lain (atau di dekat spawn utama jika sendirian) untuk mendukung permainan kooperatif.
11. **Perintah Mode Game (`/gamemode` / `/gm`):** Memungkinkan administrator untuk mengubah mode permainan (Survival, Creative, Adventure, Spectator) bagi diri sendiri maupun pemain lain.

---

## Lisensi

Proyek ini dilisensikan di bawah **GNU Affero General Public License v3.0 (AGPL-3.0)** — lisensi copyleft paling ketat yang mewajibkan penyediaan kode sumber lengkap jika server dijalankan atau dimodifikasi untuk diakses melalui jaringan. Lihat berkas [LICENSE](file:///data/data/com.termux/files/home/mc/LICENSE) untuk informasi lebih lanjut.

