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

- Go 1.19+  
- Dependency: [`godotenv`](https://github.com/joho/godotenv) for loading environment variables.  

---

## Installation  

1. **Clone the repository:**  
   ```bash  
   git clone <repository-url>  
   cd <repository-directory>  
