from fastapi import APIRouter
from schema import IncomingMessage, OutgoingMessage
from gpt2 import CodeGenerator
from config import loop, KAFKA_BOOTSTRAP_SERVERS, KAFKA_CONSUMER_GROUP, KAFKA_TOPIC_TASKS, KAFKA_TOPIC_RESPONSES
from aiokafka import AIOKafkaConsumer, AIOKafkaProducer
import json

route = APIRouter()
generator = CodeGenerator()

async def process_tasks():
    consumer = AIOKafkaConsumer(
        KAFKA_TOPIC_TASKS, 
        loop=loop,
        bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
        group_id=KAFKA_CONSUMER_GROUP,
    )
    await consumer.start()

    producer = AIOKafkaProducer(
        loop=loop, 
        bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
    )
    await producer.start()

    try:
        async for msg in consumer:
            message = IncomingMessage(**json.loads(msg.value))

            generated = generator.generate(message.code, message.instruction)
            
            message_out = {
                "task_id": message.task_id,
                "response": generated
            }
        
            response_json = json.dumps(message_out).encode()

            await producer.send_and_wait(topic=KAFKA_TOPIC_RESPONSES, value=response_json)
    finally:
        await consumer.stop()