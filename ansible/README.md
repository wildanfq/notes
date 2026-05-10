 # Ansibleuntuk mengelola dan mengotomatisasi konfigurasi server.

---

## 1. Konsep Dasar Ansible (Konfigurasi sebagai Kode)

Dalam ekosistem modern, jika Terraform bertugas membangun "bangunan" (VM, Network), maka Ansible bertugas mengisi "interior" di dalamnya. Ansible adalah alat otomatisasi yang bersifat *agentless*, artinya Anda tidak perlu menginstal software apa pun di server target. Ansible bekerja melalui jalur SSH yang sudah ada untuk mengirimkan perintah secara otomatis. Keunggulan utamanya adalah **Idempotensi**, yaitu kemampuan untuk memastikan server berada dalam kondisi yang diinginkan tanpa melakukan tindakan berulang jika kondisi tersebut sudah tercapai.

---

## 2. Instalasi dan Persiapan Kendali di Laptop

Langkah pertama dimulai dengan menginstal mesin Ansible di laptop lokal Anda (sebagai *Control Node*). Di sistem Linux seperti Nusantara/Ubuntu, Anda cukup menjalankan perintah `sudo apt update && sudo apt install ansible -y`. Setelah terinstal, Anda bisa memverifikasinya dengan perintah `ansible --version`. Di tahap ini, laptop Anda telah resmi menjadi pusat kendali yang mampu mengirimkan instruksi ke banyak server sekaligus hanya dengan satu perintah.

---

## 3. Penyusunan Inventory (Daftar Alamat Server)

Ansible membutuhkan peta untuk mengetahui server mana saja yang harus dikelola. Peta ini disebut sebagai file **Inventory**, biasanya bernama `hosts.ini`. Di dalam file ini, Anda mengelompokkan IP server dan menentukan parameter koneksi seperti user SSH dan lokasi kunci privat. Contoh penulisan yang solid adalah:
`[my_servers] 34.50.71.123 ansible_ssh_user=wildan-fq ansible_ssh_private_key_file=~/.ssh/google_compute_engine`.
Parameter tambahan seperti `ansible_ssh_common_args='-o StrictHostKeyChecking=no'` sangat disarankan agar Ansible tidak berhenti saat meminta konfirmasi keamanan sidik jari SSH (fingerprint).

---

## 4. Penulisan Playbook (Resep Otomatisasi)

Inti dari Ansible terletak pada **Playbook**, sebuah file berformat `.yml` yang berisi daftar tugas (*tasks*) yang harus dilakukan di server. Playbook ditulis dalam bahasa YAML yang sangat mudah dibaca. Setiap tugas di dalamnya mendefinisikan "kondisi akhir" yang diinginkan, seperti "Memastikan Git terinstal" atau "Memastikan folder X tersedia". Dengan Playbook, Anda tidak lagi mengetik perintah satu per satu secara manual, melainkan mendokumentasikan seluruh konfigurasi server dalam satu file yang bisa dibagikan dan dijalankan ulang kapan saja.

---

## 5. Contoh Kode Playbook (setup-backend.yml)

Berikut adalah contoh struktur Playbook yang bersih untuk menyiapkan lingkungan backend Anda. Kode ini akan melakukan update sistem, menginstal aplikasi pendukung, dan mengatur folder proyek dengan izin akses yang tepat.

```yaml
---
- name: Setup Server Backend Wildan
  hosts: my_servers
  become: yes  # Menggunakan hak akses root (sudo)

  tasks:
    - name: Memastikan sistem sudah update
      apt:
        update_cache: yes

    - name: Menginstal git, curl, dan vim
      apt:
        name: [git, curl, vim]
        state: present

    - name: Membuat folder project dengan izin user
      file:
        path: /home/wildan-fq/app-backend
        state: directory
        owner: wildan-fq
        group: wildan-fq
        mode: '0755'

```

---

## 6. Eksekusi dan Pemantauan Hasil (Deployment)

Untuk menjalankan instruksi tersebut, Anda menggunakan perintah `ansible-playbook -i hosts.ini setup-backend.yml`. Saat perintah ini berjalan, Ansible akan melakukan koneksi SSH, mengumpulkan data fakta tentang server (*gathering facts*), dan mulai menjalankan tugas satu per satu. Anda akan melihat laporan berwarna di terminal: **Kuning** berarti ada perubahan nyata, **Hijau** berarti kondisi sudah sesuai (tidak ada perubahan), dan **Merah** jika terjadi kegagalan.

---

## 7. Verifikasi dan Pengelolaan Lanjutan

Setelah proses selesai, Anda akan mendapatkan laporan akhir berupa *Play Recap*. Indikator keberhasilan yang paling utama adalah `failed=0` dan `unreachable=0`. Anda dapat memverifikasi hasil kerja Ansible dengan masuk ke server via SSH dan melihat perubahan yang terjadi, seperti mengecek versi aplikasi yang baru diinstal atau memeriksa kepemilikan folder. Dengan alur ini, Anda telah berhasil memisahkan antara pembuatan infrastruktur (Terraform) dan pengelolaan sistem (Ansible), menciptakan sistem kerja yang sangat profesional dan efisien.
