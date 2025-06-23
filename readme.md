# 🛍️ ShopHub Microservice

**ShopHub** is a scalable and modular e-commerce system built with microservices architecture.  
This monorepo contains multiple services including:

- 🧾 `order-service` – handles order creation and tracking  
- 📦 `product-service` – manages product catalogs and inventory  
- 💳 `payment-service` – processes payments and integrates with payment gateways  
- 🚚 `shipment-service` – manages shipping, couriers, and tracking  
- 👤 `user-service` – handles user registration, authentication, and profiles  

---

## 🧰 Tech Stack

### 📦 Microservice Architecture
Each domain (`Order`, `Product`, `Payment`, `Shipment`, `User`) is implemented as an independent service with its own database.  
This follows the **Database-per-Service** pattern, enabling independent deployment, fault isolation, and true service ownership.

---

### 🗃️ Database – PostgreSQL
Each service uses its own **PostgreSQL** instance.  
This setup supports **Change Data Capture (CDC)** and allows domain isolation with full ownership over its schema.

---

### 🔄 Change Data Capture – Debezium

We use **two Debezium instances** for different event processing needs:

#### 🧱 Instance 1: Debezium + Kafka (High Throughput)
- **Purpose**: Event stream processing at scale.
- **Pattern**: Soft Copy CDC – captures raw `INSERT`, `UPDATE`, `DELETE` changes.
- **Use Cases**:
  - Product inventory syncing
  - User data projection to downstream services

#### ⚡ totaly 4 instance debezium server: Debezium server + sink nats (Real-Time Outbox Pattern)
- **Purpose**: Real-time, push-based communication between services.
- **Pattern**: Outbox CDC – application writes to an `outbox` table; Debezium publishes to nats.
- **Use Cases**:
  - Order placement flow (e.g., trigger payment generation)

> ✅ This hybrid CDC strategy leverages **Kafka's durability and scalability** and **Nats's real-time responsiveness**.

---

### 📡 Message Brokers

| Broker | Role                            | Characteristics         |
| ------ | ------------------------------- | ----------------------- |
| Kafka  | High-throughput event streaming | Pull-based, durable     |
| Nats   | Real-time message delivery      | Push-based, lightweight |

---

### 📊 Observability – Tracing, Logging, and Metrics

To support traceability and debugging across **event-driven** and **synchronous** workflows, we implement full observability using:

| Tool                     | Purpose                          |
| ------------------------ | -------------------------------- |
| **OpenTelemetry (OTel)** | Distributed Tracing (span-based) |
| **Grafana Tempo**        | Tracing backend (span storage)   |
| **Grafana Loki**         | Structured, indexed logging      |
| **Grafana Dashboards**   | Visualization of traces & logs   |

#### 🔁 Span-to-Log & Log-to-Span Integration

- **Span-to-Log**: Each OpenTelemetry span includes context-aware log entries.
- **Log-to-Span**: Each structured log includes `trace_id` and `span_id`, allowing Loki to correlate logs to Tempo traces.

> 🔍 This ensures end-to-end visibility even across asynchronous, event-based interactions where HTTP headers cannot propagate context natively.

#### ✅ Benefits:
- Traceability preserved even through **Kafka/Nats events**
- Enables root cause analysis without full request/response lifecycle
- Observability is **non-blocking** and **low-latency**

---

### 🛠️ Supporting Tools

| Tool           | Purpose                                 |
| -------------- | --------------------------------------- |
| Docker         | Containerization for microservices      |
| Docker Compose | Local orchestration                     |
| gRPC / REST    | Inter-service communication             |
| JWT / OAuth2   | Authentication & authorization          |
| Redis          | Caching, rate limiting, session storage |

---

### 🎯 Learning Objective

This project aims to demonstrate **how to achieve full observability in an event-driven microservice architecture** without losing visibility across:
- Service boundaries
- Message brokers (Kafka/Nats)
- CDC pipelines
- Synchronous and asynchronous flows

By using OpenTelemetry with Tempo and Loki, developers can **trace every action** from user request to database write and downstream message consumption — **in real time** and without additional complexity or blocking I/O.

## 👤 User Flow (Service Communication via CDC & Message Brokers)

In ShopHub, services are decoupled and communicate through **Change Data Capture (CDC)** and **message brokers** (Kafka & Nats). Below are the key user flows and how data propagates between services:

---

### 🧍 1. User Account Lifecycle → Consumed by `order-service` and `payment-service`

| Action         | Source         | Consumed By                       | Method           |
| -------------- | -------------- | --------------------------------- | ---------------- |
| Create Account | `user-service` | `order-service` `payment-service` | Debezium + Kafka |
| Update Account | `user-service` | `order-service` `payment-service` | Debezium + Kafka |
| Delete Account | `user-service` | `order-service` `payment-service` | Debezium + Kafka |

