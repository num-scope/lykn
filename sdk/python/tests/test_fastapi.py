from datetime import datetime, timedelta, timezone

from fastapi import Depends, FastAPI
from fastapi.testclient import TestClient

from lykn.contrib.fastapi import RequireFeatures, RequireLicense
from lykn.exceptions import FeatureNotLicensedError, LicenseExpiredError
from lykn.schemas import LicenseData, Subject


class StubValidator:
    def __init__(self, *, error: Exception | None = None):
        self.error = error
        now = datetime.now(timezone.utc)
        self.license = LicenseData(
            id="lic-001",
            subject=Subject(name="Demo"),
            plan="pro",
            issued_at=now,
            not_before=now - timedelta(minutes=1),
            not_after=now + timedelta(days=1),
            features=["reports", "exports"],
        )

    def verify(self, required_features=None):
        if self.error:
            raise self.error
        if required_features:
            missing = [item for item in required_features if item not in self.license.features]
            if missing:
                raise FeatureNotLicensedError(f"Missing licensed features: {', '.join(missing)}")
        return self.license


def test_require_license_returns_license_data():
    validator = StubValidator()
    app = FastAPI()

    @app.get("/protected")
    def protected(license_data=Depends(RequireLicense(validator))):
        return {"plan": license_data.plan}

    response = TestClient(app).get("/protected")

    assert response.status_code == 200
    assert response.json() == {"plan": "pro"}


def test_require_license_translates_errors_to_403():
    validator = StubValidator(error=LicenseExpiredError("License has expired"))
    app = FastAPI()

    @app.get("/protected")
    def protected(license_data=Depends(RequireLicense(validator))):
        return {"plan": license_data.plan}

    response = TestClient(app).get("/protected")

    assert response.status_code == 403
    assert response.json()["detail"] == "License has expired"


def test_require_features_blocks_missing_feature():
    validator = StubValidator()
    app = FastAPI()

    @app.get("/reports")
    def reports(license_data=Depends(RequireFeatures(validator, "reports", "billing"))):
        return {"plan": license_data.plan}

    response = TestClient(app).get("/reports")

    assert response.status_code == 403
    assert response.json()["detail"] == "Missing licensed features: billing"
