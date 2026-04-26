from qdrant_client import QdrantClient
from qdrant_client.models import Filter, FieldCondition, MatchValue
from config import Config
from embedder import Embedder


class VectorStore:
    def __init__(self, cfg: Config, embedder: Embedder) -> None:
        self.collection = cfg.qdrant_collection
        self.embedder = embedder
        self.client = QdrantClient(
            host=cfg.qdrant_host,
            port=cfg.qdrant_port,
            api_key=cfg.qdrant_api_key or None,
        )


    def search(self, tenant_id, question, top_k=5):
        vector =  self.embedder.embed([question])[0]
        hits = self.client.search(
            collection_name=self.collection,
            query_vector=vector,
            query_filter=Filter(
                must=[FieldCondition(key="tenant_id", match=MatchValue(value=tenant_id))],
            ),
            limit=top_k,
            with_payload=True,
        )

        return hits