#!/usr/bin/env python3
"""
Demonit device agent (Windows / macOS / Linux).

Zero required dependencies (stdlib only). Optional: pip install psutil
for more accurate CPU/RAM on all platforms.

Examples:
  python device_agent.py --device-id <UUID>
  python device_agent.py --device-id <UUID> --api-base http://192.168.1.10:8080
  python device_agent.py --device-id <UUID> --api-key your-device-key --interval 10
"""

from __future__ import annotations

import argparse
import json
import os
import platform
import socket
import sys
import time
import urllib.error
import urllib.request
from typing import Any


def read_cpu_ram() -> tuple[float, float]:
    """Return (cpu_percent, ram_percent). Best-effort across OSes."""
    try:
        import psutil  # type: ignore

        cpu = float(psutil.cpu_percent(interval=0.3))
        ram = float(psutil.virtual_memory().percent)
        return round(cpu, 1), round(ram, 1)
    except Exception:
        pass

    system = platform.system().lower()
    if system == "linux":
        return _linux_cpu_ram()
    if system == "darwin":
        return _darwin_cpu_ram()
    if system == "windows":
        return _windows_cpu_ram()
    return 10.0, 40.0


def _linux_cpu_ram() -> tuple[float, float]:
    # CPU: sample /proc/stat twice
    def snap() -> tuple[int, int]:
        with open("/proc/stat", encoding="utf-8") as f:
            parts = f.readline().split()
        vals = list(map(int, parts[1:]))
        idle = vals[3] + (vals[4] if len(vals) > 4 else 0)
        total = sum(vals)
        return idle, total

    idle1, total1 = snap()
    time.sleep(0.2)
    idle2, total2 = snap()
    didle = idle2 - idle1
    dtotal = total2 - total1
    cpu = 0.0 if dtotal <= 0 else (1.0 - didle / dtotal) * 100.0

    mem_total = mem_avail = 0
    with open("/proc/meminfo", encoding="utf-8") as f:
        for line in f:
            if line.startswith("MemTotal:"):
                mem_total = int(line.split()[1])
            elif line.startswith("MemAvailable:"):
                mem_avail = int(line.split()[1])
    ram = 0.0 if mem_total <= 0 else (1.0 - mem_avail / mem_total) * 100.0
    return round(cpu, 1), round(ram, 1)


def _darwin_cpu_ram() -> tuple[float, float]:
    # Lightweight fallbacks via shell helpers already on macOS.
    import subprocess

    cpu = 15.0
    ram = 50.0
    try:
        out = subprocess.check_output(["ps", "-A", "-o", "%cpu"], text=True)
        vals = []
        for line in out.splitlines()[1:]:
            line = line.strip()
            if line:
                vals.append(float(line))
        if vals:
            # Rough aggregate; capped for display sanity.
            cpu = min(100.0, sum(vals) / max(os.cpu_count() or 1, 1))
    except Exception:
        pass
    try:
        out = subprocess.check_output(["vm_stat"], text=True)
        page_size = 4096
        stats: dict[str, int] = {}
        for line in out.splitlines():
            if ":" not in line:
                continue
            key, raw = line.split(":", 1)
            digits = "".join(ch for ch in raw if ch.isdigit())
            if digits:
                stats[key.strip()] = int(digits)
        free = stats.get("Pages free", 0) + stats.get("Pages speculative", 0)
        active = stats.get("Pages active", 0)
        inactive = stats.get("Pages inactive", 0)
        wired = stats.get("Pages wired down", 0)
        used = active + inactive + wired
        total = used + free
        if total > 0:
            ram = used / total * 100.0
        _ = page_size  # reserved for future byte math
    except Exception:
        pass
    return round(cpu, 1), round(ram, 1)


def _windows_cpu_ram() -> tuple[float, float]:
    import subprocess

    cpu = 20.0
    ram = 50.0
    try:
        out = subprocess.check_output(
            [
                "powershell",
                "-NoProfile",
                "-Command",
                "(Get-CimInstance Win32_Processor | Measure-Object LoadPercentage -Average).Average",
            ],
            text=True,
        ).strip()
        if out:
            cpu = float(out)
    except Exception:
        pass
    try:
        out = subprocess.check_output(
            [
                "powershell",
                "-NoProfile",
                "-Command",
                "$o=Get-CimInstance Win32_OperatingSystem; [math]::Round((($o.TotalVisibleMemorySize-$o.FreePhysicalMemory)/$o.TotalVisibleMemorySize)*100,1)",
            ],
            text=True,
        ).strip()
        if out:
            ram = float(out)
    except Exception:
        pass
    return round(cpu, 1), round(ram, 1)


def post_heartbeat(api_base: str, api_key: str, body: dict[str, Any]) -> dict[str, Any]:
    url = api_base.rstrip("/") + "/api/v1/heartbeat"
    data = json.dumps(body).encode("utf-8")
    req = urllib.request.Request(
        url,
        data=data,
        method="POST",
        headers={
            "Content-Type": "application/json",
            "Accept": "application/json",
            "X-Device-API-Key": api_key,
        },
    )
    with urllib.request.urlopen(req, timeout=10) as resp:
        return json.loads(resp.read().decode("utf-8"))


def main() -> int:
    parser = argparse.ArgumentParser(description="Demonit cross-platform device agent")
    parser.add_argument("--device-id", required=True, help="Device UUID from the dashboard/API")
    parser.add_argument(
        "--api-base",
        default=os.getenv("DEMONIT_API_BASE", "http://localhost:8080"),
        help="Backend base URL (default: http://localhost:8080)",
    )
    parser.add_argument(
        "--api-key",
        default=os.getenv("DEMONIT_DEVICE_API_KEY", "dev-device-key-change-me"),
        help="X-Device-API-Key value",
    )
    parser.add_argument("--interval", type=int, default=10, help="Seconds between heartbeats")
    args = parser.parse_args()

    print("Demonit agent started")
    print(f"  device_id : {args.device_id}")
    print(f"  api       : {args.api_base}")
    print(f"  interval  : {args.interval}s")
    print("Ctrl+C to stop.\n")

    while True:
        cpu, ram = read_cpu_ram()
        body = {
            "device_id": args.device_id,
            "cpu_usage": cpu,
            "ram_usage": ram,
            "status_payload": {
                "hostname": socket.gethostname(),
                "os": platform.system().lower(),
                "os_release": platform.release(),
                "arch": platform.machine(),
                "agent": "python",
                "python": platform.python_version(),
            },
        }
        ts = time.strftime("%H:%M:%S")
        try:
            res = post_heartbeat(args.api_base, args.api_key, body)
            status = (res.get("data") or {}).get("status", "?")
            print(f"[{ts}] OK  cpu={cpu}% ram={ram}%  status={status}")
        except urllib.error.HTTPError as e:
            detail = e.read().decode("utf-8", errors="replace")
            print(f"[{ts}] FAIL HTTP {e.code}: {detail}", file=sys.stderr)
        except Exception as e:
            print(f"[{ts}] FAIL {e}", file=sys.stderr)
        time.sleep(max(args.interval, 1))


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except KeyboardInterrupt:
        print("\nStopped.")
        raise SystemExit(0)
