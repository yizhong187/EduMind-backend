# EduMind-backend
This is the repo for the EduMind mobile application's backend.

## Stucture of Models

### Student

Represents a student in the system.

<details>
 <summary> Attributes </summary>

- **student_id** (`uuid.UUID`): The unique identifier for the student.
- **username** (`string`): The username chosen by the student.
- **email** (`string`): The email address of the student.
- **created_at** (`time.Time`): The timestamp when the student account was created.
- **name** (`string`): The full name of the student.
- **valid** (`bool`): Indicates if the student account is currently valid.
- **photo_url** (`string`): The URL of an optional profile photo associated with the student (nullable).

</details>

<details>
 <summary> Example JSON Representation </summary>

```json
{
  "student_id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "student123",
  "email": "student123@example.com",
  "created_at": "2024-07-15T12:00:00Z",
  "name": "John Doe",
  "valid": true,
  "photo_url": "https://res.cloudinary.com/dnc1q8tlu/image/upload/v1720609760/file_kyd4yc.jpg"
}
```

</details>

### Subject (Tutor's specialisation)

Represents a subject in which a tutor specializes.

<details>
 <summary> Attributes </summary>

- **subject_id** (`int32`): The ID of the subject.
- **yoe** (`int32`): The years of experience the tutor has in teaching this subject.

</details>

<details>
 <summary> Example JSON Representation </summary>

```json
{
  "subject_id": "1",
  "yoe": 5
}
```

</details>

### Tutor

Represents a tutor in the system.

<details>
 <summary> Attributes </summary>

- **tutor_id** (`uuid.UUID`): The unique identifier for the tutor.
- **username** (`string`): The username chosen by the tutor.
- **email** (`string`): The email address of the tutor.
- **created_at** (`time.Time`): The timestamp when the tutor account was created.
- **name** (`string`): The full name of the tutor.
- **valid** (`bool`): Indicates if the tutor account is currently valid.
- **subjects** (array of `Subject (Tutor's specialisation)`): An array of subjects that the tutor specializes in.
- **verified** (`bool`): Indicates whether the tutor’s account is verified (`true`) or not (`false`).
- **rating** (`float64`): The average rating given to the tutor.
- **rating_count** (`int32`): The total number of ratings received by the tutor.
- **photo_url** (`string`): The URL of an optional profile photo associated with the tutor (nullable).

</details>

<details>
 <summary> Example JSON Representation </summary>

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
      "subject_id": 2,
      "yoe": 5
    },
    {
      "subject_id": 1,
      "yoe": 3
    }
  ],
  "verified": true,
  "rating": 4.7,
  "rating_count": 35,
  "photo_url": "https://res.cloudinary.com/dnc1q8tlu/image/upload/v1720609760/file_kyd4yc.jpg"
}
```

</details>

### Chat

Represents a chat session between a student and a tutor.

<details>
 <summary> Attributes </summary>

- **chat_id** (`int32`): The unique identifier for the chat session.
- **student_id** (`uuid.UUID`): The unique identifier of the student participating in the chat.
- **tutor_id** (`uuid.UUID`): The unique identifier of the tutor participating in the chat (nullable).
- **created_at** (`time.Time`): The timestamp when the chat session was created.
- **subject_id** (`int32`): The identifier of the subject associated with the chat.
- **topics** (array of `int32`): The topics of the chat session (nullable).
- **header** (`string`): A brief header or title for the chat session.
- **photo_url** (`string`): The URL of an optional photo associated with the chat (nullable).
- **completed** (`bool`): Indicates whether the chat session is completed (`true`) or ongoing (`false`).

</details>

<details>
 <summary> Example JSON Representation </summary>

```json
{
  "chat_id": 12345,
  "student_id": "550e8400-e29b-41d4-a716-446655440000",
  "tutor_id": "352bd79d-f432-48a6-9607-3b89ac7d1452",
  "created_at": "2024-07-15T12:00:00Z",
  "subject_id": 1,
  "topics": [12,34,56],
  "header": "Stoichiometry for Redox Reactions",
  "photo_url": "https://res.cloudinary.com/dnc1q8tlu/image/upload/v1720609760/file_kyd4yc.jpg",
  "completed": false
    }
