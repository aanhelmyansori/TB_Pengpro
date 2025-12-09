package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Barang merepresentasikan satu item inventaris
type Barang struct {
	Kode     string
	Nama     string
	Kategori string
	Stok     int
}

// konstanta — opsional, tapi sesuai standar akademik (Tel-U)
const MAX_BARANG = 100

func main() {
	var inventaris []Barang // slice dinamis — lebih idiomatik Go
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("============================================")
	fmt.Println("   SISTEM MANAJEMEN INVENTARIS GUDANG")
	fmt.Println("============================================")

	for {
		fmt.Println("[1] Tambah Barang")
		fmt.Println("[2] Lihat Daftar Barang")
		fmt.Println("[3] Update Stok (Masuk/Keluar)")
		fmt.Println("[4] Hapus Barang")
		fmt.Println("[5] Laporan Ringkasan")
		fmt.Println("[6] Keluar")
		fmt.Print("Pilihan Anda: ")

		if !scanner.Scan() {
			break
		}
		pilih := strings.TrimSpace(scanner.Text())

		switch pilih {
		case "1":
			tambahBarang(&inventaris, scanner)
		case "2":
			lihatDaftar(inventaris)
		case "3":
			updateStok(&inventaris, scanner)
		case "4":
			hapusBarang(&inventaris, scanner)
		case "5":
			laporanRingkasan(inventaris)
		case "6":
			if konfirmasiKeluar(scanner) {
				fmt.Println("Terima kasih! Sampai jumpa di gudang berikutnya.")
				return
			}
		default:
			fmt.Println("Pilihan tidak valid. Silakan coba lagi.")
		}

		fmt.Println("--------------------------------------------")
	}
}

// === CASE 1: Tambah Barang ===
func tambahBarang(inventaris *[]Barang, scanner *bufio.Scanner) {
	fmt.Println("Tambah Barang Baru")

	if len(*inventaris) >= MAX_BARANG {
		fmt.Printf("Gudang penuh (maks %d barang). Tidak dapat menambah barang.\n", MAX_BARANG)
		return
	}

	fmt.Print("Kode barang (unik, contoh: BRG-001): ")
	scanner.Scan()
	kode := strings.TrimSpace(scanner.Text())
	if kode == "" {
		fmt.Println("Kode tidak boleh kosong.")
		return
	}

	// Cek keunikan
	for _, b := range *inventaris {
		if b.Kode == kode {
			fmt.Println("Kode barang sudah digunakan. Gagal menambahkan.")
			return
		}
	}

	fmt.Print("Nama barang: ")
	scanner.Scan()
	nama := strings.TrimSpace(scanner.Text())
	if nama == "" {
		fmt.Println("Nama tidak boleh kosong.")
		return
	}

	fmt.Print("Kategori: ")
	scanner.Scan()
	kategori := strings.TrimSpace(scanner.Text())
	if kategori == "" {
		fmt.Println("Kategori tidak boleh kosong.")
		return
	}

	fmt.Print("Stok awal: ")
	scanner.Scan()
	stokStr := strings.TrimSpace(scanner.Text())
	stok, err := strconv.Atoi(stokStr)
	if err != nil {
		fmt.Println("Stok harus berupa bilangan bulat.")
		return
	}
	if stok < 0 {
		fmt.Println("Stok tidak boleh negatif.")
		return
	}

	*inventaris = append(*inventaris, Barang{
		Kode:     kode,
		Nama:     nama,
		Kategori: kategori,
		Stok:     stok,
	})
	fmt.Println("Barang berhasil ditambahkan!")
}

// === CASE 2: Lihat Daftar ===
func lihatDaftar(inventaris []Barang) {
	if len(inventaris) == 0 {
		fmt.Println("Gudang masih kosong.")
		return
	}

	fmt.Printf("Daftar Barang (%d item):\n", len(inventaris))
	fmt.Println("KODE      | NAMA               | KATEGORI     | STOK")
	fmt.Println("----------|--------------------|--------------|-----")

	for _, b := range inventaris {
		status := ""
		if b.Stok <= 5 {
			status = " <-- Stok kritis!"
		}
		fmt.Printf("%-10s| %-18s| %-12s| %-4d%s\n",
			b.Kode, b.Nama, b.Kategori, b.Stok, status)
	}
}

