# order_management

# Order Management Service

This service handles order creation, processing, and reporting through a REST API. It is built using Go, Docker, and Docker Compose, and uses SQLite for data storage. The service is designed to simulate order processing with configurable parameters for worker count, order processing timeout, and report interval.

## hint

The service supports concurrency with multiple workers, but SQLite can become a bottleneck in high-concurrency scenarios due to its single-writer limitation and locking mechanism. 



## Docker Setup

This project uses Docker and Docker Compose to build and run the application in a containerized environment.

### Docker Compose Configuration

The `docker-compose.yml` file defines the worker service. The service is configured with environment variables to control the worker count, order processing timeout, and report interval.

```yaml
version: '3.8'
services:
  worker:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8099:8080"
    environment:
      DB_LOG_LEVEL: INFO
      SERVICE_WORKER_COUNT: 10
      SERVICE_ORDER_PROCESS_TIMEOUT: 5
      SERVICE_REPORT_INTERVAL: 2
```

## example 


+ create order

```
curl -X POST http://10.10.10.10:8099/api/v1/orders \
-H "Content-Type: application/json" \
-d '{
  "order_id": "1235",
  "priority": "Normal",
  "processing_time": 1
}'
```

+ get order

```
curl -X GET "http://10.10.10.10:8099/api/v1/orders/3722" \
-H "Content-Type: application/json"  
```

+ python script

```py
import requests
import random
import json
import time

# Define the URL and headers
url = "http://10.21.10.13:8099/api/v1/orders"
headers = {"Content-Type": "application/json"}

# Function to generate a random order
def generate_random_order():
    order_id = str(random.randint(1000, 9999))  # Random order_id between 1000 and 9999
    priority = random.choice(["Normal", "High"])  # Random priority
    processing_time = random.randint(1, 10)  # Random processing time between 1 and 10 seconds
    return {
        "order_id": order_id,
        "priority": priority,
        "processing_time": processing_time
    }

# Send 100 requests, 1 millisecond apart
for i in range(100):
    # Generate a random order
    order_data = generate_random_order()
    
    # Send the POST request
    response = requests.post(url, headers=headers, data=json.dumps(order_data))
    
    # Print the response (Optional)
    if response.status_code == 200:
        print(f"[{i+1}/100] Order {order_data['order_id']} sent successfully!")
    else:
        print(f"[{i+1}/100] Failed to send order {order_data['order_id']}. Status code: {response.status_code}")
    
    # Wait for 1 millisecond before sending the next request
    time.sleep(0.001)

```

## Analysis and Planning

Before diving into the development process, a thorough analysis of the requirements was performed. The following images show the initial analysis and planning stages, outlining the architecture and workflows for the system.

![tips](https://github.com/seyedmo30/order_management/blob/main/docs/1.jpg)


![tips](https://github.com/seyedmo30/order_management/blob/main/docs/2.jpg)
