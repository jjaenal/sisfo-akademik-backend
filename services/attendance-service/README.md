# Attendance Service

Service ini bertanggung jawab untuk mengelola presensi siswa dan guru.

## Fitur

- **Presensi Siswa**: Check-in, check-out, status (hadir, sakit, izin, alpa).
- **Presensi Guru**: Check-in, check-out, lokasi.
- **Bulk Check-in**: Presensi massal untuk satu kelas.
- **Validasi GPS**: Validasi lokasi saat check-in (untuk siswa/guru).
- **Laporan**: Rekap presensi per kelas/periode.

## Struktur Project

```
services/attendance-service/
├── cmd/
│   └── server/          # Entry point aplikasi
├── internal/
│   ├── domain/          # Entities & interfaces
│   ├── handler/         # HTTP handlers
│   ├── repository/      # Database access
│   └── usecase/         # Business logic
└── migrations/          # Database schemas
```

## API Endpoints

### Student Attendance

- `POST /api/v1/attendance/students` - Create student attendance
- `POST /api/v1/attendance/students/bulk` - Bulk create student attendance
- `GET /api/v1/attendance/students` - Get attendance by class & date
- `GET /api/v1/attendance/students/:id` - Get detail attendance
- `GET /api/v1/attendance/students/:id/summary` - Get attendance summary
- `PUT /api/v1/attendance/students/:id` - Update attendance

### Teacher Attendance

- `POST /api/v1/attendance/teachers/checkin` - Teacher check-in
- `PUT /api/v1/attendance/teachers/checkout` - Teacher check-out
- `GET /api/v1/attendance/teachers` - Get teacher attendance list

### Reports

- `GET /api/v1/attendance/reports/daily` - Get daily attendance report
- `GET /api/v1/attendance/reports/monthly` - Get monthly attendance report
- `GET /api/v1/attendance/reports/class/:class_id` - Get class attendance report

## Testing

### Integration Tests

To run the integration tests, ensure you have Docker installed and running. The tests will spin up a PostgreSQL container.

```bash
# Run all integration tests
go test -v ./cmd/server/...
```

## Konfigurasi

Environment variables yang dibutuhkan:

```env
APP_HTTP_PORT=9093
DB_URL=postgres://user:pass@host:port/dbname
ACADEMIC_SERVICE_URL=http://academic-service:9091
```
