import json
import logging
from confluent_kafka import Consumer, KafkaError, KafkaException, Message
from config import Config

logger = logging.getLogger(__name__)

def build_consumer(cfg: Config) -> Consumer:
    consumer = Consumer({
        "bootstrap.servers": cfg.kafka_bootstrap,
        "security.protocol": cfg.kafka_security_protocol,
        "ssl.ca.location": cfg.kafka_ssl_ca,
        "ssl.certificate.location": cfg.kafka_ssl_cert,
        "ssl.key.location": cfg.kafka_ssl_key,
        "group.id": cfg.kafka_group_id,
        "auto.offset.reset": "earliest",
        "enable.auto.commit": False,
    })
    consumer.subscribe(cfg.kafka_topics)
    logger.info("Subscribed to: %s", cfg.kafka_topics)
    return consumer


def poll(consumer: Consumer, timeout: float = 1.0):
    """Returns (topic, payload, raw_msg) or (None, None, None)."""
    msg: Message = consumer.poll(timeout)
    if msg is None:
        return None, None, None
    if msg.error():
        code = msg.error().code()
        if code == KafkaError._PARTITION_EOF:
            return None, None, None
        if code == KafkaError.UNKNOWN_TOPIC_OR_PART:
            logger.warning("Topic not yet available, waiting: %s", msg.topic())
            return None, None, None
        raise KafkaException(msg.error())
    try:
        payload = json.loads(msg.value().decode("utf-8"))
    except (json.JSONDecodeError, UnicodeDecodeError) as e:
        logger.warning("Error message on %s, skipping: %s", msg.topic(), e)
        consumer.commit(message=msg)
        return None, None, None
    return msg.topic(), payload, msg

