from pathlib import Path
from typing import Literal

from pydantic import BaseModel, Field, model_validator


EngineName = Literal["pyinstaller", "nuitka"]
BundleMode = Literal["onedir", "onefile"]


class ResourceSpec(BaseModel):
    source: Path
    target: str | None = None


class PackConfig(BaseModel):
    project_dir: Path
    engine: EngineName
    entry_script: Path | None = None
    entry_module: str | None = None
    output_dir: Path = Path("dist")
    bundle_mode: BundleMode = "onedir"
    resources: list[ResourceSpec] = Field(default_factory=list)
    exclude: list[str] = Field(default_factory=list)
    name: str | None = None
    clean: bool = True
    extra_args: list[str] = Field(default_factory=list)

    @model_validator(mode="after")
    def validate_entry(self):
        if (self.entry_script is None) == (self.entry_module is None):
            raise ValueError("exactly one of entry_script or entry_module must be provided")
        return self


class StagingResult(BaseModel):
    staging_dir: Path
    resources: list[ResourceSpec] = Field(default_factory=list)


class PackResult(BaseModel):
    engine: EngineName
    command: list[str]
    output_path: Path
    resources: list[ResourceSpec] = Field(default_factory=list)
