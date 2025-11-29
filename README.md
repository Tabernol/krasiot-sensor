# krasiot — Project Overview

> **Elevator pitch**: krasiot is a lightweight, low‑cost IoT platform to monitor and optimize agricultural environments (soil moisture, temperature, humidity) using low‑power sensor devices. Initial scope targets small farms and home gardens with Always‑Free/low‑cost infrastructure.

---

## 1) Goals & Non‑Goals

**Goals**
- Minimal infrastructure cost (prefer Always‑Free tiers).
- Reliable ingestion of sensor data (MQTT ➜ processing ➜ database).
- Near‑real‑time notifications when conditions change materially.
- Simple web UI for visibility (latest reading, trends, alerts).
- Clear path to authentication/authorization for multi‑tenant use.

**Non‑Goals (for now)**
- Heavy analytics/ML.
- Over‑engineered Kubernetes setup.
- Device OTA management beyond basic config.

---

## 2) Current State

**Hardware / Devices**
- **Prebuilt prototype**: ESP32 + DHT11 + soil moisture + CP2104 + 18650 holder. Link (AliExpress): https://www.aliexpress.com/item/32842847829.html#nav-specification
- **Custom option**: ESP32‑S3 Mini + capacitive soil‑moisture sensor (planned/partial).

**Services (microservices)**
- `krasiot-sensor` (Go) — **running (MVP)**
  - Subscribes to MQTT (HiveMQ), enriches readings (calculates moisture percentage and category), and persists to Oracle Autonomous DB (ADB).
  - Exposes **GET** `/latest` on port **8080** for the most recent message (test endpoint).
  - **Responsibility:** Data ingestion and persistence only. Does not handle notifications.
- `krasiot-notifier` (Go) — **running (MVP)**
  - Polls `sensor_raw` table in ADB, applies notification rules, sends emails.
  - **Responsibility:** Notification logic and email delivery.
- `krasiot-ui` (TypeScript/React) — **planned**
  - Dashboard for latest readings, device list, alert history, basic settings.
- `krasiot-auth` (Java/Spring Boot) — **planned**
  - Central auth (JWT), roles, basic tenant separation.

**Infrastructure**
- CI/CD: GitHub Actions (adjusted; currently building binary on VM).
- Cloud:
  - **OCI** VM `Standard.A1.Flex` (ARM) for running Go services.
  - **Oracle Autonomous DB** (ADB) for storage (connected via **godror**, Wallet, Instant Client `19.23`).
  - **HiveMQ** (MQTT broker) for ingest.

**Runtime Versions on VM**
- `go version go1.23.10 linux/arm64`
- Oracle Instant Client: `/usr/lib/oracle/19.23/client64/lib`
- Services on host: `krasiot-sensor`, `krasiot-notifier`

---

## 3) High‑Level Architecture

```
[ESP32 Device(s)] --MQTT--> [HiveMQ]
                              |
                              v
                        [krasiot-sensor]
                   (subscribe, enrich, persist)
                              |
                              | (godror)
                              v
                        [Oracle ADB]
                      (sensor_raw table)

Note: krasiot-notifier (separate service) polls sensor_raw table independently

Optional / Next: [krasiot-ui] (React) -> read APIs -> ADB (via service) / alerts
                 [krasiot-auth] (JWT) -> protect APIs/UI
```

**Data Flow (krasiot-sensor scope)**
1. Device publishes reading every **5 minutes** to HiveMQ.
2. `krasiot-sensor` subscribes to MQTT topic, normalizes/enriches data (calculates moisture %), and **persists to ADB** (`sensor_raw` table).
3. `krasiot-sensor` is stateless and does not know about downstream consumers (notifier, UI, analytics, etc.).

---

## 4) Data & Message Contracts (Proposed)

### 4.1 MQTT Topic Convention (Proposed)
- Topic: `krasiot/{env}/sensors/{hardware_uid}/readings`
  - Example: `krasiot/prod/sensors/esp32abcd/readings`

### 4.2 Sensor Reading (JSON) — Inbound from Device (Proposed)
```json
{
  "hardware_uid": "esp32abcd",
  "ts_ms": 1725600000000,
  "soil_moisture_raw": 612,
  "temperature_c": 22.4,
  "humidity_pct": 54.2
}
```

### 4.3 REST — `krasiot-sensor` Test Endpoint
- `GET /latest` (port 8080)
  - 200 OK example:
```json
{
  "hardware_uid": "esp32abcd",
  "ts_ms": 1725600000000,
  "soil_moisture_raw": 612,
  "soil_moisture_category": "DRY",
  "temperature_c": 22.4,
  "humidity_pct": 54.2
}
```

---

## 5) Notification Logic

> **Note:** Notification logic is handled by `krasiot-notifier` service, not by `krasiot-sensor`.
> `krasiot-sensor` only persists data to `sensor_raw` table. See `krasiot-notifier` documentation for notification rules and email delivery details.

---

