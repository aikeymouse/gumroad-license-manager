# Gumroad License Manager

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com)

A comprehensive web application for managing Gumroad products and license keys with advanced API call tracking, license validation, and real-time monitoring capabilities.

## 🚀 Features

### Core Functionality
- **Product Management**: View all Gumroad products with detailed information
- **License Key Management**: Browse license keys and sales data for products
- **License Validation**: Real-time license key validation with detailed results
- **API Call Monitoring**: Complete logging and tracking of all Gumroad API interactions

### Advanced Features
- **Intelligent Navigation**: Dynamic back navigation that remembers your previous page
- **Real-time Search**: Instant filtering across products and license data
- **Modal Dialogs**: Detailed API call inspection with request/response data
- **Duration Tracking**: Precise timing of API calls with millisecond accuracy
- **Error Handling**: Comprehensive error logging and user feedback
- **Professional UI**: Clean, responsive interface with hover effects and animations

### Technical Highlights
- **Zero External Dependencies**: Pure Go backend with vanilla JavaScript frontend
- **Template-based Architecture**: Modular HTML templates with inheritance
- **Static Asset Organization**: Optimized CSS and JavaScript file structure
- **Docker Support**: Containerized deployment with multi-stage builds
- **Configuration Management**: JSON-based configuration with example templates

## 📁 Project Structure

```
gumroad-license-manager/
├── main.go                    # Core application server
├── go.mod                     # Go module dependencies  
├── config.json               # Configuration file
├── config.example.json       # Example configuration
├── templates/                # HTML templates
│   ├── base.html            # Base layout with navigation
│   ├── products.html        # Products listing page
│   ├── licenses.html        # License keys page
│   ├── sales.html           # Sales data page
│   ├── api-log.html         # API call monitoring
│   └── setup.html           # Initial setup page
├── static/                   # Static web assets
│   ├── css/
│   │   └── style.css        # Application styles (900+ lines)
│   └── js/
│       ├── app.js           # General functionality
│       └── api-log-modal.js # Modal dialog handling
├── Dockerfile               # Docker image definition
├── docker-compose.yml       # Docker Compose setup
└── README.md               # This documentation
```

## 🛠️ Setup & Installation

### Prerequisites
- Go 1.19+ (for local development)
- Docker & Docker Compose (for containerized deployment)
- Gumroad account with API access

### 1. Get Your Gumroad API Token

