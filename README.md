# Author: Herasymenko Daniil

---

# Real-Time Leaderboard Backend System

This project focuses on building a backend system for a **real-time leaderboard service**. The service will allow users track their scores, and view their rankings on a leaderboard. The system will include the following key features:

- **User Authentication**
- **Score Submission**
- **Real-Time Leaderboard Updates**
- **Paginate Through Score History**


---

## Project Requirements

The system should meet the following requirements:

### 1. **User Authentication**
- Allow users to **register** and **log in** to the system securely.
- Use JWT Tokens for authentication and authorization.

### 2. **Score Submission**
- Enable users to submit their scores.
- Maintain accurate tracking of submitted scores.

### 3. **Leaderboard Updates**
- Provide a **leaderboard** that displays the top users.
- Update the leaderboard in real time to reflect new scores.

### 4. **User Rankings**
- Allow users to view their **individual rankings** on the leaderboard.
- Ensure the rankings are updated dynamically based on the scores.

### 5. **Top Players Report**
- Generate reports on the **top players**

---

### Technology Stack
- **Web framework**: Fiber
- **Database**: MongoDB
- **Authentication**: JWT
- **Deployment**: Docker for containerization
- **API Documentation**: **Swagger** for API documentation and testing
- **WebSockets**: for real-time communication
- **Logging**: Comprehensive logging for system operations and debugging - **slog**


This project aims to deliver a scalable and efficient backend system for real-time leaderboard management, suitable for a variety of competitive platforms.

---

### Documentation

Swagger documentation is available at the `http://127.0.0.1:8082/swagger/index.html` endpoint.


