class LyknError(Exception):
    """Base exception for lykn."""


class LicenseFileError(LyknError):
    """License file read or format error."""


class LicenseSignatureError(LyknError):
    """Signature verification failed."""


class LicenseExpiredError(LyknError):
    """License has expired."""


class LicenseNotYetValidError(LyknError):
    """License is not yet valid."""


class HardwareMismatchError(LyknError):
    """Hardware fingerprint mismatch."""


class FeatureNotLicensedError(LyknError):
    """Requested feature is not licensed."""
