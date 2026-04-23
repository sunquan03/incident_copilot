import os
from dataclasses import dataclass, field


@dataclass
class Config:
    kafka_bootstrap: str
    kafka_security_protocol: str
    kafka_ssl_ca: str
    kafka_ssl_cert: str
    kafka_ssl_key: str
    kafka_group_id: str
    kafka_topics: list[str]
    qdrant_host: str
    qdrant_port: int
    qdrant_api_key: str
    qdrant_collection: str
    embed_model: str
    embed_batch_size: int


def load_config() -> Config:
    return Config(
        kafka_bootstrap=os.environ["KAFKA_BOOTSTRAP_SERVERS"],
        kafka_security_protocol=os.getenv("KAFKA_SECURITY_PROTOCOL", "SSL"),
        kafka_ssl_ca=os.getenv("KAFKA_SSL_CA", "certs/ca.pem"),
        kafka_ssl_cert=os.getenv("KAFKA_SSL_CERT", "certs/service.cert"),
        kafka_ssl_key=os.getenv("KAFKA_SSL_KEY", "certs/service.key"),
        kafka_group_id=os.getenv("KAFKA_GROUP_ID", "embedding-service"),
        kafka_topics=["alert.received", "incident.created", "logdoc.received"],
        qdrant_host=os.getenv("QDRANT_HOST", "localhost"),
        qdrant_port=int(os.getenv("QDRANT_PORT", "6333")),
        qdrant_api_key=os.getenv("QDRANT_API_KEY", ""),
        qdrant_collection=os.getenv("QDRANT_COLLECTION", "incident_knowledge"),
        embed_model=os.getenv("EMBED_MODEL", "all-MiniLM-L6-v2"),
        embed_batch_size=int(os.getenv("EMBED_BATCH_SIZE", "32")),
    )
