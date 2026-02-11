import ollama
import uvicorn
from fastapi import FastAPI

class Neuro:
    def __init__(self, model: str):
        self.model = model
        super().__init__()

    async def ask(self, prompt: str):
        response = await ollama.generate(
            model=self.model,
            prompt=prompt,
            options={
                "temperature": 0.7
            }
        )
        return response['response']


class Server:
    def __init__(self):
        self.app = FastAPI()
        self.neuro = Neuro("qwen3-vl:8b")
        self.setup()

    def setup(self):
        @self.app.get("/ask")
        async def ask(data):
            answer = await self.neuro.ask(data.prompt)
            return {"response": answer}

    def run(self):
        uvicorn.run(
            self.app,
            host="0.0.0.0",
            port=10000,
        )

if __name__ == "__main__":
    server = Server()
    server.run()