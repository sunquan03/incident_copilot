from .alert import normalize_alert
from .incident import normalize_incident
from .logdoc import normalize_logdoc

NORMALIZERS = {
    "alert.received": normalize_alert,
    "incident.created": normalize_incident,
    "logdoc.received": normalize_logdoc,
}


__all__ = ["NORMALIZERS"]
