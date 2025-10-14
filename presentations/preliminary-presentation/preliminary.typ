#import "@preview/typslides:1.2.6": *

// Konfigurasi Proyek
#show: typslides.with(
  ratio: "16-9",
  theme: "bluey",
  font: "Fira Sans",
  link-style: "color",
)

// Slide Judul
#front-slide(
  title: "Proyek: URL Shortener",
  subtitle: [Pendekatan Microservices Cloud-Native],
  authors: ("Anggota Tim 1", "Anggota Tim 2", "Anggota Tim 3"),
  info: [Teknologi Cloud - Presentasi Awal],
)

// Daftar Isi
#table-of-contents()

// --- Bagian 1: Gambaran Umum Proyek ---
#title-slide[
  Gambaran Umum Proyek
]

#slide(title: "Tujuan: Lebih dari Sekadar URL Shortener")[
  - *Konsep:* Membangun backend lengkap untuk sebuah URL shortener berperforma tinggi seperti bit.ly.

  - *Tujuan Utama:* Proyek ini adalah sarana untuk mengimplementasikan dan menguasai pola-pola cloud-native. Fokus kami adalah pada #stress("arsitektur"), bukan hanya aplikasi akhirnya.

  - *Teknologi Inti yang Akan Digunakan:*
    - Arsitektur Microservices
    - Komunikasi Asinkron berbasis Event
    - Compute Serverless & Container
    - Infrastructure as Code & CI/CD
]

// --- Bagian 2: Arsitektur ---
#title-slide[
  Arsitektur yang Diusulkan
]

#slide(title: "Rancangan Sistem")[
  Kami merancang sistem yang *decoupled* di mana setiap layanan memiliki satu tanggung jawab. Ini mendukung pengembangan, deployment, dan scaling yang independen.

  #cols(columns: (1fr, 1.2fr))[
    [
    *Komponen Kunci:*
    - *API Gateway:* Single entry point.
    - *Link Service:* Menangani pembuatan link.
    - *Redirect Service:* High-performance read path.
    - *Event Hub:* Melakukan ingest data klik mentah.
    - *Analytics Service:* Memproses data secara asinkron.
    - *Multiple Databases:* SQL, NoSQL, dan Cache.
    ]
  ][
    #framed[
      *Diagram Sederhana:*
      ```mermaid
      graph TD
          User -- Klik Link --> APIGateway
          APIGateway --> RedirectSvc["Redirect Service"]
          RedirectSvc -- Kirim Event --> EventHub
          RedirectSvc -- Redirect User --> ExternalSite

          EventHub -- Memicu --> AnalyticsSvc["Analytics Function"]
          AnalyticsSvc -- Menulis ke --> AnalyticsDB
      ```
    ]
  ]
]

// --- Bagian 3: Peran & Tanggung Jawab ---
#title-slide[
  Peran & Tanggung Jawab
]

#slide(title: "Struktur Tim")[
  Kami membagi kepemilikan berdasarkan domain arsitektur untuk memastikan fokus dan kejelasan.

  #cols(columns: (1fr, 1fr, 1fr), gutter: 1em)[
    #framed(title: "Person 1: API & Gateway")[
      - Bertanggung jawab atas `Link Management Service` dan `Analytics Query Service`.
      - Mengelola konfigurasi API Gateway (APIM) dan kontrak API publik.
    ]
  ][
    #framed(title: "Person 2: Redirect & Data")[
      - Bertanggung jawab atas `Redirect Service` yang berkinerja tinggi.
      - Mengelola database (SQL & Redis).
      - Mengimplementasikan data ingestion ke Event Hubs.
    ]
  ][
    #framed(title: "Person 3: Backend & DevOps")[
      - Bertanggung jawab atas `Analytics Processing Service` (Azure Function) yang asinkron.
      - Mengelola pipeline CI/CD (GitHub Actions).
      - Menyiapkan shared infrastructure (ACR, Key Vault).
    ]
  ]
]

// --- Bagian 4: Progres Awal ---
#title-slide[
  Progres Awal (Minggu ke-1)
]

#slide(title: "Fondasi Telah Dibangun")[
  Kami telah menyelesaikan penyiapan fondasi untuk mempercepat pengembangan di minggu-minggu mendatang.

  - *[x] Version Control:* Repositori GitHub telah dibuat dan dikonfigurasi.

  - *[x] Penyiapan Azure:*
    - Sebuah #stress("Resource Group") utama telah disiapkan.
    - Layanan bersama telah di-deploy:
      - #strong("Azure Container Registry (ACR)") untuk Docker image.
      - #strong("Azure Key Vault") untuk manajemen secret yang aman.
    - Aturan networking awal telah dipertimbangkan.

  - *[x] Desain API:*
    - Draf awal #stress("spesifikasi OpenAPI (Swagger)") telah dibuat. Ini memberikan kontrak yang jelas untuk diikuti semua layanan.
]

// --- Bagian 5: Rencana Kerja ---
#title-slide[
  Rencana Kerja 6 Minggu
]

#slide(title: "Lini Masa Proyek")[
  - *Minggu 1:* Fondasi & Penyiapan. #greeny("[Selesai]")

  - *Minggu 2:* #yelly("Implementasi 'Write Path'.") Tujuannya adalah membuat dan menyimpan link melalui API yang sudah di-deploy.

  - *Minggu 3:* #yelly("Implementasi 'Read Path'.") Tujuannya adalah sebuah `Redirect Service` berkecepatan tinggi yang berfungsi.

  - *Minggu 4:* #yelly("Integrasi Gateway & Analytics Ingestion.") Semua layanan di belakang APIM; event klik mulai ditangkap.

  - *Minggu 5:* #yelly("Menyelesaikan Alur Analitik.") Memproses event dan membuat data analitik bisa di-query.

  - *Minggu 6:* #yelly("Otomatisasi Penuh (CI/CD) & Finalisasi.") Semua deployment diotomatisasi melalui GitHub Actions.
]

// --- Bagian 6: Langkah Berikutnya ---
#title-slide[
  Langkah Berikutnya
]

#slide(title: "Fokus Kami Saat Ini")[
  Tujuan kami untuk pertemuan berikutnya adalah mendemonstrasikan "Write Path" yang berfungsi penuh.

  - *Apa yang akan kami tunjukkan:*
    - Sebuah `Link Management` microservice yang di-deploy di Azure Container Apps.
    - Panggilan API yang berhasil melalui Postman yang membuat record di database Azure SQL.
    - Semua secret dikelola dengan aman melalui Azure Key Vault.

  - Kami terbuka untuk pertanyaan dan masukan awal mengenai rencana kami.
]

#focus-slide[
  Terima Kasih.
]