🔄 Order service listens to user changes to maintain denormalized customer data and order ownership.

---

### 🏠 2. User Address CRUD → Consumed by `shipment-service`

| Action         | Source         | Consumed By        | Method           |
| -------------- | -------------- | ------------------ | ---------------- |
| Create Address | `user-service` | `shipment-service` | Debezium + kafka |
| Update Address | `user-service` | `shipment-service` | Debezium + kafka |
| Delete Address | `user-service` | `shipment-service` | Debezium + kafka |

📦 Shipment service uses this data for delivery information and logistics routing.

---

### 🛒 3. Product CRUD → Consumed by `order-service`

| Action         | Source            | Consumed By     | Method           |
| -------------- | ----------------- | --------------- | ---------------- |
| Create Product | `product-service` | `order-service` | Debezium + Kafka |
| Update Product | `product-service` | `order-service` | Debezium + Kafka |
| Delete Product | `product-service` | `order-service` | Debezium + Kafka |

🧾 Order service listens to product changes to maintain accurate product info in orders (e.g., price, name, availability).

---

> ⚙️ These flows ensure that each service remains autonomous yet synchronized via event streams. No direct API calls between services are required, promoting loose coupling and high resilience.

## 🔄 Order Flow

The order lifecycle in **ShopHub** is designed using a mix of **synchronous** and **event-driven** communication patterns to balance **data consistency**, **scalability**, and **user experience**.

---

### 🛍️ 1. User Places an Order

- The `order-service` receives the order submission.
- It performs local validations using its own data:
  - Validates product availability and stock (replicated product data).
  - Verifies that the user exists and is verified.
- Then it makes a **synchronous API call** to `shipment-service` to **calculate shipping cost**.
  - 🧷 **Note**: This is a tightly coupled integration due to the need for immediate pricing data.

---

### 📦 2. Order Initialization (Outbox Pattern)

- After validation and shipping cost calculation, `order-service` creates an **Outbox Event** with `aggregateType: OrderInitialization`.
- This event is picked up by **Debezium + Nats** and consumed by the `payment-service`.

---

### 💳 3. Payment Service: Create Payment Link

- `payment-service` creates a **payment link** (e.g., Midtrans, Xendit).
- Frontend client polls via **SSE** or **long polling** to monitor payment status.
- Once user completes payment, the **payment gateway sends a callback** to `payment-service`.

---

### ✅ 4. On Payment Success

- `payment-service` publishes a **PaymentSuccess** event.
- `order-service` consumes this and:
  - Updates the order status to **"Being Packed"**.
  - Calls `product-service` to re-check and **deduct stock**.
- If stock is **not available**:
  - `product-service` publishes a **StockFailed** message.
  - `order-service` updates status to **Failed**, triggers rollback.
  - `payment-service` auto-triggers a **refund** via gateway.

---

### 🚚 5. Shipment Scheduling & Compensation

- On successful stock verification, `shipment-service` schedules the delivery (asynchronously).
- If shipment **fails** due to external issues (e.g., courier unavailability, invalid address):
  - `shipment-service` publishes a **ShipmentFailed** event via outbox pattern.
  - This triggers a **compensation workflow**, including:
    - `payment-service` → Automatically issues a **refund** to the user.
    - `order-service` → Updates order status to **"Failed"**.
    - `product-service` → **Restores the previously deducted stock** to maintain data consistency.

---

### 🔁 Product Compensation Logic

When `product-service` receives the **ShipmentFailed** event:

- It identifies the associated order and the deducted stock.
- Restores the stock quantities to their previous state (before the deduction).
- Publishes a **StockCompensationSuccess** or **StockRestoreSuccess** event for auditing (optional).

> ✅ This ensures that inventory remains accurate even in multi-step failure scenarios, aligning with **Saga Pattern's compensation strategy**.

---

### 🧠 Summary of Key Mechanisms

| Step                      | Mechanism                                  |
| ------------------------- | ------------------------------------------ |
| Product & User Validation | Local DB in `order-service`                |
| Shipping Fee              | Sync call to `shipment-service`            |
| Payment Flow              | Outbox → Nats + Callback                   |
| Stock Final Check         | Sync in `product-service` + rollback logic |
| Shipment Scheduling       | Event-driven via outbox                    |
| Compensation (Failure)    | Event-based rollback + refund              |

> ✅ The design blends sync and async communication to balance user responsiveness and system decoupling.  
> 🧩 Compensation logic ensures system resilience in failure scenarios.
