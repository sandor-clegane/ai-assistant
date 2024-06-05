from sqlalchemy import Boolean, Column, String

from database import Base


class Task(Base):
    __tablename__ = "tasks"

    id = Column(String, primary_key=True)
    instruction = Column(String)
    code= Column(String)
    is_done = Column(Boolean, default=False)

