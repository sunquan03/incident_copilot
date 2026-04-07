def normalize_logdoc(event: dict) -> tuple[str, dict]:
    tags = event.get("tags") or []
    tag_str = " ".join(tags)
    text = f"logdoc {event.get('source_type', '')} {event.get('service_name', '')} {tag_str} {event.get('title', '')} {event.get('content', '')}".strip()

    payload = {
        "event_type": "logdoc",
        "source_id": event.get("source_id"),
        "service_name": event.get("service_name"),
        "title": event.get("title"),
        "content": event.get("content"),
        "source_type": event.get("source_type"),
        "tags": tags,
        "status": event.get("status"),
        "created_at": event.get("created_at"),
    }
    return text, payload


