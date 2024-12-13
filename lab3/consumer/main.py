import pika
import json
import requests

def callback(ch, method, properties, body):
    arr = json.loads(body.decode())
    for product in arr['products']:
        response = requests.post("http://localhost:8080/monitor/", json=product)
        
        if response.status_code == 200:
            print("Successfully sent to the URL")
        else:
            print(f"Failed to send, HTTP {response.status_code}, Response: {response.text}")


def consumer():
    # Connect to RabbitMQ server
    connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
    channel = connection.channel()

    # Declare a queue
    channel.queue_declare(queue='monitors')
    channel.basic_consume(queue='monitors', on_message_callback=callback, auto_ack=True)

    print('Consuming')
    channel.start_consuming()

if __name__ == '__main__':
    consumer()