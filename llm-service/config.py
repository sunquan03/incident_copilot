import os
from dataclasses import dataclass
from dotenv import load_dotenv


@dataclass
class Config:
    service_host: str
    service_port: int
    llm_base_url: str
    llm_api_key: str
    llm_model: str


def load_config() -> Config:
    load_dotenv()
    return Config(
        service_host=os.getenv('SERVICE_HOST'),
        service_port=int(os.getenv('SERVICE_PORT')),
        llm_base_url=os.getenv('LLM_BASE_URL'),
        llm_api_key=os.getenv('LLM_API_KEY'),
        llm_model=os.getenv('LLM_MODEL'),
    )
