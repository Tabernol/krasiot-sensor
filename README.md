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
  - Subscribes to MQTT (HiveMQ), enriches readings, saves to Oracle Autonomous DB (ADB), forwards events to AWS SQS.
  - Exposes **GET** `/latest` on port **8080** for the most recent message (test endpoint).
- `krasiot-notifier` (Go) — **running (MVP)**
  - Consumes messages from AWS SQS, checks latest sent notification in DB, sends email if changes are important.
  - Connected to ADB for state checks.
- `krasiot-ui` (TypeScript/React) — **planned**
  - Dashboard for latest readings, device list, alert history, basic settings.
- `krasiot-auth` (Java/Spring Boot) — **planned**
  - Central auth (JWT), roles, basic tenant separation.

**Infrastructure**
- CI/CD: GitHub Actions (adjusted; currently building binary on VM).
- Cloud:
  - **OCI** VM `Standard.A1.Flex` (ARM) for running Go services.
  - **Oracle Autonomous DB** (ADB) for storage (connected via **godror**, Wallet, Instant Client `19.23`).
  - **AWS SQS** for event queueing (notifications pipeline).
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
                      |                 |
                      | (godror)        | (SQS send)
                      v                 v
                  [Oracle ADB] <---- [AWS SQS]
                                          |
                                          v
                                  [krasiot-notifier]
                                    (dedupe + rules)
                                          |
                                          v
                                    [Email Sender]

Optional / Next: [krasiot-ui] (React) -> read APIs -> ADB (via service) / alerts
                 [krasiot-auth] (JWT) -> protect APIs/UI
```

**Data Flow (today)**
1. Device publishes reading every **5 minutes** to HiveMQ.
2. `krasiot-sensor` subscribes, normalizes/enriches, **persists to ADB**, and **pushes to SQS**.
3. `krasiot-notifier` consumes SQS, checks last notification in ADB, and sends an email when rules trigger.

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

### 4.3 Enriched Event (JSON) — Outbound to SQS (Proposed)
```json
{
  "event_id": "7f2c7c20-1c7f-4ff0-9c2f-3a6a2ec8b1b9",
  "hardware_uid": "esp32abcd",
  "ts_ms": 1725600000000,
  "soil_moisture_raw": 612,
  "soil_moisture_category": "DRY",  
  "temperature_c": 22.4,
  "humidity_pct": 54.2,
  "ingest_src": "mqtt/hivemq",
  "ingest_seq": 145903,
  "env": "prod"
}
```
> `soil_moisture_category` thresholds to be configured per device/profile.

### 4.4 REST — `krasiot-sensor` Test Endpoint
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

## 5) Notification Logic (Current/Planned)

**Business Rules (target)**
- Send a notification **if**:
  1) **No previous** notification exists for the `hardware_uid`, **or**
  2) `soil_moisture_category` **changed** and **≥ 30 minutes** have passed since the last notification, **or**
  3) `soil_moisture_category` **unchanged** but **≥ 24 hours** have passed since the last notification.

**Email Content (MVP)**
- Subject: `krasiot alert — {hardware_uid} — {soil_moisture_category}`
- Body: short status summary + last reading + link to dashboard (future).

**Delivery**
- SMTP or AWS SES (TBD). Retry with backoff. Store send outcome in ADB.

---

## 6) Storage Model (Proposed, ADB)

**Tables**
- `readings` — immutable time‑series
  - (`id` PK, `hardware_uid`, `ts`, `soil_moisture_raw`, `soil_moisture_category`, `temperature_c`, `humidity_pct`, `ingest_src`, `ingest_seq`, `env`)
- `alerts` — sent notifications
  - (`id` PK, `hardware_uid`, `ts`, `category`, `subject`, `status`, `meta_json`)
- `devices` — device metadata/config
  - (`hardware_uid` PK, `name`, `owner_id`, `thresholds_json`, `created_at`, `status`)

**Indexes**
- `readings(hardware_uid, ts)`
- `alerts(hardware_uid, ts)`

**Retention (Proposed)**
- Raw readings: 180–365 days (TBD).
- Alerts: 365+ days (TBD).

---

## 7) Configuration & Secrets

**Environment variables (examples)**
- MQTT: `MQTT_BROKER_URI`, `MQTT_USER`, `MQTT_PASS`, `MQTT_TOPIC`
- DB: `DB_USER`, `DB_PASS`, `DB_WALLET_DIR`, `DB_CONNECT_STR`
- SQS: `AWS_REGION`, `SQS_QUEUE_URL`
- Email: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` (or SES params)
- App: `ENV`, `LOG_LEVEL`, `NOTIFY_COOLDOWN_MIN=30`, `NOTIFY_DAILY_HEARTBEAT_H=24`

**Secrets handling**: GitHub Actions Secrets; consider OCI Vault/AWS Secrets Manager later.

---

## 8) Deployment

**Today**
- Build artifacts on OCI VM (ARM) and run services directly.

**Next (Proposed)**
- Cross‑compile via GitHub Actions (Linux/ARM64) and upload artifacts.
- Systemd services or Docker (single host) for better process mgmt.
- Versioned releases, health checks, log rotation.

---

## 9) Observability

- Logging: structured logs to STDOUT + file (rotation).
- Metrics (Proposed): Expose Prometheus endpoint for sensor/queue DB metrics.
- Dashboards (Proposed): Grafana (self‑hosted or managed) for key graphs.
- Alerts (Proposed): On dead letter rate, DB errors, MQTT disconnects.

---

## 10) Security

- Use Oracle Wallet for ADB (already in place).
- Restrict inbound ports on OCI VM; only required HTTP ports exposed.
- Rotate credentials; least‑privilege AWS IAM for SQS consumer.
- JWT‑based auth via `krasiot-auth` (future) for UI/API.

---

## 11) Roadmap (Next 4–6 Weeks)

1. **Stabilize ingestion**: backpressure, retry policies, idempotency keys.
2. **Notifier rules**: finalize thresholds & schedules; implement daily heartbeat.
3. **CI/CD**: move builds to GitHub Actions (ARM64), produce signed artifacts.
4. **Packaging**: systemd units or Docker; env‑file based config.
5. **UI MVP**: latest reading per device, alert list, simple chart (24h).
6. **Auth MVP**: single‑tenant JWT + password login.
7. **Docs**: runbooks (on‑call), infra as code (Terraform — minimal).

---

## 12) Open Questions (Need Input)

- Exact MQTT topic(s) and payload from the **prebuilt** prototype (fields, units)?
- Final thresholds for `soil_moisture_category` by device or global?
- Email provider choice (SMTP vs SES) and sender domain?
- AWS region & SQS queue names/URLs used in prod?
- ADB schema names and any constraints already created?
- Desired data retention window (raw vs. aggregated)?
- Are temperature/humidity required for notifications or only moisture?
- Will UI be read‑only initially, or include device threshold editing?
- Device provisioning: how to register/approve new `hardware_uid`?

---

## 13) Glossary

- **ADB**: Oracle Autonomous Database.
- **HiveMQ**: Managed MQTT broker.
- **MQTT**: Lightweight pub/sub protocol for IoT.
- **SQS**: AWS Simple Queue Service.

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

### B. Basic Health Checklist
- [ ] VM has Oracle Instant Client `19.23` & Wallet configured
- [ ] `krasiot-sensor` connected to MQTT and ADB, can write `readings`
- [ ] SQS messages produced per reading
- [ ] `krasiot-notifier` consumes SQS, writes to `alerts`, sends email
- [ ] `/latest` returns the most recent enriched reading

---

*This is a living document. Edit inline as the system evolves.*

