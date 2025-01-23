# OAuth and User Info Proxy Server  

This Go application acts as a proxy server to handle OAuth token requests and user information retrieval. It forwards incoming HTTP requests to an external OAuth server and returns the responses to the client. The server also supports Cross-Origin Resource Sharing (CORS) to allow requests from any origin.  

---

## Features  

- **CORS Support:** Handles preflight (OPTIONS) requests and allows POST and GET requests from any origin.  
- **OAuth Token Proxy:** Proxies POST requests to obtain an OAuth access token.  
- **User Info Proxy:** Proxies GET requests to retrieve user information.  
- **Environment Variables:** URLs for external services are configurable through a `.env` file.  

---

## Requirements  

- Go 1.18+  
- Dependency: [`godotenv`](https://github.com/joho/godotenv) for loading environment variables.  

---

## Installation  

1. **Clone the repository:**  
   ```bash  
   git clone https://github.com/kema-ray/sso-middleware.git  
   cd sso-middleware
   ```

2. **Install dependencies:**  
   ```bash
   go mod download
   ```

3. **Set up environment variables:**  
   Create a `.env` file in the project root with the following variables:
   ```
   OAUTH_TOKEN_URL=https://your-oauth-server.com/token
   USER_INFO_URL=https://your-oauth-server.com/userinfo
   ```

---

## Configuration

The application uses environment variables to configure external service endpoints:
- `OAUTH_TOKEN_URL`: The OAuth server's token endpoint
- `USER_INFO_URL`: The OAuth server's user information endpoint

---

## Running the Server

```bash
go run main.go
```

The server will start on `localhost:8080` with two proxy endpoints:
- `/proxy/oauth/token`: OAuth token request proxy
- `/proxy/user-info`: User information request proxy

---

## Endpoints

### OAuth Token Proxy `/proxy/oauth/token`
- **Method:** POST
- **Description:** Forwards OAuth token requests to the configured token endpoint
- **CORS:** Supports cross-origin requests

### User Info Proxy `/proxy/user-info`
- **Method:** GET
- **Description:** Retrieves user information using the provided access token
- **Parameters:** `access_token` (query parameter)
- **CORS:** Supports cross-origin requests

---

## Error Handling

The server provides appropriate HTTP status codes and error messages for various scenarios:
- Method Not Allowed (405)
- Bad Request (400)
- Internal Server Error (500)
- Bad Gateway (502)

---

## Security Considerations

- All requests are proxied through the server
- CORS is enabled with a wildcard origin (`*`)
- Access tokens are passed securely via Authorization header

---

<!-- ## License

[Add your license information here] -->