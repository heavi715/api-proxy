# API Proxy

This project is an API proxy server built using Go and the Gin framework. It is designed to forward requests to various platforms, handling authentication and request logging.

## Features

- Proxy requests to multiple platforms like OpenAI, Claude, Google, etc.
- IP filtering to allow requests only from specific IP ranges.
- Supports CORS for cross-origin requests.
- Configurable via a JSON configuration file.
- Logs requests and responses for monitoring and debugging.

## Prerequisites

- Go 1.16 or later
- A valid configuration file (`config.json`)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/api-proxy.git
   cd api-proxy
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the project:

   ```bash
   go build -o api-proxy
   ```

## Configuration

The server is configured using a `config.json` file. Below is an example configuration:

```json
{
    "server_addr": "0.0.0.0:8080",
    "source_list": ["prod", "dev", "test"],
    "server_key_list": [],
    "proxy_url": "http://127.0.0.1:7890",
    "platform_list": {
        "openai": {
            "name": "openai",
            "url": "https://api.openai.com",
            "header_key": "Authorization",
            "header_values": ["Bearer sk-proj-1234567890"]
        },
        ...
    }
}
```

## Usage

1. Start the server:

   ```bash
   ./api-proxy
   ```

2. The server will listen on the address specified in the configuration file (e.g., `0.0.0.0:8080`).

3. Make requests to the proxy endpoints, for example:

   ```bash
   curl -X GET http://localhost:8080/proxy/openai/prod/v1/your-endpoint
   ```

## Logging

The server logs requests and responses using a custom logger. Logs are printed to the console.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
