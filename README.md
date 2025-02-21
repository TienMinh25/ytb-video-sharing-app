# Youtube Video Sharing App

## Introduction

**Youtube Video Sharing App** is a web application that allows users to register, log in, and share their favorite YouTube videos by simply pasting a link. Users can also browse videos shared by others and receive real-time notifications when new videos are added.

### Purpose
The purpose of this application is to create an interactive platform where users can easily share and discover YouTube videos while engaging with a real-time notification system.

## Key Features
- **User Registration & Authentication**: Users can sign up and log in using a token-based authentication system (no OAuth required).
- **Video Sharing**: Users can share YouTube video links, which will be stored and displayed in the app.
- **Video Listing**: Users can browse a list of videos shared by others.
- **Real-Time Notifications**: When a new video is shared, all logged-in users receive real-time notifications via WebSockets.

## System Design

### Assumptions
- **Total registered users**: 10,000
- **Daily Active Users (DAU)**: 5,000
- **Concurrent users**: 1,000
- **Videos shared per day**: 5,000
- **Each user shares at least 1 video per day**

### Estimations
- **Database Write Load**: 5,000 QPS (queries per second)
- **Database Read Load**: 100,000 QPS (assuming 20x read requests per write)
- **Average video size**: 300MB - 500MB (1080p, 10-15 minutes)
- **Estimated storage needed**:
  - **Per day**: 500MB * 1,000 videos ≈ 488GB
  - **Per year**: 488GB * 365 ≈ 174TB
- **Metadata Storage**:
  - **Video metadata**: ~49MB/day → ~17.5GB/year
  - **User metadata**: ~9.7MB for 10K users

### High-Level Architecture
```
                            |-> save expire key in redis (to invalidate refresh token)
                            |
                            | Log out request           
                            |
FE -> Load Balancer ----> BE 1  -------   <share new video, push event to Kafka>
 |      | |         |---> BE 2  ------| -----------------------------------------> Kafka --------------------------- 
 -------| |         |---> BE 3  ------|                                                                            |
 |        |                           |                                                                            |
 |        |        -------------------|                                                                            |
 |        |        |                   MySQL (Account, Video metadata, Account password) (~=17.5GB + 9.7MB)        |
 |        |        |                                                                                               |
 |--------|----->Blob Storage (video) (1 year) ~ 174TB                                                             |
          |                                                                                                        |
          ---------------------------------------------------------------------------------------------------------|
          <Background job consumes messages from Kafka, checking online users via WebSockets>
```

## API Design

### Authentication
#### Login
- **Endpoint**: `POST /api/v1/login`
- **Payload**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Responses**:
  - `200 OK`
    ```json
    {
      "accessToken": "string",
      "refreshToken": "string",
      "metadata": { "email": "string", "avatar": "string" }
    }
    ```
  - `400 Bad Request` (invalid credentials)

#### Logout
- **Endpoint**: `POST /api/v1/logout`
- **Request Headers**:
  - `Authorization`: Bearer Token
- **Response**:
  - `200 OK`

#### Register
- **Endpoint**: `POST /api/v1/register`
- **Payload**:
  ```json
  {
    "email": "string",
    "password": "string",
    "confirmPassword": "string",
    "fullName": "string",
    "avatarURL": "string (optional)"
  }
  ```
- **Responses**:
  - `200 OK`
  - `409 Conflict` (email already exists)

### Video Management
#### List Videos
- **Endpoint**: `GET /api/v1/videos`
- **Headers**: `Authorization: Bearer Token`
- **Response**:
  ```json
  [
    "data": { 
      "id": "int", 
      "title": "string", 
      "description": "string", 
      "thumbnail": "string" 
    },
    "metadata": { 
      "page": 2, 
      "limit": 10,
      "total": 103,
      "totalPage": 11,
      "isNext": true,
      "isPrevious": true
    }
  ]
  ```

#### Share a Video
- **Endpoint**: `POST /api/v1/videos`
- **Headers**: `Authorization: Bearer Token`
- **Payload**:
  ```json
  {
    "description": "string",
    "thumbnail": "string",
    "video_url": "string",
    "account_id": "int"
  }
  ```
- **Response**:
  - `200 OK`

## Database Schema

### Account Table
| Column     | Type   |
|------------|--------|
| id         | int    |
| email      | string |
| fullname   | string |
| avatarURL  | string |

### Account_Password Table
| Column     | Type    |
|------------|--------|
| id         | int    |
| password   | string |

### Video Table
| Column      | Type   |
|------------|--------|
| id         | int    |
| description| string |
| upvote     | int64  |
| downvote   | int64  |
| thumbnail  | string |
| video_url  | string |
| account_id | int    |

### Kafka Message Schema
```json
{
  "video_title": "string",
  "username": "string" // or using fullname instead
}
```

## Installation & Configuration (update later)

1. **Clone the repository:**
   ```sh
   git clone https://github.com/your-repo/youtube-video-sharing.git
   cd youtube-video-sharing
   ```

2. **Set up environment variables:**
   ```sh
   cp .env.example .env
   ```

3. **Install dependencies:**
   ```sh
   go mod tidy
   ```

## Running the Application (update later)

1. **Start Kafka (if not running):**
   ```sh
   docker-compose up -d kafka zookeeper
   ```

2. **Run the backend server:**
   ```sh
   go run main.go
   ```

## Docker Deployment (update later)

1. **Build and run the Docker container:**
   ```sh
   docker-compose up --build
   ```

## Troubleshooting (update later)

### Common Issues
1. **Kafka not running:** Ensure Kafka is up and running before starting the application.
   ```sh
   docker ps | grep kafka
   ```
2. **Database connection issues:** Verify that your MySQL database URL in `.env` is correct and MySQL is running.

## Conclusion
This application demonstrates a scalable approach to building a YouTube video-sharing platform with real-time notifications using Go, Kafka, WebSockets, and MySQL.

---
### **🚀 Happy Coding!**

