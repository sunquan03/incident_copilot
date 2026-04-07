def normalize_incident(event: dict) -> tuple[str, dict]:
    tags = event.get("tags") or []
    tag_str = " ".join(tags)
    text = f"incident {event.get('service_name', '')} {tag_str} {event.get('message', '')}".strip()

    payload = {
        "event_type": "incident",
        "source_id": event.get("source_id"),
        "service_name": event.get("service_name"),
        "message": event.get("message"),
        "tags": tags,
        "status": event.get("status"),
        "created_at": event.get("created_at"),
    }
    return text, payload