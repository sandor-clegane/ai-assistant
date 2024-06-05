FROM python:3.12.3

WORKDIR /code

COPY ./requirements.txt /code/requirements.txt

RUN python -m pip install --upgrade pip
RUN pip install -r /code/requirements.txt

COPY . /code/

CMD ["uvicorn", "main:app", "--reload"]
