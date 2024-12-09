# ZopDev

Zop is a comprehensive tool for managing cloud infrastructure. It consists of three main components:

1. **zop-cli**: Command-line interface for developers and admins.
2. **zop-api**: Backend API service.
3. **zop-ui**: User interface for managing and monitoring cloud resources.

---

## Installation

### Prerequisites

- Docker installed on your system.

---

### Running Locally

#### zop-api
Run the following command to pull and start the Docker image for the zop-api:

```bash
    docker run -d -p 8000:8000 --name zop-api zop.dev/zop-api:v0.0.1
```

#### zop-ui
Run the following command to pull and start the Docker image for the zop-ui:
```bash
    docker run -d -p 3000:3000 -e NEXT_PUBLIC_API_BASE_URL='http://localhost:8000' --name zop-ui zop.dev/zop-ui:v0.0.1
```

> **Note:** The environment variable `NEXT_PUBLIC_API_BASE_URL` is used by zop-ui to connect to the zop-api. Ensure that the value matches the API's running base URL.
#### zop-cli

Run the following command install zop-cli:
```bash
   go install zop.dev/clizop@latest
```

> **Note:** Set the environment variable `ZOP_API_URL`, used by zop-cli to connect to the zop-api. Ensure that the value matches the API's running base URL.

### zop-api

#### Commands

1. **cloud import**  
   Imports all the cloud accounts present on the local system to the zop-api.

   ```bash
    zop cloud import
    ```