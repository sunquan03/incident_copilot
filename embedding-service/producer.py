import json
import logging
from confluent_kafka import Producer
from config import Config

logger = logging.getLogger(__name__)

COMPLETED_TOPIC = "embedding.completed"


def build_producer(cfg: Config) -> Producer:
    return Producer({
        "bootstrap.servers": cfg.kafka_bootstrap,
        "security.protocol": cfg.kafka_security_protocol,
        "ssl.ca.location": cfg.kafka_ssl_ca,
        "ssl.certificate.location": cfg.kafka_ssl_cert,
        "ssl.key.location": cfg.kafka_ssl_key,
        "acks": "all",
    })


def emit_completed(producer: Producer, payload: dict) -> None:
    event = {
        "event_type": "embedding.completed",
        "source_id": payload.get("source_id"),
        "event_type_origin": payload.get("event_type"),
        "service_name": payload.get("service_name"),
        "created_at": payload.get("created_at"),
    }

    def _cb(err, _msg):
        if err:
            logger.error("Delivery failed: %s", err)

    producer.produce(
        topic=COMPLETED_TOPIC,
        key=str(payload.get("source_id", "")),
        value=json.dumps(event).encode("utf-8"),
        callback=_cb,
    )
    producer.poll(0)
