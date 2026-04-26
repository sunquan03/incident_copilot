import os
from dataclasses import dataclass, field


@dataclass
class Config:
    qdrant_host: str
    qdrant_port: int
    qdrant_api_key: str
    qdrant_collection: str
    embed_model: str
    embed_batch_size: int


def load_config() -> Config:
    return Config(
        qdrant_host=os.getenv("QDRANT_HOST", "localhost"),
        qdrant_port=int(os.getenv("QDRANT_PORT", "6333")),
        qdrant_api_key=os.getenv("QDRANT_API_KEY", ""),
        qdrant_collection=os.getenv("QDRANT_COLLECTION", "incident_knowledge"),
        embed_model=os.getenv("EMBED_MODEL", "all-MiniLM-L6-v2"),
        embed_batch_size=int(os.getenv("EMBED_BATCH_SIZE", "32")),
    )
