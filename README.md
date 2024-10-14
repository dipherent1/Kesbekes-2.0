
# Kesbekes 2.0

Kesbekes 2.0 is a Go-based application that leverages the Telegram Database Library (TDLib) to create a robust, fully functional Telegram bot. It integrates PostgreSQL for database management and relies on environment variables for configuration, allowing seamless setup across different environments.

## Features

- **Telegram Bot Integration**: Uses TDLib for advanced interactions with the Telegram API.
- **Database Management**: Connects to PostgreSQL, utilizing GORM for ORM capabilities.
- **Environment-Driven Configuration**: Easily configurable for different environments via `.env` files.
  
## Prerequisites

- [Go](https://golang.org/) (version 1.16 or higher recommended)
- [PostgreSQL](https://www.postgresql.org/) database
- [TDLib](https://core.telegram.org/tdlib)
- CMake, OpenSSL, zlib, and gperf for building TDLib dependencies

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/dipherent1/Kesbekes-2.0.git
   cd Kesbekes-2.0
   ```

2. **Install TDLib**:
   Follow these steps to install TDLib for Telegram API interactions:
   ```bash
   git clone https://github.com/tdlib/td.git --depth 1
   cd td
   mkdir build
   cd build
   cmake -DCMAKE_BUILD_TYPE=Release ..
   cmake --build . -- -j5
   sudo make install
   ```
   Configure your environment for TDLib access:
   ```bash
   export CGO_CFLAGS="-I/usr/local/include"
   export CGO_LDFLAGS="-L/usr/local/lib -ltdjson"
   ```

3. **Install Go Dependencies**:
   In the project directory, install the required Go packages:
   ```bash
   go mod tidy
   ```

4. **Set Up Environment Variables**:
   Kesbekes 2.0 relies on the following environment variables, which you can set in a `.env` file:
   ```bash
   DB_PORT=<your-db-port>
   DB_HOST=<your-db-host>
   DB_USER=<your-db-username>
   DB_NAME=<your-db-name>
   DB_PASSWORD=<your-db-password>
   
   API_ID=<your-telegram-api-id>
   API_HASH=<your-telegram-api-hash>
   PhoneNumber=<your-phone-number>
   BotToken=<your-bot-token>
   GeminAPIKey=<your-gemini-api-key>
   ```

5. **Run the Application**:
   Start the application with:
   ```bash
   go run main.go
   ```

## Usage

After running the application, the Telegram bot will listen for updates and execute configured actions via the webhook specified in `bot.go`. You can customize the webhook URL and bot behaviors as needed.

## Contributing

To contribute to Kesbekes 2.0:
1. Fork the repository.
2. Create a new branch:
   ```bash
   git checkout -b feature-name
   ```
3. Make your changes and commit:
   ```bash
   git commit -m "Add new feature"
   ```
4. Push to your branch:
   ```bash
   git push origin feature-name
   ```
5. Submit a pull request for review.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contact

For questions, please open an issue on GitHub or contact the repository maintainer directly.