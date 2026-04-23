import hashlib
import logging
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams, PointStruct
from config import Config

logger = logging.getLogger(__name__)


def _make_point_id(source_id: str) -> str:
    digest = hashlib.md5(source_id.encode()).hexdigest()
    return f"{digest[:8]}-{digest[8:12]}-{digest[12:16]}-{digest[16:20]}-{digest[20:]}"


class VectorStore:
    def __init__(self, cfg: Config) -> None:
        self.collection = cfg.qdrant_collection
        self.client = QdrantClient(
            host=cfg.qdrant_host,
            port=cfg.qdrant_port,
            api_key=cfg.qdrant_api_key or None,
        )
        self._ensure_collection()

    def _ensure_collection(self) -> None:
        names = {c.name for c in self.client.get_collections().collections}
        if self.collection not in names:
            logger.info("Creating collection: %s", self.collection)
            self.client.create_collection(
                collection_name=self.collection,
                vectors_config=VectorParams(size=384, distance=Distance.COSINE),
            )

    def upsert(self, text: str, payload: dict, vector: list[float]) -> None:
        point_id = _make_point_id(str(payload.get("source_id", text)))
        self.client.upsert(
            collection_name=self.collection,
            points=[
                PointStruct(
                    id=point_id,
                    vector=vector,
                    payload={"text": text, **payload},
                )
            ],
        )
