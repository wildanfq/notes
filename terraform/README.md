# Struktur Project Terraform (Google Cloud)

## A. Struktur Kode Terraform

Buat sebuah folder untuk proyek Anda (misal: `proyek-vps`), lalu buat dua file di bawah ini di dalam folder tersebut.

### 1. `variables.tf`

File ini menyimpan semua konfigurasi dasar agar kode utama Anda tetap bersih dan mudah diubah di satu tempat.

```hcl
variable "project_id" {
  description = "ID Project Google Cloud Anda"
  default     = "coral-box-495211-p8"
}

variable "region" {
  description = "Region utama untuk resource"
  default     = "asia-southeast2"
}

variable "zone" {
  description = "Zona spesifik untuk penempatan server"
  default     = "asia-southeast2-a"
}

variable "machine_type" {
  description = "Spesifikasi ukuran server (Compute Engine)"
  default     = "e2-micro"
}

```

### 2. `main.tf`

Ini adalah cetak biru (*blueprint*) utama infrastruktur Anda, disusun dari fondasi (Provider & Jaringan) hingga ke atas (Server/VM).

```hcl
# ==========================================
# 0. KONFIGURASI PROVIDER
# ==========================================
provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

# ==========================================
# 1. FONDASI JARINGAN (VPC)
# ==========================================
resource "google_compute_network" "vpc_network" {
  name                    = "backend-vpc"
  auto_create_subnetworks = true
}

# ==========================================
# 2. KEAMANAN (FIREWALL)
# ==========================================
resource "google_compute_firewall" "ijinkan_akses" {
  name    = "buka-port-backend"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = ["22", "80", "443", "8080"] # SSH, HTTP, HTTPS, & Port API Backend
  }

  source_ranges = ["0.0.0.0/0"] # Mengizinkan akses dari seluruh internet
}

# ==========================================
# 3. SERVER (COMPUTE ENGINE)
# ==========================================
resource "google_compute_instance" "backend_server" {
  name         = "server-backend-utama"
  machine_type = var.machine_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2204-lts"
      size  = 20 # Kapasitas disk (GB)
    }
  }

  network_interface {
    network = google_compute_network.vpc_network.name
    access_config {
      # Blok ini kosong sudah cukup untuk meminta IP Publik Ephemeral
    }
  }

  # Script otomatis saat server pertama kali menyala
  metadata_startup_script = <<-EOT
    sudo apt-get update
    sudo apt-get install -y podman
  EOT

  labels = {
    env = "production"
  }
}

# ==========================================
# 4. OUTPUT (INFORMASI PENTING)
# ==========================================
output "ip_publik_server" {
  value       = google_compute_instance.backend_server.network_interface.0.access_config.0.nat_ip
  description = "Gunakan IP Publik ini untuk SSH atau mengakses API backend di browser"
}

```

---

## B. Alur Kerja (Langkah Eksekusi)

Berikut adalah panduan dari tahap persiapan di lokal hingga server Anda siap digunakan di Cloud.

### Tahap 1: Persiapan Lingkungan Kerja (Lokal)

Siapkan "alat tempur" di laptop Anda:

1. Unduh dan instal **Terraform**.
2. Unduh dan instal **Google Cloud CLI (gcloud)**.
3. Pastikan keduanya sudah terdaftar di *System PATH*. Verifikasi dengan mengetik `terraform --version` dan `gcloud --version` di terminal.

### Tahap 2: Autentikasi (Penghubung Akun)

Hubungkan laptop Anda dengan Google Cloud menggunakan metode *Application Default Credentials* (ADC) agar aman tanpa perlu mengunduh file JSON manual.

* Jalankan perintah ini di terminal:
```bash
gcloud auth application-default login

```


* Browser akan terbuka. Pilih akun Google Anda dan berikan izin (*Allow*).

### Tahap 3: Inisialisasi Proyek Terraform

Masuk ke folder proyek tempat Anda menyimpan file `main.tf` dan `variables.tf`, lalu jalankan:

```bash
terraform init

```

*Langkah ini akan mengunduh plugin/driver Google Cloud yang dibutuhkan oleh Terraform ke dalam folder Anda.*

### Tahap 4: Simulasi (Plan)

Sebelum membangun server yang sesungguhnya, lihat pratinjau apa saja yang akan dibuat oleh kode Anda:

```bash
terraform plan

```

*Pastikan tidak ada error dan rencana pembuatan infrastruktur (Jaringan, Firewall, VM) sudah sesuai.*

### Tahap 5: Eksekusi (Deploy)

Jika rencana sudah benar, perintahkan Terraform untuk mulai membangun semuanya di Google Cloud:

```bash
terraform apply

```

*Ketik `yes` saat diminta konfirmasi. Tunggu beberapa menit hingga selesai. Di akhir proses, terminal akan menampilkan **`ip_publik_server`**.*

### Tahap 6: Pengujian & Akses Masuk

Setelah server menyala, mari kita uji apakah konfigurasi berjalan sempurna:

1. **Masuk ke Server:**
Gunakan perintah gcloud untuk *remote* ke dalam server:

```bash
   gcloud compute ssh server-backend-utama --zone=asia-southeast2-a

```

2. **Cek Podman:**
Karena kita menggunakan *startup script*, Podman seharusnya sudah terinstal otomatis. Cek dengan:
```bash
podman --version

```


3. **Tes Jaringan Publik:**
Jalankan web server ringan untuk memastikan IP Publik dan Firewall berfungsi:
```bash
sudo podman run -d -p 80:80 nginx

```


Buka browser di laptop Anda, lalu ketikkan **IP Publik** yang didapat dari langkah 5. Jika halaman *"Welcome to nginx!"* muncul, selamat! Infrastruktur Anda sudah solid dan siap menampung kode *backend* Anda.