1. Visit [Gumroad API Settings](https://gumroad.com/api)
2. Generate an access token
3. Copy the token for configuration

### 2. Configure the Application

Create `config.json` from the example:

```json
{
  "gumroad_token": "your-actual-gumroad-access-token-here"
}
```

### 3. Choose Your Deployment Method

#### Option A: Docker Compose (Recommended)

```bash
# Clone or navigate to the project directory
cd gumroad-license-manager

# Start with Docker Compose
docker-compose up --build

# Application will be available at http://localhost:8086
```

#### Option B: Local Development

```bash
# Install Go dependencies
go mod tidy

# Run the application
go run main.go

# Server starts on http://localhost:8086
```

### 4. Initial Setup

1. Open http://localhost:8086 in your browser
2. If no configuration is found, you'll see a setup page
3. Enter your Gumroad API token
4. Click "Save Configuration" to complete setup

## 🎯 How to Use

### Main Dashboard
- **Products Tab**: View all your Gumroad products
- **API Call Log Tab**: Monitor all API interactions in real-time

### Product Management
1. Navigate to the main page to see all products
2. Products display with names, descriptions, and pricing
3. Click on any product to view its license keys or sales data

### License Key Validation
1. On any product page, find the "Validate License Key" section
2. Enter a license key in the input field
3. Click "Validate" to check the key's status
4. View detailed validation results including:
   - Validity status
   - Usage count
   - Purchase information
   - Buyer details

### API Call Monitoring
1. Click "API Call Log" in the navigation
2. View all API calls with:
   - **Time**: When the call was made
   - **Method**: HTTP method (GET, POST)
   - **URL**: Gumroad API endpoint
   - **Status**: Response code
   - **Duration**: Request time in milliseconds
   - **Error**: Any error messages

3. Click any row to view detailed information:
   - Complete request and response data
   - Headers and timing information
   - Error details if applicable

### Navigation Features
- **Smart Back Button**: Returns to your previous page
- **Breadcrumb Navigation**: Top-level tabs for easy switching
- **Modal Dialogs**: Detailed views without page reloads

## 🔧 API Endpoints

### Gumroad API Integration
The application integrates with these Gumroad endpoints:

- `GET /v2/products` - Fetch all products
- `GET /v2/products/{product_id}/subscribers` - Get license keys
- `GET /v2/sales?product_id={product_id}` - Get sales data
- `POST /v2/licenses/verify` - Validate license keys

### Internal API Endpoints
- `GET /` - Main products dashboard
- `GET /licenses/{product_id}` - License keys for product
- `GET /sales/{product_id}` - Sales data for product
- `GET /api-log` - API call monitoring page
- `GET /api/api-calls` - JSON API for call data
- `POST /validate-license` - License validation endpoint
- `GET /setup` - Initial configuration page
- `POST /setup` - Save configuration

## ⚙️ Configuration

### Environment Variables
- `PORT` - Server port (default: 8086)

### Configuration File (`config.json`)
```json
{
  "gumroad_token": "your-gumroad-access-token"
}
```

### Features Configuration
- **API Rate Limiting**: Built-in request throttling
- **Error Handling**: Comprehensive error logging and user feedback
- **Cache Management**: Optimized data loading and caching
- **Security**: Token validation and secure configuration handling

## 🐳 Docker Deployment

The application uses a multi-stage Docker build for optimal performance:

### Development
```bash
docker-compose up --build
```

### Production
```bash
# Build production image
docker build -t gumroad-license-manager .

# Run container
docker run -p 8086:8086 -v $(pwd)/config.json:/app/config.json gumroad-license-manager
```

## 📊 Monitoring & Logging

### API Call Tracking
- **Automatic Logging**: All Gumroad API calls are logged
- **Performance Metrics**: Response times and success rates
- **Error Tracking**: Detailed error messages and stack traces
- **Historical Data**: Last 100 API calls stored in memory

### License Validation Logging
- **Validation Attempts**: All license validation requests
- **Success/Failure Rates**: Track validation patterns
- **Usage Analytics**: Monitor license key usage

## 🎨 User Interface

### Design Features
- **Responsive Design**: Works on desktop and mobile devices
- **Professional Styling**: Clean, modern interface
- **Interactive Elements**: Hover effects and smooth animations
- **Status Indicators**: Visual feedback for API states
- **Monospace Typography**: Technical data displayed clearly

### Accessibility
- **Keyboard Navigation**: Full keyboard support
- **Screen Reader Friendly**: Semantic HTML structure
- **High Contrast**: Clear visual hierarchy
- **Error Messages**: User-friendly error communication

## 🚀 Development

### Local Development
```bash
# Install dependencies
go mod tidy

# Run with hot reload (if using air)
air

# Or run directly
go run main.go
```

### File Structure for Development
- **Backend Logic**: Edit `main.go`
- **HTML Templates**: Modify files in `templates/`
- **Styling**: Update `static/css/style.css`
- **JavaScript**: Edit files in `static/js/`

### Key Development Areas
1. **Template Functions**: Custom Go template functions for formatting
2. **API Integration**: Gumroad API client and error handling
3. **Frontend Interactions**: Modal dialogs and form handling
4. **Configuration Management**: Setup and token validation

## 📋 Troubleshooting

### Common Issues

1. **"Invalid Token" Error**
   - Verify your Gumroad API token is correct
   - Check token permissions in Gumroad dashboard

2. **"No Products Found"**
   - Ensure your Gumroad account has products
   - Verify API token has product access permissions

3. **License Validation Fails**
   - Check that the license key format is correct
   - Verify the product ID matches the license

4. **API Call Log Empty**
   - API calls are logged after they occur
   - Try refreshing products or validating a license

### Performance Tips
- **Docker**: Use Docker Compose for consistent deployment
- **Memory**: Application keeps 100 recent API calls in memory
- **Caching**: Templates and static assets are cached efficiently

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

The MIT License allows you to:
- ✅ Use the software for any purpose
- ✅ Modify and distribute the software
- ✅ Include in commercial projects
- ✅ Private use

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
