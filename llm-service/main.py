from fastapi import FastAPI
from pydantic import BaseModel
from llm import LLMClient
from config import load_config
import uvicorn

app = FastAPI()


cfg = load_config()

llmClient = LLMClient(cfg.llm_base_url, cfg.llm_api_key, cfg.llm_model)

class ContextRequest(BaseModel):
    question: str
    context: list
@app.post("/llm")
async def ask_llm(ctx_req: ContextRequest):
    res = llmClient.chat(ctx_req.question, ctx_req.context)
    return {"response": res}

@app.get("/health")
async def health():
    return {"status": "ok"}


if __name__ == "__main__":
    uvicorn.run(app, host=cfg.service_host, port=cfg.service_port)