// === CASE 3: Update Stok ===
func updateStok(inventaris *[]Barang, scanner *bufio.Scanner) {
	fmt.Print("Masukkan kode barang: ")
	scanner.Scan()
	kode := strings.TrimSpace(scanner.Text())

	idx := -1
	for i, b := range *inventaris {
		if b.Kode == kode {
			idx = i
			break
		}
	}

	if idx == -1 {
		fmt.Println("Barang tidak ditemukan.")
		return
	}

	barang := (*inventaris)[idx]
	fmt.Printf("Barang ditemukan: %s (Stok: %d)\n", barang.Nama, barang.Stok)
	fmt.Println("[1] Stok Masuk (tambah)")
	fmt.Println("[2] Stok Keluar (kurangi)")
	fmt.Print("Pilihan: ")
	scanner.Scan()
	pilih := strings.TrimSpace(scanner.Text())

	fmt.Print("Jumlah: ")
	scanner.Scan()
	jumlahStr := strings.TrimSpace(scanner.Text())
	jumlah, err := strconv.Atoi(jumlahStr)
	if err != nil || jumlah <= 0 {
		fmt.Println("Jumlah harus bilangan bulat > 0.")
		return
	}

	switch pilih {
	case "1":
		(*inventaris)[idx].Stok += jumlah
		fmt.Printf("Stok bertambah --> %d unit.\n", (*inventaris)[idx].Stok)
	case "2":
		if (*inventaris)[idx].Stok-jumlah < 0 {
			fmt.Printf("Stok tidak cukup! Sisa: %d\n", (*inventaris)[idx].Stok)
		} else {
			(*inventaris)[idx].Stok -= jumlah
			fmt.Printf("Stok berkurang --> %d unit.\n", (*inventaris)[idx].Stok)
		}
	default:
		fmt.Println("Pilihan tidak valid.")
	}
}

// === CASE 4: Hapus Barang ===
func hapusBarang(inventaris *[]Barang, scanner *bufio.Scanner) {
	fmt.Print("Kode barang yang akan dihapus: ")
	scanner.Scan()
	kode := strings.TrimSpace(scanner.Text())

	idx := -1
	for i, b := range *inventaris {
		if b.Kode == kode {
			idx = i
			break
		}
	}

	if idx == -1 {
		fmt.Println("Barang tidak ditemukan.")
		return
	}

	nama := (*inventaris)[idx].Nama
	fmt.Printf("Hapus barang '%s'? (y/n): ", nama)
	scanner.Scan()
	konfirmasi := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if konfirmasi != "y" {
		return
	}

	// Hapus dengan slice copy (lebih aman daripada manual loop)
	*inventaris = append((*inventaris)[:idx], (*inventaris)[idx+1:]...)
	fmt.Println("Barang berhasil dihapus.")
}

// === CASE 5: Laporan Ringkasan ===
func laporanRingkasan(inventaris []Barang) {
	totalStok := 0
	stokKritis := 0
	for _, b := range inventaris {
		totalStok += b.Stok
		if b.Stok <= 5 {
			stokKritis++
		}
	}

	fmt.Println("LAPORAN RINGKASAN:")
	fmt.Printf("- Total jenis barang : %d\n", len(inventaris))
	fmt.Printf("- Total stok         : %d\n", totalStok)
	fmt.Printf("- Barang stok kritis : %d\n", stokKritis)
}

// === CASE 6: Konfirmasi Keluar ===
func konfirmasiKeluar(scanner *bufio.Scanner) bool {
	fmt.Print("Keluar dari program? (y/n): ")
	scanner.Scan()
	konfirmasi := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return konfirmasi == "y"
}