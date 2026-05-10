# Terraform untuk mengelola infrastruktur Google Cloud

## 1. Persiapan Lingkungan Kerja di Laptop

Langkah paling dasar adalah menyiapkan "alat tempur" di laptop lokal Anda. Anda perlu mengunduh **Terraform** sebagai alat orkestrasi infrastruktur dan **Google Cloud CLI (gcloud)** sebagai jembatan komunikasi antara laptop dan server Google. Pastikan setelah instalasi, folder binari keduanya sudah masuk ke dalam *System PATH* agar perintah tersebut bisa dipanggil dari folder mana pun di terminal Anda.

Untuk memulainya, instal Google Cloud CLI dengan menjalankan file `install.sh` yang telah Anda unduh sebelumnya. Setelah itu, verifikasi instalasi dengan mengetik `gcloud --version` dan `terraform --version`. Jika keduanya memunculkan angka versi, berarti fondasi di laptop Anda sudah siap untuk melangkah ke tahap autentikasi.


---

## 2. Autentikasi dan Penghubung Akun

Setelah alat terpasang, Anda harus memberikan izin akses agar Terraform bisa masuk ke proyek Google Cloud Anda. Karena adanya kebijakan organisasi yang melarang pembuatan file kunci JSON secara manual, kita menggunakan metode **Application Default Credentials (ADC)**. Metode ini jauh lebih aman karena menggunakan sesi login browser Anda untuk memberikan token akses sementara kepada Terraform.

Jalankan perintah `gcloud auth application-default login` di terminal laptop. Browser akan terbuka secara otomatis; silakan pilih akun Google Anda dan klik "Allow". Setelah berhasil, terminal akan menampilkan pesan bahwa kredensial telah disimpan dalam file lokal di folder `.config`. Dengan langkah ini, laptop Anda sekarang sudah memiliki "kunci duplikat" yang sah untuk membangun apa pun di Google Cloud tanpa perlu membuat kunci akun layanan manual di konsol web.

---

## 3. Inisialisasi Proyek Terraform

Langkah selanjutnya adalah membuat direktori kerja khusus untuk proyek Anda, misalnya `mkdir proyek-vps && cd proyek-vps`. Di dalam folder ini, Anda harus membuat file konfigurasi utama, biasanya bernama `main.tf`. File ini berfungsi sebagai cetak biru (blueprint) yang menjelaskan kepada Google Cloud jenis server apa yang ingin Anda buat, berapa kapasitas RAM-nya, di mana lokasinya (seperti Jakarta/`asia-southeast2`), hingga aturan keamanan firewall-nya.

Setelah file `main.tf` dibuat, jalankan perintah `terraform init`. Perintah ini sangat krusial karena Terraform akan mendeteksi bahwa Anda menggunakan provider Google dan secara otomatis mengunduh plugin atau "driver" yang diperlukan ke dalam folder proyek Anda. Jika muncul pesan berwarna hijau bertuliskan "Terraform has been successfully initialized!", berarti mesin Terraform Anda sudah panas dan siap untuk mengeksekusi kode.

---



## 4. Penulisan Kode Konfigurasi (Blueprint)

Berikut adalah contoh kode `main.tf` yang solid dan rapi untuk kebutuhan server backend Anda. Kode ini menggunakan provider Google, mendefinisikan jaringan VPC agar server memiliki jalur komunikasi sendiri, serta mengatur firewall agar port penting seperti SSH (22) dan HTTP (80) terbuka untuk umum.



```hcl

# Definisi Koneksi ke Google Cloud

provider "google" {

  project = "ID-PROJECT-ANDA"

  region  = "asia-southeast2"

}


# Membuat Jaringan Virtual (VPC)

resource "google_compute_network" "vpc_utama" {

  name = "jaringan-backend"

}


# Aturan Firewall agar Server Bisa Diakses

resource "google_compute_firewall" "ijinkan_akses" {

  name    = "buka-akses-publik"

  network = google_compute_network.vpc_utama.name

  allow {

    protocol = "tcp"

    ports    = ["22", "80", "8080"]

  }

  source_ranges = ["0.0.0.0/0"]

}


# Spesifikasi Server VPS (Virtual Machine)

resource "google_compute_instance" "vm_backend" {

  name         = "server-backend-utama"

  machine_type = "e2-medium"

  zone         = "asia-southeast2-a"

  boot_disk {

    initialize_params {

      image = "ubuntu-os-cloud/ubuntu-2204-lts"

      size  = 20

    }

  }



  network_interface {

    network = google_compute_network.vpc_utama.name

    access_config {} # Agar dapat IP Publik

  }



  metadata_startup_script = "sudo apt-get update && sudo apt-get install -y podman"

}



output "ip_publik" {

  value = google_compute_instance.vm_backend.network_interface.0.access_config.0.nat_ip

}

```

---

## 5. Simulasi dan Eksekusi (Deploy)

Sebelum benar-benar membangun server yang mungkin memakan biaya, Terraform menyediakan fitur simulasi melalui perintah `terraform plan`. Perintah ini akan membandingkan kondisi laptop Anda saat ini dengan kondisi di Google Cloud. Terraform akan menampilkan daftar "rencana" tentang apa saja yang akan ditambah, diubah, atau dihapus. Ini adalah tahap krusial untuk memastikan tidak ada kesalahan ketik atau pengaturan yang salah sebelum uang/kuota Anda terpakai.

Jika rencana sudah sesuai, jalankan perintah pamungkas: `terraform apply`. Ketik `yes` saat diminta konfirmasi. Terraform akan bekerja di latar belakang, menghubungi API Google Cloud, membuat jaringan, mengatur firewall, dan menyalakan mesin VPS Anda. Tunggu hingga proses selesai dan alamat IP publik server muncul di terminal Anda. IP inilah yang menjadi "pintu masuk" utama ke server baru Anda.

---

## 6. Pengujian dan Akses Masuk

Setelah server aktif, langkah terakhir adalah memverifikasi apakah server tersebut bisa digunakan. Gunakan perintah `gcloud compute ssh server-backend-utama --zone=asia-southeast2-a` untuk masuk ke dalam terminal server. Karena kita menyertakan *startup script* di dalam kode Terraform, server tersebut akan otomatis menginstal Podman sesaat setelah menyala. Anda bisa mengetesnya dengan perintah `podman --version` di dalam server.

Untuk pengujian akhir secara visual, jalankan container sederhana dengan `sudo podman run -d -p 80:80 nginx`. Ambil IP publik yang muncul di output Terraform tadi, lalu tempelkan di browser laptop Anda. Jika halaman "Welcome to nginx" muncul, berarti alur dari bawah (laptop) ke atas (server cloud) telah terjalin dengan sangat solid dan infrastruktur Anda siap digunakan untuk tahap pengembangan backend selanjutnya.
