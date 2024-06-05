from fastapi import FastAPI
import router
import asyncio

app = FastAPI()

@app.get("/")
async def home():
    return "Welcome Home"

app.include_router(router.route)
asyncio.create_task(router.process_tasks())
asyncio.create_task(router.process_queue())