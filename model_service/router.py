from fastapi import APIRouter
from schema import IncomingMessage
from gpt2 import CodeGenerator
from config import loop, KAFKA_BOOTSTRAP_SERVERS, KAFKA_CONSUMER_GROUP, KAFKA_TOPIC_TASKS, KAFKA_TOPIC_RESPONSES
from aiokafka import AIOKafkaConsumer, AIOKafkaProducer
import json
import asyncio
from multiprocessing import Lock

mutex = Lock()

route = APIRouter()
generator = CodeGenerator()

TaskQueue = []

async def process_tasks():
    while True:
        consumer = AIOKafkaConsumer(
            KAFKA_TOPIC_TASKS, 
            loop=loop,
            bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
            group_id=KAFKA_CONSUMER_GROUP,
        )
        await consumer.start()

        try:
            async for msg in consumer:
                message = IncomingMessage(**json.loads(msg.value))
                with mutex:
                    TaskQueue.append(message)
                print("message added to task queue")
        finally:
            await consumer.stop()

        await asyncio.sleep(1)

async def process_queue():
    producer = AIOKafkaProducer(
        loop=loop, 
        bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
    )
    await producer.start()

    while True:
        tasks  = []
        with mutex:
            tasks = TaskQueue[:10]
            del TaskQueue[:10]

        for i in range(len(tasks)):
            message = tasks[i]
            generated = generator.generate(message.code, message.instruction)
            message_out = {
                "task_id": message.task_id,
                "response": generated
            }
        
            response_json = json.dumps(message_out).encode()

            await producer.send_and_wait(topic=KAFKA_TOPIC_RESPONSES, value=response_json)
            print("task done")

        await asyncio.sleep(30)
