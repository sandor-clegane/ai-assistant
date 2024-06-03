from transformers import GPT2Tokenizer, GPT2LMHeadModel
from config import FINAL_MODEL, FINAL_TOKENIZER, DEVICE
import torch
from os.path import dirname


class CodeGenerator:
    def __init__(self):
        absolute_path_final_model = f'{dirname(__file__)}' + FINAL_MODEL
        absolute_path_final_tokenizer = f'{dirname(__file__)}' + FINAL_TOKENIZER

        self.model = GPT2LMHeadModel.from_pretrained(absolute_path_final_model, torch_dtype=torch.float16).to(DEVICE)
        self.tokenizer = GPT2Tokenizer.from_pretrained(absolute_path_final_tokenizer)
        self.tokenizer.pad_token = self.tokenizer.eos_token

    def generate(self, code, instruction):
        text = f"""
        ### Question:
        {instruction}
        ### Context:
        {code}
        ### Answer:
        """

        input_ids = self.tokenizer.encode(text, return_tensors="pt").to(DEVICE)
    
        out = self.model.generate(input_ids, max_new_tokens=512,
                        do_sample=True, num_beams=6,
                        temperature=1, repetition_penalty=3.0)

        generated_text = self.tokenizer.batch_decode(out, skip_special_tokens=True)[0]
        return generated_text
    
