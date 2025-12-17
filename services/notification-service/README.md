# Notification Service

Service ini bertanggung jawab untuk mengirimkan notifikasi kepada pengguna melalui berbagai channel (Email, WhatsApp).

## Fitur

- **Multi-channel Support**: Mengirim notifikasi via Email (SMTP) dan WhatsApp (HTTP API).
- **Template Engine**: Mendukung penggunaan template dinamis untuk notifikasi.
- **Async Processing**: Menggunakan RabbitMQ untuk pemrosesan notifikasi secara asynchronous.
- **Retry Mechanism**: Mekanisme retry otomatis untuk notifikasi yang gagal.
- **Webhook Support**: Menerima status update dari provider (misal: WhatsApp delivery status).

## Struktur Project

```
services/notification-service/
├── cmd/
│   └── server/          # Entry point aplikasi
├── internal/
│   ├── domain/          # Entities & interfaces
│   ├── handler/         # HTTP handlers
│   ├── infrastructure/  # External services (Email, WA)
│   ├── repository/      # Database access
│   └── usecase/         # Business logic
└── migrations/          # Database schemas
```

## Konfigurasi

Environment variables yang dibutuhkan:

```bash
# App
APP_ENV=development
APP_HTTP_PORT=9097

# Database
POSTGRES_URL=postgres://user:password@localhost:5432/sisfo_notification?sslmode=disable

# Message Broker
RABBIT_URL=amqp://dev:dev@localhost:5672/

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
SMTP_FROM=noreply@school.id

# WhatsApp
WA_API_URL=https://api.whatsapp.provider.com
WA_API_KEY=your-api-key
```

## API Endpoints

### Templates

- `POST /api/v1/notifications/templates` - Membuat template baru
- `GET /api/v1/notifications/templates` - List template
- `GET /api/v1/notifications/templates/:id` - Detail template
- `PUT /api/v1/notifications/templates/:id` - Update template
- `DELETE /api/v1/notifications/templates/:id` - Hapus template

### Notifications

- `POST /api/v1/notifications/send` - Mengirim notifikasi
- `GET /api/v1/notifications/:id` - Detail notifikasi
- `GET /api/v1/notifications/recipient` - List notifikasi berdasarkan penerima

### Webhooks

- `POST /webhooks/:provider` - Handle webhook dari provider

## Cara Menjalankan

### Lokal

```bash
# Jalankan service
make run-notification

# Jalankan test
make test-notification
```

### Docker

```bash
docker-compose up -d notification-service
```
