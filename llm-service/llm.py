from openai import OpenAI
from prompt import prompt_messages


class LLMClient:
    def __init__(self, base_url: str, api_key: str, model: str):
        self.client = OpenAI(api_key=api_key, base_url=base_url)
        self.model = model

    def chat(self, question: str, context: list) -> str:
        prompt_msg = prompt_messages(question, context)
        response = self.client.chat.completions.create(
            model=self.model,
            messages=prompt_msg,
            temperature=0.0,
        )

        return response.choices[0].message.content