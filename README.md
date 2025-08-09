# MatchPoint

[](https://opensource.org/licenses/MIT)
[](https://golang.org)
[](https://postgresql.org)

**MatchPoint** is an online sports court booking system that allows users to book sports facilities easily and efficiently. The system is built using **Clean Architecture** for optimal maintainability and high scalability.

-----


## Features

  - **Court Management**: Manage various types of sports courts
  - **Booking System**: Real-time reservation system with conflict detection
  - **User Management**: Registration, login, and profile management
  - **Payment Integration**: Support for multiple payment gateways
  - **Dashboard**: Admin dashboard for monitoring and analytics
  - **Notifications**: Email and push notifications for booking updates
  - **Mobile Responsive**: Optimized for both mobile and desktop
  - **Authentication**: JWT-based authentication with role-based access

-----

##  Tech Stack

  - **Backend**: Go (Golang)
  - **Database**: PostgreSQL
  - **Cache**: Redis (optional)
  - **Logging**: Structured logging with logrus/zap
  - **API Documentation**: Swagger/OpenAPI 3.0
  - **Testing**: Testify, GoMock
  - **Migration**: golang-migrate
  - **Containerization**: Docker & Docker Compose

-----

## 🚀 Getting Started

### Prerequisites

  - Go 1.21 or newer
  - PostgreSQL 13+
  - Docker & Docker Compose (optional)
  - Make (optional, for running Makefile commands)

### Installation

1.  **Clone the repository**

    ```bash
    git clone https://github.com/yourusername/matchpoint.git
    cd matchpoint
    ```

2.  **Install dependencies**

    ```bash
    go mod tidy
    ```

3.  **Set up environment variables**

    ```bash
    cp .env.example .env
    # Edit the .env file with your specific configuration
    ```

4.  **Set up the database**

    ```bash
    # Create PostgreSQL database
    createdb matchpoint_db
    ```

5.  **Run migrations**

    ```bash
    go run main.go --migrate
    ```

6.  **Run the application**

    ```bash
    go run main.go
    # or
    air
    ```

-----


## 📝 Environment Variables

Create a `.env` file based on `.env.example`:

```env
# Database Configuration
DB_HOST=localhost
DB_USER=root
DB_PASS=
DB_NAME=database_name
DB_PORT=3306

# Application Settings
APP_TIMEZONE=Asia/Jakarta
APP_ENV=development

# Server Ports
PORT=8000
NGINX_PORT=8080
GOLANG_PORT=8888

# JWT Secret Key
# Change this to a random string of at least 32 characters
JWT_SECRET=your_jwt_secret_key

# SMTP Mailer Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_SENDER_NAME="Go.Gin.Template <no-reply@yourdomain.com>"
SMTP_AUTH_EMAIL=your_email@gmail.com
SMTP_AUTH_PASSWORD=your_email_password
```

-----


## 📊 Logging

The application uses structured logging with a JSON format for production and a human-readable format for development.

### Log Levels

  - **ERROR**: Errors that require immediate attention
  - **WARN**: Warnings that should be noted
  - **INFO**: General application information
  - **DEBUG**: Detailed information for debugging

-----


## 🤝 Contributing

1.  Fork the repository
2.  Create your feature branch (`git checkout -b feature/amazing-feature`)
3.  Commit your changes (`git commit -m 'Add amazing feature'`)
4.  Push to the branch (`git push origin feature/amazing-feature`)
5.  Open a Pull Request

### Coding Standards

  - Follow Go coding conventions
  - Use `gofmt` for formatting
  - Write unit tests for every new feature
  - Update documentation as needed

-----

## 📄 License

Distributed under the MIT License. See the `LICENSE` file for more information.

-----

## 👥 Team

  - **Maulana Ferdiansyah** - [mferdian](https://github.com/mferdian)

-----

## 🙏 Acknowledgments

  - Clean Architecture by Ahmad Mirza - [Amierza](https://github.com/Amierza)
  - The Go community for amazing libraries
  - The PostgreSQL team for a robust database system

-----

## 📞 Support

If you have questions or issues, please create a [GitHub issue](https://github.com/mferdian/matchpoint/issues) or contact:

  - Email: eikhapoetra@gmail.com
  - LinkedIn: [Maulana Ferdiansyah]([https://linkedin.com/in/yourprofile](https://www.linkedin.com/in/maulana-ferdiansyah-eka-putra-08a4b0289?lipi=urn%3Ali%3Apage%3Ad_flagship3_profile_view_base_contact_details%3BRfKQvp50RDGxgDAGkUWFFg%3D%3D))

-----

⭐ **Don't forget to star this repository if you find it helpful\!**
