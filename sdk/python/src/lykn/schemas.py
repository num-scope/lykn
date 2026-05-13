from datetime import datetime

from pydantic import BaseModel, Field


class Subject(BaseModel):
    name: str
    email: str = ""
    organization: str = ""


class Hardware(BaseModel):
    mac_addresses: list[str] = Field(default_factory=list)
    ip_addresses: list[str] = Field(default_factory=list)
    hostname: str = ""
    cpu_id: str = ""
    disk_serial: str = ""


class LicenseLimits(BaseModel):
    max_users: int = 0
    max_devices: int = 0


class LicenseData(BaseModel):
    id: str
    version: int = 1
    subject: Subject
    plan: str = ""
    plan_name: str = ""
    issued_at: datetime
    not_before: datetime
    not_after: datetime
    hardware: Hardware | None = None
    features: list[str] = Field(default_factory=list)
    limits: LicenseLimits = Field(default_factory=LicenseLimits)
    metadata: dict = Field(default_factory=dict)
