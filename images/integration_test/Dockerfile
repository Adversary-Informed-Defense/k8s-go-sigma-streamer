FROM python:3.10.0-slim

RUN apt-get update && \
    apt-get install -y gcc librdkafka-dev

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY . .
