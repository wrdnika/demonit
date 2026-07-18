# Demonit device agents

Send heartbeats from any machine to your Demonit backend.

## 1. Register a device

From the dashboard (**Register device**), or via API:

```bash
curl -X POST http://localhost:8080/api/v1/devices \
  -H "Content-Type: application/json" \
  -H "X-Admin-API-Key: YOUR_ADMIN_KEY" \
  -d '{"name":"My Laptop","type":"LAPTOP"}'
```

Copy the returned `id` (UUID).

## 2. Pick an agent

| Platform | Script | Notes |
|----------|--------|--------|
| **Any (recommended)** | `device_agent.py` | Windows / macOS / Linux, stdlib only |
| **Windows** | `device-agent.ps1` | Native PowerShell metrics |
| **macOS / Linux** | `device-agent.sh` | Wraps Python if available, else curl fallback |

Default API key for local dev: `dev-device-key-change-me`  
(change `DEVICE_API_KEY` in backend `.env` for real deployments)

## 3. Run

### Python (recommended)

```bash
# Windows / macOS / Linux
python scripts/device_agent.py --device-id YOUR-UUID

# Point at a remote backend on your LAN
python scripts/device_agent.py \
  --device-id YOUR-UUID \
  --api-base http://192.168.1.10:8080 \
  --api-key YOUR_DEVICE_API_KEY \
  --interval 10
```

Optional better metrics:

```bash
pip install psutil
```

### Windows PowerShell

```powershell
cd scripts
.\device-agent.ps1 -DeviceId "YOUR-UUID"
# or from Git Bash:
powershell.exe -ExecutionPolicy Bypass -File ./device-agent.ps1 -DeviceId "YOUR-UUID"
```

### macOS / Linux shell

```bash
chmod +x scripts/device-agent.sh
./scripts/device-agent.sh --device-id YOUR-UUID --api-base http://localhost:8080
```

## 4. What you should see

- Dashboard device flips to **ONLINE**
- CPU / RAM / `status_payload` update each interval
- Stop the agent → after ~30s deadman marks **OFFLINE** (SSE alert)

## 5. Network tips

- `localhost` only works on the same machine as the API
- Phones / Pi / other PCs must use your computer’s LAN IP, e.g. `http://192.168.x.x:8080`
- Backend listens on `0.0.0.0:8080` by default so LAN devices can connect
- Firewall must allow inbound TCP `8080`