```

</details>

### Message

Represents a message within a chat session.

<details>
 <summary> Attributes </summary>

- **message_id** (`uuid.UUID`): The unique identifier for the message.
- **chat_id** (`int32`): The identifier for the chat session to which the message belongs.
- **user_id** (`uuid.UUID`): The unique identifier of the user who sent the message.
- **created_at** (`time.Time`): The timestamp when the message was created.
- **updated_at** (`time.Time`): The timestamp when the message was last updated.
- **deleted** (`bool`): Indicates if the message is deleted (`true`) or not (`false`).
- **content** (`string`): The content of the message.

</details>

<details>
 <summary> Example JSON Representation </summary>

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

</details>

### Subject (ID-Name)

Maps the ID to the name of a subject.

<details>
 <summary> Attributes </summary>

- **subject_id** (`int32`): The ID of the subject.
- **name** (`string`): The name of the subject.

</details>

<details>
 <summary> Example JSON Representation </summary>

```json
{
  "subject_id": "1",
  "name": "Chemistry"
}
```

</details>

### Topic (SubjectID-TopicID-Name)

Maps the subject ID and topic ID to the name of a subject.

<details>
 <summary> Attributes </summary>

- **subject_id** (`int32`): The ID of the subject that the topic belongs to.
- **topic_id** (`int32`): The ID of the topic.
- **name** (`string`): The name of the topic.

</details>

<details>
 <summary> Example JSON Representation </summary>

```json
{
  "subject_id": "1",
  "topic_id": "22",
  "name": "Nitrogen Compounds"
}
```
</details>

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
> | `200`         | `"Service ready"`                                |

</details>

<details>
 <summary><code>GET</code> <code><b>/subjects</b></code> Retrieves all subjects within the database. </summary>

##### Parameters

> None

##### Responses

> | HTTP Code     | Response                                                  |
> |---------------|-----------------------------------------------------------|
> | `200`         | `Array of Subject (ID-Name)`                                       |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `401`         | `{"error": "Authentication required."}`                    |
> | `500`         | `{"error": "Internal server error."}`                      |

</details>

<details>
 <summary><code>GET</code> <code><b>/topics</b></code> Retrieves all topics within the database.</summary>

##### Parameters

> None

##### Responses

> | HTTP Code     | Response                                                  |
> |---------------|-----------------------------------------------------------|
> | `200`         | `Array of Topic (SubjectID-Topic-ID-Name)`         |
> | `500`         | `{"error": "Internal server error."}`                      |

</details>

<details>
 <summary><code>GET</code> <code><b>/{subjectID}</b></code> Retrieves all topics of a subject within the database.</summary>

##### Path Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | `subjectID` | Required | Integer    | ID of subject.         |

##### Responses

> | HTTP Code     | Response                                                  |
> |---------------|-----------------------------------------------------------|
> | `200`         | `Array of Topic (SubjectID-Topic-ID-Name)`         |
> | `500`         | `{"error": "Internal server error."}`                      |

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
> | `200`         | `"Service ready."`       |

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
> | `username` | Required | String    | Student's username.           |
> | `password` | Required | String    | Student's password.           |
> | `name`     | Required | String    | Student's name.               |
> | `email`    | Required | String    | Student's email address.      |
> | `photo_url`    | Optional | String    | URL of student's profile photo (if any). |

##### Responses

> | HTTP Code     | Response                            |
> |---------------|-------------------------------------|
> | `201`         | `"Registration successful."`          |
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
> | `username` | Required | String    | Student's username.           |
> | `password` | Required | String    | Student's password.           |

##### Responses

> | HTTP Code     | Response                                                          |
> |---------------|-------------------------------------------------------------------|
> | `200`         | `{"token": jwt_token_string, "student": student_model}` |
> | `400`         | `{"error": "Missing one or more required parameters."}`            |
> | `401`         | `{"error": "Wrong password"}`                                      |
> | `500`         | `{"error": "Internal server error"}`                               |


</details>

<details>
 <summary><code>GET</code> <code><b>/profile/{student_username}</b></code> Retrieve the profile of a student by username.</summary>

##### Path Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | `student_username` | Required | String    | Student's username.         |

##### Responses

> | HTTP Code     | Response                            |
> |---------------|-------------------------------------|
> | `200`         |​ `student_model`                      |
> | `404`         |​ `{"error": "Student profile not found"}`    |
> | `500`         |​ `{"error": "Internal server error"}`    |

</details>

<details>
 <summary><code>GET</code> <code><b>/profile</b></code> Retrieve the profile of a student by ID.</summary>

##### Body Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | `tutor_id` | Required | UUID    | Student's ID.         |

##### Responses

> | HTTP Code     | Response                            |
> |---------------|-------------------------------------|
> | `200`         |​ `student_model`                      |
> | `404`         |​ `{"error": "Student profile not found"}`    |
> | `500`         |​ `{"error": "Internal server error"}`    |

</details>

<details>
 <summary><code>PUT</code> <code><b>/profile</b></code> Update the profile of the authenticated student.</summary>

##### Body Parameters

> | Name       | Type     | Data Type | Description                  |
> |------------|----------|-----------|------------------------------|
> | `username` | Required | String    | Student's new/updated username.           |
> | `name`     | Required | String    | Student's new/updated name.               |
> | `email`    | Required | String    | Student's new/updated email address.      |

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
 <summary><code>PUT</code> <code><b>/update-password</b></code> Update the password of the authenticated student.</summary>

##### Body Parameters

> | Name       | Type     | Data Type | Description                  |
> |------------|----------|-----------|------------------------------|

> | `old_password`     | Required | String    | Student's current password.               |
> | `new_password`    | Required | String    | Student's updated password.      |

##### Responses

> | HTTP Code     | Response                                 |
> |---------------|------------------------------------------|
> | `200`         | `"Password updated successfully"`                         |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `401`         | `{"error": "Incorrect password."}`  |
> | `500`         | `{"error": "Internal server error."}`    |

</details>

<details>
 <summary><code>POST</code> <code><b>/new-question</b></code> Submit a new question.</summary>

##### Body Parameters

> | Name         | Type     | Data Type | Description                   |
> |--------------|----------|-----------|-------------------------------|
> | `subject_id` | Required | Integer   | ID of the subject for the question. |
> | `header`     | Required | String    | Header/title of the question.  |
> | `photo_url`  | Optional | String    | URL of a photo related to the question (if any). |
> | `content`    | Required | String    | Content/body of the question.  |

##### Responses

> | HTTP Code     | Response                                 |
> |---------------|------------------------------------------|
> | `201`         | `"Question submitted successfully."`       |
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
> | `200`         | `"Service ready."`       |

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
> | `username`   | Required | String    | Tutor's username.                       |
> | `password`   | Required | String    | Tutor's password.                       |
> | `name`       | Required | String    | Tutor's name.                           |
> | `subjects`   | Required | Object    | Map of subjects and years of experience. Keys are subject ID, values are years of experience (integer). |
> | `email`      | Required | String    | Tutor's email address.                 |
> | `photo_url`    | Optional | String    | URL of student's profile photo (if any). |

##### Responses

> | HTTP Code     | Response                            |
> |---------------|-------------------------------------|
> | `201`         | `"Registration successful."`          |
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
> | `username` | Required | String    | Tutor's username.           |
> | `password` | Required | String    | Tutor's password.           |

##### Responses

> | HTTP Code     | Response                                                          |
> |---------------|-------------------------------------------------------------------|
> | `200`         | `{"token": jwt_token_string, "tutor": tutor_model}` |
> | `400`         | `{"error": "Missing one or more required parameters."}`            |
> | `401`         | `{"error": "Wrong password"}`                                      |
> | `500`         | `{"error": "Internal server error"}`                               |


</details>

<details>
 <summary><code>GET</code> <code><b>/profile/{tutor_username}</b></code> Retrieve the profile of a tutor by username.</summary>

##### Path Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | `tutor_username` | Required | String    | Tutor's username.         |

##### Responses

> | HTTP Code     | Response                            |
> |---------------|-------------------------------------|
> | `200`         |​ `tutor_model`                      |
> | `404`         |​ `{"error": "Tutor profile not found"}`    |
> | `500`         |​ `{"error": "Internal server error"}`    |

</details>

<details>
 <summary><code>GET</code> <code><b>/profile</b></code> Retrieve the profile of a tutor by ID.</summary>

##### Body Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | `tutor_id` | Required | UUID    | Tutor's ID.         |

##### Responses

> | HTTP Code     | Response                            |
> |---------------|-------------------------------------|
> | `200`         |​ `tutor_model`                      |
> | `404`         |​ `{"error": "Tutor profile not found"}`    |
> | `500`         |​ `{"error": "Internal server error"}`    |

</details>

<details>
 <summary><code>PUT</code> <code><b>/profile</b></code> Update the profile of the authenticated tutor.</summary>

##### Body Parameters

> | Name       | Type     | Data Type | Description                  |
> |------------|----------|-----------|------------------------------|
> | `username` | Required | String    | Tutor's new/unchanged username.           |
> | `name`     | Required | String    | Tutor's new/unchanged name.               |
> | `email`    | Required | String    | Tutor's new/unchanged email address.      |

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
 <summary><code>PUT</code> <code><b>/update-password</b></code> Update the password of the authenticated tutor.</summary>

##### Body Parameters

> | Name       | Type     | Data Type | Description                  |
> |------------|----------|-----------|------------------------------|

> | `old_password`     | Required | String    | Tutor's current password.               |
> | `new_password`    | Required | String    | Tutor's updated password.      |

##### Responses

> | HTTP Code     | Response                                 |
> |---------------|------------------------------------------|
> | `200`         | `"Password updated successfully"`                         |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `401`         | `{"error": "Incorrect password."}`  |
> | `500`         | `{"error": "Internal server error."}`    |

</details>

<details>
 <summary><code>GET</code> <code><b>/student-profile</b></code> Retrieve the profile of a student for the authenticated tutor.</summary>

##### Body Parameters

> | Name       | Type     | Data Type | Description                  |
> |------------|----------|-----------|------------------------------|
> | `student_id` | Required | UUID    | Student's ID.          |

##### Responses

> | HTTP Code     | Response                                 |
> |---------------|------------------------------------------|
> | `200`         | `updated tutor_model`                         |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `401`         | `{"error": "Authentication required."}`  |
> | `409`         | `{"error": "Username or email already taken."}`|
> | `500`         | `{"error": "Internal server error."}`    |

</details>

### General Chat Routes
Base URL: `/v1/chats`

<details>
 <summary><code>GET</code> <code><b>/</b></code> Retrieve all chats for the authenticated user.</summary>

##### Parameters

> None

##### Responses

> | HTTP Code     | Response                                                  |
> |---------------|-----------------------------------------------------------|
> | `200`         | `JSON arry of chat_model (null if no chats)`   |
> | `401`         | `{"error": "Authentication required."}`                    |
> | `500`         | `{"error": "Internal server error."}`                      |

</details>

<details>
 <summary><code>GET</code> <code><b>/{chat_id}</b></code> Retrieve all messages of a specific chat for the authenticated user.</summary>

##### Path Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | `chat_id`  | Required | Integer       | Chat's ID.          |

##### Responses

> | HTTP Code     | Response                                                  |
> |---------------|-----------------------------------------------------------|
> | `200`         | `JSON array of message_model`   |
> | `401`         | `{"error": "Authentication required."}`                    |
> | `500`         | `{"error": "Internal server error."}`                      |

</details>

<details>
 <summary><code>POST</code> <code><b>/{chat_id}</b></code> Send a new message into a specific chat for the authenticated user.</summary>

##### Path Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | `chat_id`  | Required | Integer       | Chat's ID.          |

##### Body Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | `content`  | Required | String       | Content of message.          |

##### Responses

> | HTTP Code     | Response                                                  |
> |---------------|-----------------------------------------------------------|
> | `200`         | `"Message sent."`   |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `401`         | `{"error": "Authentication required."}`                    |
> | `500`         | `{"error": "Internal server error."}`                      |

</details>

### Tutor Specific Chat Routes
Base URL: `/v1/chats`

<details>
 <summary><code>POST</code> <code><b>/pending</b></code> Retrieve all available questions for the authenticated tutor. </summary>

##### Parameters

> None

##### Responses

> | HTTP Code     | Response                                                  |
> |---------------|-----------------------------------------------------------|
> | `200`         | `JSON array of chat_model (null if no chats)`                                       |
> | `401`         | `{"error": "Authentication required."}`                    |
> | `500`         | `{"error": "Internal server error."}`                      |

</details>

<details>
 <summary><code>POST</code> <code><b>/{chat_id}/accept</b></code> Accept an available question for the authenticated tutor. </summary>

##### Path Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | `chat_id` | Required | Integer    | Chat ID of question.         |

##### Responses

> | HTTP Code     | Response                                                  |
> |---------------|-----------------------------------------------------------|
> | `200`         | `"Question accepted."`                                       |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `401`         | `{"error": "Authentication required."}`                    |
> | `500`         | `{"error": "Internal server error."}`                      |

</details>

<details>
 <summary><code>PUT</code> <code><b>/{chat_id}/update-topics</b></code> Update the topics of the specific question for the authenticated tutor. </summary>

##### Path Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | `chat_id` | Required | Integer    | Chat ID of question.         |

##### Body Parameters

> | Name  | Type     | Data Type | Description                     |
> |-------|----------|-----------|---------------------------------|
> | `topics`  | Required | Integer Array       | Updated topics (topic IDs) of the question.       |

##### Responses

> | HTTP Code     | Response                                                  |
> |---------------|-----------------------------------------------------------|
> | `200`         | `updated chat_model`                                       |
> | `400`         | `{"error": "Missing one or more required parameters."}`|
> | `401`         | `{"error": "Authentication required."}`                    |
> | `500`         | `{"error": "Internal server error."}`                      |

</details>


## Deployment to Heroku
Use `env GOOS=linux GOARCH=amd64 GOARM=7 go build` to compile into linux based binary code for heroku
