from fastapi import HTTPException

from lykn.exceptions import LyknError


class RequireLicense:
    def __init__(self, validator):
        self.validator = validator

    def __call__(self):
        try:
            return self.validator.verify()
        except LyknError as exc:
            raise HTTPException(status_code=403, detail=str(exc)) from exc


class RequireFeatures:
    def __init__(self, validator, *features: str):
        self.validator = validator
        self.features = list(features)

    def __call__(self):
        try:
            return self.validator.verify(required_features=self.features)
        except LyknError as exc:
            raise HTTPException(status_code=403, detail=str(exc)) from exc