## 6) Storage Model (krasiot-sensor scope)

**Tables used by krasiot-sensor**
- `sensor_raw` — immutable time‑series (INSERT only)
  - Contains: `id`, `hardware_uid`, `measured_at_utc`, `ip`, `firmware_version`, `adc_resolution`, `battery_voltage`, `soil_moisture`, `moisture_percent`, `moisture_category`
  - **Responsibility:** krasiot-sensor writes; other services read

**Other tables** (managed by other services)
- `device_alert` — notification history (managed by krasiot-notifier)
- `device` — device metadata/config (future)

**Indexes**
- `sensor_raw(hardware_uid, measured_at_utc)` — for queries by device and time range

**Retention (Proposed)**
- Raw sensor readings: 180–365 days (TBD).

---

## 7) Configuration & Secrets

**Environment variables (krasiot-sensor)**
- MQTT: `MQTT_BROKER_URI`, `MQTT_USER`, `MQTT_PASS`, `MQTT_TOPIC`, `MQTT_CLIENT_ID`
- DB: `DB_USER`, `DB_PASS`, `DB_WALLET_DIR`, `DB_CONNECT_STR`
- App: `ENV`, `LOG_LEVEL`

**Secrets handling**: GitHub Actions Secrets; consider OCI Vault later.

> **Note:** Email and notification configs are handled by `krasiot-notifier` service.

---

## 8) Deployment

**Today**
- Build artifacts on OCI VM (ARM) and run services directly.

**Next (Proposed)**
- Cross‑compile via GitHub Actions (Linux/ARM64) and upload artifacts.
- Systemd services or Docker (single host) for better process mgmt.
- Versioned releases, health checks, log rotation.

---

## 9) Observability (krasiot-sensor)

**Current:**
- Logging: Structured logs to STDOUT (MQTT messages received, DB insert success/failure).

**Proposed:**
- Metrics: Expose Prometheus endpoint for:
  - `mqtt_messages_received_total` (counter)
  - `db_inserts_total{status="success|failure"}` (counter)
  - `db_insert_duration_seconds` (histogram)
  - `mqtt_connection_status` (gauge)
- Dashboards: Grafana dashboard showing message ingestion rate, DB insert latency.
- Alerts: On MQTT disconnects, DB connection failures, high insert error rate (>5%).

---

## 10) Security

- Use Oracle Wallet for ADB (already in place).
- Restrict inbound ports on OCI VM; only required HTTP ports exposed.
- Rotate MQTT and DB credentials regularly.
- JWT‑based auth via `krasiot-auth` (future) for UI/API.

---

## 11) Roadmap (Next 4–6 Weeks)

**krasiot-sensor specific:**
1. **Stabilize ingestion**: Add retry logic for DB inserts, handle MQTT reconnection gracefully.
2. **Monitoring**: Add Prometheus metrics (messages received, DB insert success/failure rate).
3. **CI/CD**: Move builds to GitHub Actions (ARM64), produce signed artifacts.
4. **Packaging**: Systemd unit with auto-restart, env-file based config.
5. **Testing**: Add integration tests for MQTT → DB flow.

**Other services:**
- Notifier rules → see `krasiot-notifier` roadmap
- UI MVP → see `krasiot-ui` roadmap
- Auth MVP → see `krasiot-auth` roadmap

---

## 12) Open Questions (Need Input)

**krasiot-sensor specific:**
- Exact MQTT topic(s) and payload from the **prebuilt** prototype (fields, units)?
- Final thresholds for `soil_moisture_category` by device or global?
- ADB schema names and any constraints already created?
- Desired data retention window for `sensor_raw` table?
- Device provisioning: how to register/approve new `hardware_uid`?

**Other services:**
- Email/notification questions → see `krasiot-notifier` docs
- UI questions → see `krasiot-ui` docs

---

## 13) Glossary

- **ADB**: Oracle Autonomous Database.
- **HiveMQ**: Managed MQTT broker (cloud-hosted).
- **MQTT**: Lightweight pub/sub protocol for IoT.
- **Enrichment**: Process of calculating moisture percentage and category from raw ADC values.

---

## 14) Appendix

### A. Example Systemd Unit (Proposed)
```ini
[Unit]
Description=krasiot-sensor
After=network.target

[Service]
ExecStart=/opt/krasiot/krasiot-sensor
WorkingDirectory=/opt/krasiot
EnvironmentFile=/opt/krasiot/.env
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### B. Basic Health Checklist (krasiot-sensor)
- [ ] VM has Oracle Instant Client `19.23` & Wallet configured
- [ ] `krasiot-sensor` connected to MQTT broker (HiveMQ)
- [ ] `krasiot-sensor` connected to Oracle ADB
- [ ] `krasiot-sensor` can write to `sensor_raw` table
- [ ] `/latest` endpoint returns the most recent enriched reading
- [ ] Logs show successful MQTT message reception and DB inserts

---

*This is a living document. Edit inline as the system evolves.*

