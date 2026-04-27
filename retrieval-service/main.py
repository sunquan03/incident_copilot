from fastapi import FastAPI
from pydantic import BaseModel
from store import VectorStore
from config import load_config
from embedder import Embedder
app = FastAPI()

cfg = load_config()

_embedder = Embedder(cfg.embed_model, cfg.embed_batch_size)
vectorstore = VectorStore(cfg, _embedder)

class Question(BaseModel):
    question: str
    tenant_id: str
    top_k: int | None = 5

@app.post("/retrieve")
async def retrieve(question: Question):
    res = vectorstore.search(tenant_id=question.tenant_id, question=question.question, top_k=question.top_k)
    return {"chunks": res}

@app.get("/health")
async def health():
    return {"status": "ok"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host=cfg.service_host, port=cfg.service_port)