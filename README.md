# ZopDev

Zop is a comprehensive tool for managing cloud infrastructure. It consists of three main components:

1. **zop-cli**: Command-line interface for developers and admins.
2. **zop-api**: Backend API service.
3. **zop-ui**: User interface for managing and monitoring cloud resources.

---

## Installation

### Prerequisites

- Docker installed on user system.

---

### Running Locally

#### zop-api
Run the following command to pull and start the Docker image for the zop-api:

```bash
    docker run -d -p 8000:8000 --name zop-api zopdev/zop-api:v0.0.2
```

#### zop-ui
Run the following command to pull and start the Docker image for the zop-ui:
```bash
    docker run -d -p 3000:3000 -e NEXT_PUBLIC_API_BASE_URL='http://localhost:8000' --name zop-ui zopdev/zop-ui:v0.0.2
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
2. **cloud list**  
   Lists all the cloud accounts present in the zop-api.

   ```bash
    zop cloud list
    ```
3. **application add -name=<app_name>**

   Adds a new application to the zop-api. This lets users add environment is ascending order of
   their continuous delivery sequence.
    
    ```bash
     zop application add -name=<app_name>
     ```
4. **application list**
   
   Lists all the applications present in the zop-api for a selected application.

    ```bash
     zop application list
     ```
5. **environment add**

   Adds a new environment to the zop-api. This lets user add deployment in ascending order of
   their continuous delivery sequence. Users can add multiple environments to an application.

    ```bash
     zop environment add
     ```
6. **environment list**

   Lists all the environments present in the zop-api for a selected application.

    ```bash
     zop environment list
     ```
   
7. **deployment add**

   Adds a new deployment to the zop-api. The users are needed to select cloud-account and the application
   environment where the deployment space is needed to be configured. Then users can select from a list of
   available options(ex, GKE cluster, AWS EC2 instance, etc.) to deploy their application.

    ```bash
     zop deployment add
     ```