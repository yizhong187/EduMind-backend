# EduMind-backend
This is the repo for the EduMind mobile application's backend.

## Stucture of Models

### Student

Represents a student in the system.

#### Attributes

- **student_id** (`uuid.UUID`): The unique identifier for the student.
- **username** (`string`): The username chosen by the student.
- **email** (`string`): The email address of the student.
- **created_at** (`time.Time`): The timestamp when the student account was created.
- **name** (`string`): The full name of the student.
- **valid** (`bool`): Indicates if the student account is currently valid.

#### Example JSON Representation

```json
{
  "student_id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "student123",
  "email": "student123@example.com",
  "created_at": "2024-07-15T12:00:00Z",
  "name": "John Doe",
  "valid": true
}
```

### Subject (Tutor's specialisation)

Represents a subject in which a tutor specializes.

#### Attributes

- **subject** (`string`): The name of the subject.
- **yoe** (`int32`): The years of experience the tutor has in teaching this subject.

#### Example JSON Representation

```json
{
  "subject": "Mathematics",
  "yoe": 5
}
```

### Tutor

Represents a tutor in the system.

#### Attributes

- **tutor_id** (`uuid.UUID`): The unique identifier for the tutor.
- **username** (`string`): The username chosen by the tutor.
- **email** (`string`): The email address of the tutor.
- **created_at** (`time.Time`): The timestamp when the tutor account was created.
- **name** (`string`): The full name of the tutor.
- **valid** (`bool`): Indicates if the tutor account is currently valid.
- **subjects** (array of `Subject`): An array of subjects that the tutor specializes in.
- **verified** (`bool`): Indicates whether the tutor’s account is verified (`true`) or not (`false`).
- **rating** (`float64`): The average rating given to the tutor.
- **rating_count** (`int32`): The total number of ratings received by the tutor.

#### Example JSON Representation

```json
{
  "tutor_id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "tutor456",
  "email": "tutor456@example.com",
  "created_at": "2024-07-15T12:00:00Z",
  "name": "Jane Smith",
  "valid": true,
  "subjects": [
    {
      "subject": "Mathematics",
      "yoe": 5
    },
    {
      "subject": "Physics",
      "yoe": 3
    }
  ],
  "verified": true,
  "rating": 4.7,
  "rating_count": 35
}
```

### Chat

Represents a chat session between a student and a tutor.

#### Attributes

- **chat_id** (`int32`): The unique identifier for the chat session.
- **student_id** (`uuid.UUID`): The unique identifier of the student participating in the chat.
- **tutor_id** (`uuid.NullUUID`): The unique identifier of the tutor participating in the chat (nullable).
- **created_at** (`time.Time`): The timestamp when the chat session was created.
- **subject** (`int32`): The identifier of the subject associated with the chat.
- **topic** (`sql.NullString`): The topic of the chat session (nullable).
- **header** (`string`): A brief header or title for the chat session.
- **photo_url** (`sql.NullString`): The URL of an optional photo associated with the chat (nullable).
- **completed** (`bool`): Indicates whether the chat session is completed (`true`) or ongoing (`false`).

#### Example JSON Representation

```json
{
  "chat_id": 12345,
  "student_id": "550e8400-e29b-41d4-a716-446655440000",
  "tutor_id": null,
  "created_at": "2024-07-15T12:00:00Z",
  "subject": 1,
  "topic": null,
  "header": "Mathematics Tutoring Session",
  "photo_url": null,
  "completed": false
}
```

### Message

Represents a message within a chat session.

#### Attributes

- **message_id** (`uuid.UUID`): The unique identifier for the message.
- **chat_id** (`int32`): The identifier for the chat session to which the message belongs.
- **user_id** (`uuid.UUID`): The unique identifier of the user who sent the message.
- **created_at** (`time.Time`): The timestamp when the message was created.
- **updated_at** (`time.Time`): The timestamp when the message was last updated.
- **deleted** (`bool`): Indicates if the message is deleted (`true`) or not (`false`).
- **content** (`string`): The content of the message.

#### Example JSON Representation

```json
{
  "message_id": "550e8400-e29b-41d4-a716-446655440000",
  "chat_id": 12345,
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "created_at": "2024-07-15T12:00:00Z",
  "updated_at": "2024-07-15T12:05:00Z",
  "deleted": false,
  "content": "Hello, how are you?"
}
```

## API Endpoints

### General Routes
Base URL: `/v1`

1. Health Check
   - Route: `/v1/healthz`
   - Method: GET
   - Purpose: Check the readiness of the service.
   - Request Parameters: None
   - Responses:
     - 200 OK: Service is ready.

2. Error Testing
   - Route: `/v1/error`
   - Method: GET
   - Purpose: Simulate an error response for testing.
   - Request Parameters: None
   - Responses:
     - 500 Internal Server Error: Error simulated successfully.

### Student Routes
Base URL: `/v1/students`

1. Health Check
   - Route: `/v1/students/healthz`
   - Method: GET
   - Purpose: Check the readiness of the student service.
   - Request Parameters: None
   - Responses:
     - 200 OK: Service is ready.

2. Error Testing
   - Route: `/v1/students/err`
   - Method: GET
   - Purpose: Simulate an error response for testing.
   - Request Parameters: None
   - Responses:
     - 500 Internal Server Error: Error simulated successfully.

