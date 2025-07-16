# 🧠 image-catalog-service

The `image-catalog-service` is a backend service for managing and retrieving medical image metadata stored in Firestore. It provides filtering, retrieval, and update endpoints for images processed by the image-processing pipeline. Additionally, it includes a secure GCS proxy endpoint to serve tile-based image data and related resources (e.g., DZI, thumbnails) required by clients such as OpenSeadragon.

---

## 🧩 Features

- 🔍 Filter and retrieve image records from Firestore
- 🔄 Update or delete image metadata
- 🧵 Serve GCS-based resources (e.g., Deep Zoom tiles) via a secure proxy
- 🛡️ Designed to sit behind an authentication gateway

---

## ⚙️ Dependencies

Make sure the following tools and services are configured before running this project:

- [Go 1.21+](https://golang.org/dl/)
- [Google Cloud Project](https://console.cloud.google.com/)
- Firestore (in Native mode)
- Google Cloud Storage (GCS) bucket
- Valid Service Account JSON key with access to GCS and Firestore



---

## 📦 Installation

```bash
git clone https://github.com/your-org/image-catalog-service.git
cd image-catalog-service

# Set up .env file based on .env.example
cp .env.example .env
```

Install dependencies and run:

```bash
go mod tidy
go run cmd/main.go
```

---

## 🔧 Environment Variables

```env
ENV=LOCAL
GIN_MODE=debug                   # Use "release" in production
PORT=3232

# Google Cloud Configuration
GCP_PROJECT_ID=your-gcp-project-id
GCP_REGION=us-central1
GCS_BUCKET_NAME=your-image-catalog-bucket
GCS_BUCKET_LOCATION=US-CENTRAL1
GCS_STORAGE_CLASS=STANDARD
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json

# Timeouts
READ_TIMEOUT=15m
WRITE_TIMEOUT=60s
IDLE_TIMEOUT=5m
```

---

## 📡 Sample API Requests

### 🔎 Get Image by ID

```bash
curl -X GET http://localhost:3232/api/v1/images/{image_id}
```

---

### 📋 List Images with Filters

```bash
curl -X GET "http://localhost:3232/api/v1/images?dataset_name=CMB-BRCA&organ_type=breast"
```

---

### ✏️ Update Image Metadata

```bash
curl -X PUT http://localhost:3232/api/v1/images/{image_id} \
  -H "Content-Type: application/json" \
  -d '{
    "disease_type": "carcinoma",
    "classification": "carcinoma",
    "sub_type": "ductal",
    "grade": ""
  }'
```

---

### 🗑️ Delete an Image

```bash
curl -X DELETE http://localhost:3232/api/v1/images/{image_id}
```

---

### 🌐 Proxy a GCS Object (e.g., tiles, thumbnails)

```bash
curl -X GET http://localhost:3232/api/v1/proxy/1752612491902535632/image_files/10/0_0.jpeg
```

This proxy handles all GCS paths and streams the file directly, avoiding public signed URLs.

---

## 🧑‍💻 Developer Notes

- Tile, thumbnail, and DZI resources are private and **proxied** through this service.
- You can later enhance the system by:
  - Adding job tracking (`job_id`) support
  - Enabling pagination or sorting for image lists
  - Integrating full text search via Firestore indexing

---

## 📜 License

MIT — feel free to use, adapt, and extend.
