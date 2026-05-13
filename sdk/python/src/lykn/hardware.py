import platform
import socket
import subprocess
import uuid
from ipaddress import ip_address

import psutil

from lykn.schemas import Hardware


def _safe_call(func, default):
    try:
        value = func()
    except Exception:
        return default
    return default if value is None else value


def _normalize_mac(value: str) -> str:
    text = value.strip().replace("-", ":").upper()
    parts = [part.zfill(2) for part in text.split(":") if part]
    return ":".join(parts) if len(parts) == 6 else ""


def _first_non_empty_line(text: str) -> str:
    for line in text.splitlines():
        line = line.strip()
        if line and "serial" not in line.lower() and "processorid" not in line.lower():
            return line
    return ""


def _run_command(command: list[str]) -> str:
    completed = subprocess.run(command, capture_output=True, check=True, text=True)
    return completed.stdout.strip()


def get_mac_addresses() -> list[str]:
    macs: set[str] = set()

    for interface_addrs in psutil.net_if_addrs().values():
        for address in interface_addrs:
            normalized = _normalize_mac(str(address.address))
            if normalized and normalized != "00:00:00:00:00:00":
                macs.add(normalized)

    fallback = uuid.getnode()
    fallback_mac = ":".join(f"{(fallback >> shift) & 0xFF:02X}" for shift in range(40, -1, -8))
    if fallback_mac != "00:00:00:00:00:00":
        macs.add(fallback_mac)

    return sorted(macs)


def get_ip_addresses() -> list[str]:
    addresses: set[str] = set()
    for interface_addrs in psutil.net_if_addrs().values():
        for address in interface_addrs:
            raw = str(address.address).split("%", 1)[0]
            try:
                parsed = ip_address(raw)
            except ValueError:
                continue
            if parsed.is_loopback:
                continue
            addresses.add(str(parsed))
    return sorted(addresses)


def get_hostname() -> str:
    return socket.gethostname().strip()


def get_cpu_id() -> str:
    system = platform.system().lower()
    if system == "linux":
        return _run_command(["sh", "-c", "cat /proc/cpuinfo | awk -F': ' '/Serial/{print $2; exit}'"])
    if system == "darwin":
        return _run_command(["sysctl", "-n", "machdep.cpu.brand_string"])
    if system == "windows":
        return _first_non_empty_line(_run_command(["wmic", "cpu", "get", "ProcessorId"]))
    return ""


def get_disk_serial() -> str:
    system = platform.system().lower()
    if system == "linux":
        return _run_command(["sh", "-c", "lsblk -ndo SERIAL | head -n 1"])
    if system == "darwin":
        return _run_command(
            ["sh", "-c", "system_profiler SPNVMeDataType | awk -F': ' '/Serial Number/{print $2; exit}'"]
        )
    if system == "windows":
        return _first_non_empty_line(_run_command(["wmic", "diskdrive", "get", "SerialNumber"]))
    return ""


def normalize_hardware(hw: Hardware) -> Hardware:
    mac_addresses = sorted({mac for mac in (_normalize_mac(item) for item in hw.mac_addresses) if mac})
    ip_addresses = sorted({item.strip() for item in hw.ip_addresses if item.strip()})
    return Hardware(
        mac_addresses=mac_addresses,
        ip_addresses=ip_addresses,
        hostname=hw.hostname.strip(),
        cpu_id=hw.cpu_id.strip(),
        disk_serial=hw.disk_serial.strip(),
    )


def collect_hardware() -> Hardware:
    hardware = Hardware(
        mac_addresses=_safe_call(get_mac_addresses, []),
        ip_addresses=_safe_call(get_ip_addresses, []),
        hostname=_safe_call(get_hostname, ""),
        cpu_id=_safe_call(get_cpu_id, ""),
        disk_serial=_safe_call(get_disk_serial, ""),
    )
    return normalize_hardware(hardware)
