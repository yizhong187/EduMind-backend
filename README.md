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
  "topic": {
            "String": "",
            "Valid": false
        },
  "header": "Mathematics Tutoring Session",
  "photo_url": {
            "String": "https://res.cloudinary.com/dnc1q8tlu/image/upload/v1720609760/file_kyd4yc.jpg",
            "Valid": true
        },
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

### General

**Base URL**: `/v1`

<details>
 <summary><code>GET</code> <code><b>/healthz</b></code> Check the readiness of the service.</summary>

##### Parameters

> None


##### Responses

> | HTTP Code     | Response                                                            |
> |---------------|---------------------------------------------------------------------|
> | `200`         | `Service ready`                                |

</details>

<details>
 <summary><code>GET</code> <code><b>/error</b></code> Simulate an error response for testing.</summary>

##### Parameters

> None

##### Responses

> | HTTP Code     | Response                               |
> |---------------|----------------------------------------|
> | `400`         | `Something went wrong :(`        |

</details>

### Student Routes
Base URL: `/v1/students`

<details>
 <summary><code>GET</code> <code><b>/healthz</b></code> Check the readiness of the student service.</summary>

##### Parameters

> None

##### Responses

> | HTTP Code     | Response                  |
> |---------------|---------------------------|
> | `200`         | `Service ready.`       |

</details>

<details>
 <summary><code>GET</code> <code><b>/err</b></code> Simulate an error response for testing.</summary>

##### Parameters

> None

##### Responses

> | HTTP Code     | Response                               |
> |---------------|----------------------------------------|
> | `400`         | `{"error": "Something went wrong :("}`         |

</details>

<details>
 <summary><code>POST</code> <code><b>/register</b></code> Register a new student.</summary>

##### Body Parameters

> | Name       | Type     | Data Type | Description                  |
> |------------|----------|-----------|------------------------------|
> | `username` | Required | String    | Student's username           |
> | `password` | Required | String    | Student's password           |
> | `name`     | Required | String    | Student's name               |
> | `email`    | Required | String    | Student's email address      |

##### Responses

> | HTTP Code     | Response                            |
> |---------------|-------------------------------------|
> | `201`         | `Registration successful.`          |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `409`         | `{"error": "Email already taken."}`              |
> | `409`         | `{"error": "Username already taken."}`           |
> | `500`         | `{"error": Internal server error."}`            |

</details>

<details>
 <summary><code>POST</code> <code><b>/login</b></code> Login as a registered student.</summary>

##### Body Parameters

> | Name       | Type     | Data Type | Description                  |
> |------------|----------|-----------|------------------------------|
> | `username` | Required | String    | Student's username           |
> | `password` | Required | String    | Student's password           |

##### Responses

> | HTTP Code     | Response                                                          |
> |---------------|-------------------------------------------------------------------|
> | `200`         | `{"token": jwt_token_string, "student": student_model}` |
> | `400`         | `{"error": "Missing one or more required parameters."}`            |
> | `401`         | `{"error": "Wrong password"}`                                      |
> | `500`         | `{"error": "Internal server error"}`                               |


</details>

<details>
 <summary><code>GET</code> <code><b>/profile</b></code> Retrieve the profile of the authenticated student.</summary>

##### Parameters

> None

##### Responses

> | HTTP Code     | Response                            |
> |---------------|-------------------------------------|
> | `200`         |​ `student_model`                      |
> | `500`         |​ `{"error": "Internal server error"}`    |


<details>
 <summary><code>PUT</code> <code><b>/profile</b></code> Update the profile of the authenticated student.</summary>

##### Parameters

> | Name       | Type     | Data Type | Description                  |
> |------------|----------|-----------|------------------------------|
> | Body       | Required | JSON      |                              |
> | `username` | Required | String    | Student's username           |
> | `name`     | Required | String    | Student's name               |
> | `email`    | Required | String    | Student's email address      |

##### Responses

> | HTTP Code     | Response                                 |
> |---------------|------------------------------------------|
> | `200`         | `updated student_model`                         |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `401`         | `{"error": "Authentication required."}`  |
> | `409`         | `{"error": "Username or email already taken."}`|
> | `500`         | `{"error": "Internal server error."}`    |

</details>

<details>
 <summary><code>POST</code> <code><b>/new-question</b></code> Submit a new question.</summary>

##### Parameters

> | Name         | Type     | Data Type | Description                   |
> |--------------|----------|-----------|-------------------------------|
> | Body         | Required | JSON      |                               |
> | `subject_id` | Required | Integer   | ID of the subject for the question |
> | `header`     | Required | String    | Header/title of the question  |
> | `photo_url`  | Optional | String    | URL of a photo related to the question (if any) |
> | `content`    | Required | String    | Content/body of the question  |

##### Responses

