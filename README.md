# 👤 FlyHorizons - User Service

This is the **User Service** for **FlyHorizons**, an enterprise-grade airline booking system. The service manages users and administrators, including authentication, profile data, and role-based access control across the platform.

---

## 🚀 Overview

This microservice handles all **user identity and access management** within the FlyHorizons ecosystem. It supports passenger registration, admin onboarding, authentication, and role management. The User Service also integrates with other services—such as Booking, Email, and Payment—through **RabbitMQ** to propagate user-related events.

Built with **Go** and the **Gin** framework, the service ensures secure and scalable user data handling for both frontend portals and backend services.

---

## 🛠️ Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin
- **Database**: Microsoft SQL Server (Azure Hosted)
- **Messaging**: RabbitMQ
- **Authentication**: JWT
- **Architecture**: Microservices

---

## 📦 Features

- 👥 **User registration and login**
- 🛂 **Admin role management**
- 🔐 **JWT-based authentication**
- 🔄 **Publishes events** (e.g., user deleted, user registered)
- 📬 **Integrates with Email and Booking Services**
- 🧩 **Role-based access control (RBAC)**
- 🧾 **User profile updates and preferences**
- ⚠️ **Centralized error handling and input validation**

---

## 📄 License
This project is shared for educational and portfolio purposes only. Commercial use, redistribution, or modification is not allowed without explicit written permission. All rights reserved © 2025 Beatrice Marro.

## 👤 Author
Beatrice Marro GitHub: https://github.com/beamarro
