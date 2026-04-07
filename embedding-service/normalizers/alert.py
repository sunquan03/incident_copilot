def normalize_alert(event: dict) -> tuple[str, dict]:
    labels = event.get("labels") or {}
    label_str = " ".join(f"{k}={v}" for k, v in labels.items())
    text = f"alert {event.get('source_name', '')} {label_str} {event.get('message', '')}".strip()

    payload = {
        "event_type": "alert",
        "source_id": event.get("source_id"),
        "source_name": event.get("source_name"),
        "message": event.get("message"),
        "labels": labels,
        "severity": labels.get("severity"),
        "service": labels.get("service"),
        "region": labels.get("region"),
        "created_at": event.get("created_at"),
    }
    return text, payload

