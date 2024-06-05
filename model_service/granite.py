from config import GRANITE
import torch
from peft import (
    LoraConfig,
    get_peft_model,
    get_peft_model_state_dict,
    prepare_model_for_kbit_training,
    set_peft_model_state_dict,
)
from transformers import AutoTokenizer, AutoModelForCausalLM, DataCollatorForSeq2Seq

model_path = f'{dirname(__file__)}' + GRANITE

model = AutoModelForCausalLM.from_pretrained(
    model_path.format('base'),
    load_in_8bit=True,
    torch_dtype=torch.float16,
    device_map="auto",
)
tokenizer = AutoTokenizer.from_pretrained(model_path.format('base'))

config = LoraConfig.from_pretrained(model_path)
config.inference_mode = True

model = get_peft_model(model, config)

tokenizer.add_eos_token = True
tokenizer.pad_token_id = 0
tokenizer.padding_side = "left"
