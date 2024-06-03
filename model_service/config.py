import os
from dotenv import load_dotenv
import asyncio

# Load environment variables from .env file
load_dotenv()

# Read environment variables
#   kafka
KAFKA_BOOTSTRAP_SERVERS = os.getenv("KAFKA_BOOTSTRAP_SERVERS")
KAFKA_TOPIC_TASKS = os.getenv("KAFKA_TOPIC_TASKS")
KAFKA_TOPIC_RESPONSES = os.getenv("KAFKA_TOPIC_RESPONSES")
KAFKA_CONSUMER_GROUP = os.getenv("KAFKA_CONSUMER_GROUP")
#   model
FINAL_MODEL = os.getenv("FINAL_MODEL")
FINAL_TOKENIZER = os.getenv("FINAL_TOKENIZER")
DEVICE = os.getenv("DEVICE")

# Create event loop
loop = asyncio.get_event_loop()
