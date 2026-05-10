# CI/CD Pipeline: Deploy API Go-SQLite ke VPS via GitHub Actions & Podman

## 1. Persiapan Struktur Folder & Dockerfile

Langkah pertama adalah memastikan `Dockerfile` berada di folder yang benar agar aplikasi bisa dibuild dengan dukungan SQLite (CGO).

**Struktur Folder:**

```text
notes/
├── server/
│   └── sqllite/
│       ├── Dockerfile
│       ├── main.go
│       ├── go.mod
│       └── go.sum
└── .github/
    └── workflows/
        └── deploy.yml

```

**Dockerfile (`server/sqllite/Dockerfile`):**

```dockerfile
FROM docker.io/library/golang:alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o api-binary .

FROM docker.io/library/alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
RUN mkdir -p /root/data
COPY --from=builder /app/api-binary .
EXPOSE 8080
CMD ["./api-binary"]

```

---

## 2. Uji Coba Build Lokal (Podman)

Sebelum melakukan push ke GitHub, wajib melakukan build di laptop lokal untuk memastikan tidak ada library yang kurang atau error pada kode Go.

**Perintah Build:**

```bash
cd server/sqllite
podman build -t test-api .

```

**Perintah Menjalankan Container:**

```bash
podman run -d --name api-lokal -p 8080:8080 -v ~/data-lokal:/root/data:Z test-api

```

> **Catatan:** Gunakan flag `:Z` pada volume agar Podman memiliki izin akses (SELinux) ke folder database di host.

---

## 3. Konfigurasi GitHub Secrets

Agar GitHub Actions bisa masuk ke VPS dan melakukan build, anda perlu mendaftarkan "kunci-kunci" di menu **Settings > Secrets and variables > Actions > New repository secret**.

| Nama Secret | Nilai / Isi |
| --- | --- |
| **HOST** | IP Public VPS anda (contoh: `34.133.107.178`) |
| **USERNAME** | Username SSH VPS (contoh: `pronus_programmernusantara`) |
| **SSH_KEY** | Isi file **Private Key** (`id_ed25519`) yang dibuat di server |

**PENTING:** Berikan izin tulis ke GitHub Actions di **Settings > Actions > General > Workflow permissions**, pilih **"Read and write permissions"**.

---

## 4. Automasi CI/CD (GitHub Actions)

File ini berfungsi sebagai otak yang melakukan build image secara otomatis saat anda melakukan `git push`.

**File Workflow (`.github/workflows/deploy.yml`):**

```yaml
name: Deploy API SQLite
on:
  push:
    branches: [main]
env:
  FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4
      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and Push Image
        uses: docker/build-push-action@v6
        with:
          context: ./server/sqllite
          push: true
          tags: ghcr.io/${{ github.repository }}:latest

      - name: Deploy to VPS via SSH
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.HOST }}
          username: ${{ secrets.USERNAME }}
          key: ${{ secrets.SSH_KEY }}
          script: |
            podman login ghcr.io -u ${{ github.actor }} -p ${{ secrets.GITHUB_TOKEN }}
            podman pull ghcr.io/${{ github.repository }}:latest
            podman stop go-api || true
            podman rm go-api || true
            podman run -d \
              --name go-api \
              -p 8080:8080 \
              -v /home/pronus_programmernusantara/sqlite-data:/root/data:Z \
              ghcr.io/${{ github.repository }}:latest

```

---

## 5. Setup Keamanan & SSH di VPS

Agar GitHub Actions bisa masuk tanpa password, anda harus membuat kunci SSH di dalam VPS itu sendiri.

**Langkah-langkah di Terminal VPS:**

1. **Generate Key:** `ssh-keygen -t ed25519 -C "github-actions"`
2. **Daftarkan Key:** `cat ~/.ssh/id_ed25519.pub >> ~/.ssh/authorized_keys`
3. **Ambil Private Key:** `cat ~/.ssh/id_ed25519`
* Copy teks yang muncul (mulai dari `-----BEGIN...` sampai `...END-----`) lalu masukkan ke GitHub Secret **SSH_KEY**.



**Persiapan Folder Database:**

```bash
mkdir -p ~/sqlite-data

```

---

## 6. Cara Penggunaan

Setiap kali anda ingin memperbarui aplikasi di server, anda cukup menjalankan perintah standar Git di laptop lokal:

```bash
git add .
git commit -m "Update fitur baru"
git push origin main

```

**Hasil Akhir:**

* GitHub akan build image aplikasi anda.
* Image disimpan di **GHCR** (GitHub Container Registry).
* GitHub masuk ke VPS via SSH.
* Podman menarik image terbaru dan merestart container secara otomatis.
* Aplikasi dapat diakses di `http://IP-VPS:8080`.
