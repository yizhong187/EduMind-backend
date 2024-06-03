# EduMind-backend
This is the repo for the EduMind mobile application's backend.

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

3. Login
   - Route: `/v1/login`
   - Method: POST
   - Purpose: Authenticate a user and start a session.
   - Request Parameters:
     - Body:
       - 'name': name
       - 'password': password
   - Responses:
     - 200 OK: Login successful.
     - 401 Unauthorized: Invalid credentials.

4. Logout
   - Route: `/v1/logout`
   - Method: GET
   - Purpose: Logout the current user and end the session.
   - Request Parameters: None
   - Responses:
     - 200 OK: Logout successful.

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
   - Responses:
     - 201 Created: Registration successful.
     - 400 Bad Request: Invalid registration details.

4. Get Student Profile
   - Route: `/v1/students/profile`
   - Method: GET
   - Purpose: Retrieve the profile of the authenticated student.
   - Request Parameters: None
   - Responses:
     - 200 OK: Profile retrieved successfully.
     - 401 Unauthorized: Authentication required.

5. Update Student Profile
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
   - Route: `/v1/chat/{chatID}`
   - Method: GET
   - Middleware: MiddlewareChatAuth
   - Purpose: Retrieve all messages in a specific chat.
   - Request Parameters:
     - Path Parameter: chatID (ID of the chat session)
   - Responses:
     - 200 OK: Messages retrieved successfully.
     - 401 Unauthorized: Authentication required.
     - 404 Not Found: Chat not found.

4. Join Chat Room
   - Route: `/v1/chat/joinChat/{chatID}`
   - Method: GET
   - Middleware: MiddlewareChatAuth
   - Purpose: Join a chat room for WebSocket communication.
   - Request Parameters:
     - Path Parameter: chatID (ID of the chat session)
   - Responses:
     - 200 OK: Joined chat room successfully.
     - 401 Unauthorized: Authentication required.
     - 404 Not Found: Chat room not found.
