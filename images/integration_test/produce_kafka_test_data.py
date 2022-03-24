# TODO: Add consume and comparison with expected output from alerts topic to this test

"""Produce test events to kafka"""

from os import environ as env

from pathlib import Path
import confluent_kafka


def main():
    """Prepare data and produce to kafka"""

    # Create producer, supply producer configurations using librdkafka reference
    kafka_endpoint = f"{env.get('KAFKA_SERVICE_HOST')}:{env.get('KAFKA_SERVICE_PORT')}"
    producer_conf = {"bootstrap.servers": kafka_endpoint}
    producer = confluent_kafka.Producer(**producer_conf)
    total_rows = 0

    # Iterate over data files
    files = Path.cwd().glob("data/*/*/*.*")
    for file_path in files:

        # Get Kafka topic name from data directory structure
        topic = "-".join(str(file_path).split("data")[1].split("/")[1:3])

        with open(file_path) as rows:
            for row in rows:
                # Produce each line of each file as text
                producer.produce(topic, value=row)
                total_rows += 1

    # Block until producer finishes publishing
    producer.flush()
    print(f"Published {total_rows} rows")


if __name__ == "__main__":
    main()
