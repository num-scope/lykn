from lykn.hardware import collect_hardware, normalize_hardware
from lykn.schemas import Hardware


def test_normalize_hardware_normalizes_and_sorts_values():
    raw = Hardware(
        mac_addresses=["aa-bb-cc-dd-ee-ff", "AA:BB:CC:DD:EE:FF"],
        ip_addresses=["192.168.1.10", "10.0.0.1", "192.168.1.10"],
        hostname="  Demo Host  ",
        cpu_id="  CPU-1  ",
        disk_serial="  DISK-1  ",
    )

    normalized = normalize_hardware(raw)

    assert normalized.mac_addresses == ["AA:BB:CC:DD:EE:FF"]
    assert normalized.ip_addresses == ["10.0.0.1", "192.168.1.10"]
    assert normalized.hostname == "Demo Host"
    assert normalized.cpu_id == "CPU-1"
    assert normalized.disk_serial == "DISK-1"


def test_collect_hardware_builds_hardware_from_getters(monkeypatch):
    monkeypatch.setattr("lykn.hardware.get_mac_addresses", lambda: ["aa-bb-cc-dd-ee-ff"])
    monkeypatch.setattr("lykn.hardware.get_ip_addresses", lambda: ["192.168.1.20", "192.168.1.20"])
    monkeypatch.setattr("lykn.hardware.get_hostname", lambda: " demo-host ")
    monkeypatch.setattr("lykn.hardware.get_cpu_id", lambda: " cpu-123 ")
    monkeypatch.setattr("lykn.hardware.get_disk_serial", lambda: " disk-123 ")

    hardware = collect_hardware()

    assert hardware == Hardware(
        mac_addresses=["AA:BB:CC:DD:EE:FF"],
        ip_addresses=["192.168.1.20"],
        hostname="demo-host",
        cpu_id="cpu-123",
        disk_serial="disk-123",
    )


def test_collect_hardware_swallows_optional_getter_errors(monkeypatch):
    monkeypatch.setattr("lykn.hardware.get_mac_addresses", lambda: ["AA:BB:CC:DD:EE:FF"])
    monkeypatch.setattr("lykn.hardware.get_ip_addresses", lambda: ["10.0.0.8"])
    monkeypatch.setattr("lykn.hardware.get_hostname", lambda: "demo-host")

    def boom():
        raise RuntimeError("not supported")

    monkeypatch.setattr("lykn.hardware.get_cpu_id", boom)
    monkeypatch.setattr("lykn.hardware.get_disk_serial", boom)

    hardware = collect_hardware()

    assert hardware.cpu_id == ""
    assert hardware.disk_serial == ""
