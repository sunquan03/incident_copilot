import logging
import signal

from config import load_config
from consumer import build_consumer, poll
from embedder import Embedder
from normalizers import NORMALIZERS
from producer import build_producer, emit_completed
from store import VectorStore

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s [%(name)s] %(message)s",
)
logger = logging.getLogger("embedding-service")

_running = True


def _stop(sig, _):
    global _running
    logger.info("Signal %d received, stopping...", sig)
    _running = False


def run():
    cfg = load_config()
    embedder = Embedder(cfg.embed_model, cfg.embed_batch_size)
    store = VectorStore(cfg)
    consumer = build_consumer(cfg)
    producer = build_producer(cfg)

    signal.signal(signal.SIGINT, _stop)
    signal.signal(signal.SIGTERM, _stop)

    logger.info("embedding-service running")

    try:
        while _running:
            topic, event, raw_msg = poll(consumer)
            if topic is None:
                continue

            normalizer = NORMALIZERS.get(topic)
            if normalizer is None:
                logger.warning("No normalizer for topic: %s", topic)
                consumer.commit(message=raw_msg)
                continue

            try:
                text, payload = normalizer(event)
            except Exception as e:
                logger.error("Normalize error on %s: %s", topic, e)
                consumer.commit(message=raw_msg)
                continue

            if not text:
                consumer.commit(message=raw_msg)
                continue

            try:
                vector = embedder.embed([text])[0]
                store.upsert(text, payload, vector)
            except Exception as e:
                # Don't commit — let Kafka redeliver
                logger.error("Embed/upsert failed on %s: %s", topic, e)
                continue

            emit_completed(producer, payload)
            consumer.commit(message=raw_msg)
            logger.info("Done: topic=%s source_id=%s", topic, payload.get("source_id"))

    finally:
        producer.flush(timeout=10)
        consumer.close()
        logger.info("embedding-service stopped")


if __name__ == "__main__":
    run()


