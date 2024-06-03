class IncomingMessage:
  def __init__(self, instruction, code, task_id):
    self.instruction = instruction
    self.code = code
    self.task_id = task_id

class OutgoingMessage:
  def __init__(self, generated, task_id):
    self.generated = generated
    self.task_id = task_id