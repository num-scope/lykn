from datetime import datetime, timezone
from pathlib import Path
from threading import Event, Thread
from typing import Callable

from lykn.crypto import load_public_key, verify_license_file
from lykn.exceptions import (
    FeatureNotLicensedError,
    HardwareMismatchError,
    LicenseExpiredError,
    LicenseFileError,
    LicenseNotYetValidError,
)
from lykn.hardware import collect_hardware, normalize_hardware
from lykn.schemas import Hardware, LicenseData


class LicenseValidator:
    def __init__(
        self,
        public_key: str | Path,
        license_path: str | Path | None = None,
        license_content: str | bytes | None = None,
        check_interval: int = 0,
    ):
        if license_path is None and license_content is None:
            raise LicenseFileError("license_path or license_content is required")

        self.public_key = public_key
        self.license_path = Path(license_path) if license_path is not None else None
        self.license_content = license_content
        self.check_interval = check_interval

        self.license: LicenseData | None = None
        self._last_error: Exception | None = None
        self._invalid = False
        self._thread: Thread | None = None
        self._stop_event = Event()
        self._on_invalid_callbacks: list[Callable[[str], None]] = []

    def _read_license_content(self) -> str | bytes:
        if self.license_content is not None:
            return self.license_content
        if self.license_path is None:
            raise LicenseFileError("license_path or license_content is required")
        try:
            return self.license_path.read_text()
        except OSError as exc:
            raise LicenseFileError(f"Unable to read license file: {exc}") from exc

    def _as_utc(self, value: datetime) -> datetime:
        if value.tzinfo is None:
            return value.replace(tzinfo=timezone.utc)
        return value.astimezone(timezone.utc)

    def _validate_time(self, data: LicenseData) -> None:
        now = datetime.now(timezone.utc)
        not_before = self._as_utc(data.not_before)
        not_after = self._as_utc(data.not_after)
        if not_before > now:
            raise LicenseNotYetValidError("License is not yet valid")
        if not_after < now:
            raise LicenseExpiredError("License has expired")

    def _validate_hardware(self, expected: Hardware | None) -> None:
        if expected is None:
            return

        current = normalize_hardware(collect_hardware())
        required = normalize_hardware(expected)

        if required.hostname and required.hostname != current.hostname:
            raise HardwareMismatchError("Hardware hostname mismatch")
        if required.cpu_id and required.cpu_id != current.cpu_id:
            raise HardwareMismatchError("Hardware CPU id mismatch")
        if required.disk_serial and required.disk_serial != current.disk_serial:
            raise HardwareMismatchError("Hardware disk serial mismatch")
        if required.mac_addresses and not set(required.mac_addresses).issubset(current.mac_addresses):
            raise HardwareMismatchError("Hardware MAC addresses mismatch")
        if required.ip_addresses and not set(required.ip_addresses).issubset(current.ip_addresses):
            raise HardwareMismatchError("Hardware IP addresses mismatch")

    def _validate_features(self, data: LicenseData, required_features: list[str] | None) -> None:
        if not required_features:
            return
        missing = [feature for feature in required_features if feature not in data.features]
        if missing:
            raise FeatureNotLicensedError(f"Missing licensed features: {', '.join(missing)}")

    def _verify_once(self, required_features: list[str] | None = None) -> LicenseData:
        public_key = load_public_key(self.public_key)
        payload = verify_license_file(self._read_license_content(), public_key)
        data = LicenseData.model_validate(payload)
        self._validate_time(data)
        self._validate_hardware(data.hardware)
        self._validate_features(data, required_features)
        return data

    def verify(self, required_features: list[str] | None = None) -> LicenseData:
        data = self._verify_once(required_features)
        self.license = data
        self._last_error = None
        self._invalid = False
        return data

    def has_feature(self, feature: str) -> bool:
        if self.license is None:
            self.verify()
        return feature in (self.license.features if self.license else [])

    def _mark_invalid(self, exc: Exception) -> None:
        self._invalid = True
        self._last_error = exc
        for callback in self._on_invalid_callbacks:
            callback(str(exc))

    def _run_loop(self) -> None:
        while not self._stop_event.wait(self.check_interval):
            try:
                self.verify()
            except Exception as exc:
                self._mark_invalid(exc)
                break

    def start(self) -> LicenseData:
        result = self.verify()
        if self.check_interval <= 0:
            return result
        if self._thread and self._thread.is_alive():
            return result

        self._stop_event.clear()
        self._thread = Thread(target=self._run_loop, daemon=True)
        self._thread.start()
        return result

    def stop(self) -> None:
        self._stop_event.set()
        if self._thread and self._thread.is_alive():
            self._thread.join(timeout=self.check_interval + 1)

    def on_invalid(self, func):
        self._on_invalid_callbacks.append(func)
        return func
