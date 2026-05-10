provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

resource "google_compute_network" "vpc_network" {
  name                    = "backend-vpc"
  auto_create_subnetworks = true
}

resource "google_compute_firewall" "ijinkan_akses" {
  name    = "buka-port-backend"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = ["22", "80", "443", "8080"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_instance" "backend_server" {
  name         = "server-backend-utama"
  machine_type = var.machine_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2204-lts"
      size  = 20
    }
  }

  network_interface {
    network = google_compute_network.vpc_network.name
    access_config {}
  }

  metadata_startup_script = <<-EOT
    sudo apt-get update
    sudo apt-get install -y podman
  EOT

  labels = {
    env = "production"
  }
}

output "ip_publik_server" {
  value = google_compute_instance.backend_server.network_interface.0.access_config.0.nat_ip
}