> | HTTP Code     | Response                                 |
> |---------------|------------------------------------------|
> | `201`         | `Question submitted successfully.`       |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `500`         | `{"error": "Internal server error."}`    |

</details>

### Tutor Routes
Base URL: `/v1/tutors`

<details>
 <summary><code>GET</code> <code><b>/healthz</b></code> Check the readiness of the tutor service.</summary>

##### Parameters

> None

##### Responses

> | HTTP Code     | Response                  |
> |---------------|---------------------------|
> | `200`         | `Service ready.`       |

</details>

<details>
 <summary><code>GET</code> <code><b>/err</b></code> Simulate an error response for testing.</summary>

##### Parameters

> None

##### Responses

> | HTTP Code     | Response                               |
> |---------------|----------------------------------------|
> | `400`         | `{"error": "Something went wrong :("}`         |

</details>

<details>
 <summary><code>POST</code> <code><b>/register</b></code> Register a new tutor.</summary>

##### Body Parameters

> | Name         | Type     | Data Type | Description                           |
> |--------------|----------|-----------|---------------------------------------|
> | Body         | Required | JSON      |                                       |
> | `username`   | Required | String    | Tutor's username                       |
> | `password`   | Required | String    | Tutor's password                       |
> | `name`       | Required | String    | Tutor's name                           |
> | `subjects`   | Required | Object    | Map of subjects and years of experience. Keys are subject ID, values are years of experience (integer). |
> | `email`      | Required | String    | Tutor's email address                  |

##### Responses

> | HTTP Code     | Response                            |
> |---------------|-------------------------------------|
> | `201`         | `Registration successful.`          |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `409`         | `{"error": "Email already taken."}`              |
> | `409`         | `{"error": "Username already taken."}`           |
> | `500`         | `{"error": Internal server error."}`            |

</details>

<details>
 <summary><code>POST</code> <code><b>/login</b></code> Login a registered tutor.</summary>

##### Body Parameters

> | Name       | Type     | Data Type | Description                  |
> |------------|----------|-----------|------------------------------|
> | `username` | Required | String    | Student's username           |
> | `password` | Required | String    | Student's password           |

##### Responses

> | HTTP Code     | Response                                                          |
> |---------------|-------------------------------------------------------------------|
> | `200`         | `{"token": jwt_token_string, "tutor": tutor_model}` |
> | `400`         | `{"error": "Missing one or more required parameters."}`            |
> | `401`         | `{"error": "Wrong password"}`                                      |
> | `500`         | `{"error": "Internal server error"}`                               |


</details>

<details>
 <summary><code>GET</code> <code><b>/profile</b></code> Retrieve the profile of the authenticated tutor.</summary>

##### Parameters

> None

##### Responses

> | HTTP Code     | Response                            |
> |---------------|-------------------------------------|
> | `200`         |​ `tutor_model`                      |
> | `500`         |​ `{"error": "Internal server error"}`    |


<details>
 <summary><code>PUT</code> <code><b>/profile</b></code> Update the profile of the authenticated tutor.</summary>

##### Body Parameters

> | Name       | Type     | Data Type | Description                  |
> |------------|----------|-----------|------------------------------|
> | `username` | Required | String    | Tutor's username           |
> | `name`     | Required | String    | Tutor's name               |
> | `email`    | Required | String    | Tutor's email address      |

##### Responses

> | HTTP Code     | Response                                 |
> |---------------|------------------------------------------|
> | `200`         | `updated tutor_model`                         |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `401`         | `{"error": "Authentication required."}`  |
> | `409`         | `{"error": "Username or email already taken."}`|
> | `500`         | `{"error": "Internal server error."}`    |

</details>

<details>
 <summary><code>PUT</code> <code><b>/student-profile</b></code> Retrieve the profile of a student for the authenticated tutor.</summary>

##### Body Parameters

> | Name       | Type     | Data Type | Description                  |
> |------------|----------|-----------|------------------------------|
> | `student_id` | Required | UUID    | Student's ID          |

##### Responses

> | HTTP Code     | Response                                 |
> |---------------|------------------------------------------|
> | `200`         | `updated tutor_model`                         |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `401`         | `{"error": "Authentication required."}`  |
> | `409`         | `{"error": "Username or email already taken."}`|
> | `500`         | `{"error": "Internal server error."}`    |

</details>

<details>
 <summary><code>GET</code> <code><b>/pending-questions</b></code> Retrieve available questions for the authenticated tutor.</summary>

##### Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | None  | Required | N/A       | No parameters required.          |

##### Responses

> | HTTP Code     | Response                                                  |
> |---------------|-----------------------------------------------------------|
> | `200`         | `JSON arry of chat_model (empty array if no available questions)`   |
> | `401`         | `{"error": "Authentication required."}`                    |
> | `500`         | `{"error": "Internal server error."}`                      |

</details>

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