3. Student Registration
   - Route: `/v1/students/register`
   - Method: POST
   - Purpose: Register a new student.
   - Request Parameters:
     - Body:
       - 'username': username
       - 'password': password
       - 'name': name
       - 'email': email
   - Responses:
     - 201 Created: Registration successful.
     - 400 Bad Request: Invalid registration details.

4. Student Registration
   - Route: `/v1/students/register`
   - Method: POST
   - Purpose: Register a new student.
   - Request Parameters:
     - Body:
       - 'username': username
       - 'password': password
       - 'name': name
       - 'email': email
   - Responses:
     - 201 Created: Registration successful.
     - 400 Bad Request: Invalid registration details.

5. Get Student Profile
   - Route: `/v1/students/profile`
   - Method: GET
   - Purpose: Retrieve the profile of the authenticated student.
   - Request Parameters: None
   - Responses:
     - 200 OK: Profile retrieved successfully.
     - 401 Unauthorized: Authentication required.

6. Update Student Profile
   - Route: `/v1/students/profile`
   - Method: PUT
   - Purpose: Update the profile of the authenticated student.
   - Request Parameters:
     - Body:
       - 'username': username
       - 'name': name
   - Responses:
     - 200 OK: Profile updated successfully.
     - 400 Bad Request: Invalid profile details.
     - 401 Unauthorized: Authentication required.

### Tutor Routes
Base URL: `/v1/tutors`

1. Health Check
   - Route: `/v1/tutors/healthz`
   - Method: GET
   - Purpose: Check the readiness of the tutor service.
   - Request Parameters: None
   - Responses:
     - 200 OK: Service is ready.

2. Error Testing
   - Route: `/v1/tutors/err`
   - Method: GET
   - Purpose: Simulate an error response for testing.
   - Request Parameters: None
   - Responses:
     - 500 Internal Server Error: Error simulated successfully.

3. Get Tutor Profile
   - Route: `/v1/tutors/profile`
   - Method: GET
   - Purpose: Retrieve the profile of the authenticated tutor.
   - Request Parameters: None
   - Responses:
     - 200 OK: Profile retrieved successfully.
     - 401 Unauthorized: Authentication required.

4. Update Tutor Profile
   - Route: `/v1/tutors/profile`
   - Method: PUT
   - Purpose: Update the profile of the authenticated tutor.
   - Request Parameters:
     - Body:
       - 'username': username
       - 'name': name
   - Responses:
     - 200 OK: Profile updated successfully.
     - 400 Bad Request: Invalid profile details.
     - 401 Unauthorized: Authentication required.

5. Get Student Profile
   - Route: `/v1/tutors/studentProfile`
   - Method: GET
   - Purpose: Retrieve the profile of a student for the authenticated tutor.
   - Request Parameters: None
   - Responses:
     - 200 OK: Profile retrieved successfully.
     - 401 Unauthorized: Authentication required.

### Chat Routes
Base URL: `/v1/chat`
Middleware: MiddlewareUserAuth for all routes

1. Start New Chat
   - Route: `/v1/chat/new`
   - Method: POST
   - Middleware: MiddlewareStudentAuth
   - Purpose: Start a new chat session.
   - Request Parameters:
     - Body:
       - 'subject': subject
       - 'header': header
   - Responses:
     - 201 Created: New chat started successfully.
     - 400 Bad Request: Invalid chat details.
     - 401 Unauthorized: Authentication required.

2. Get All Chats
   - Route: `/v1/chat/`
   - Method: GET
   - Purpose: Retrieve all chat sessions.
   - Request Parameters: None
   - Responses:
     - 200 OK: Chats retrieved successfully.
     - 401 Unauthorized: Authentication required.

3. Get All Messages in a Chat
   - Route: `/v1/chat/{chatID}/view`
   - Method: GET
   - Middleware: MiddlewareChatAuth
   - Purpose: Retrieve all messages in a specific chat.
   - Request Parameters:
     - Path Parameter: chatID (ID of the chat session)
   - Responses:
     - 200 OK: Messages retrieved successfully.
     - 401 Unauthorized: Authentication required.
     - 404 Not Found: Chat not found.

5. Post New Message
   - Route: `/v1/chat/{chatID}/new`
   - Method: POST
   - Middleware: MiddlewareChatAuth
   - Purpose: Post a new message in a specific chat.
   - Request Parameters:
     - Path Parameter: chatID (ID of the chat session)
     - Body:
       - 'content': content
   - Responses:
     - 200 OK: Messages retrieved successfully.
     - 401 Unauthorized: Authentication required.
     - 404 Not Found: Chat not found.

4. Join Chat Room
   - Route: `/v1/chat/{chatID}/join`
   - Method: GET
   - Middleware: MiddlewareChatAuth
   - Purpose: Join a chat room for WebSocket communication.
   - Request Parameters:
     - Path Parameter: chatID (ID of the chat session)
   - Responses:
     - 200 OK: Joined chat room successfully.
     - 401 Unauthorized: Authentication required.
     - 404 Not Found: Chat room not found.

## Deployment to Heroku
Use `env GOOS=linux GOARCH=amd64 GOARM=7 go build` to compile into linux based binary code for heroku