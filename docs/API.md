# TourHelper API Documentation

## Overview
TourHelper API is a RESTful service for managing and suggesting tour destinations.

## Base URL
```
http://localhost:8080
```

## Authentication
Currently, the API does not require authentication. This may be added in future versions.

## Endpoints

### Health Checks

#### GET /health
Check the overall health of the service.

**Response:**
```json
{
  "status": "healthy",
  "database": "healthy",
  "version": "1.0.0"
}
```

#### GET /ready
Check if the service is ready to accept requests.

**Response:**
```json
{
  "success": true,
  "message": "Service is ready"
}
```

### Tours

#### POST /api/v1/tours
Create a new tour destination.

**Request Body:**
```json
{
  "name": "Tokyo Cherry Blossom Tour",
  "description": "Experience the beautiful cherry blossoms in Tokyo",
  "location": "Tokyo",
  "country": "Japan",
  "category": "cultural",
  "duration": 5,
  "season": "spring",
  "budget": "high",
  "image_url": "https://example.com/tokyo.jpg",
  "rating": 4.8
}
```

**Response:** (201 Created)
```json
{
  "success": true,
  "message": "Tour created successfully",
  "data": {
    "id": 1,
    "name": "Tokyo Cherry Blossom Tour",
    "description": "Experience the beautiful cherry blossoms in Tokyo",
    "location": "Tokyo",
    "country": "Japan",
    "category": "cultural",
    "duration": 5,
    "season": "spring",
    "budget": "high",
    "image_url": "https://example.com/tokyo.jpg",
    "rating": 4.8,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /api/v1/tours
List all tours with optional filters.

**Query Parameters:**
- `category` (optional): Filter by category (e.g., "beach", "mountain", "cultural", "city")
- `country` (optional): Filter by country
- `budget` (optional): Filter by budget ("low", "medium", "high")
- `season` (optional): Filter by best season
- `min_rating` (optional): Minimum rating (0-5)
- `limit` (optional): Number of results per page (default: 20)
- `offset` (optional): Pagination offset

**Example Request:**
```
GET /api/v1/tours?category=beach&budget=medium&limit=10
```

**Response:** (200 OK)
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Bali Beach Paradise",
      "description": "Relax on pristine beaches",
      "location": "Bali",
      "country": "Indonesia",
      "category": "beach",
      "duration": 7,
      "season": "summer",
      "budget": "medium",
      "image_url": "https://example.com/bali.jpg",
      "rating": 4.7,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "per_page": 10
}
```

#### GET /api/v1/tours/:id
Get a specific tour by ID.

**Response:** (200 OK)
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Tokyo Cherry Blossom Tour",
    "description": "Experience the beautiful cherry blossoms in Tokyo",
    "location": "Tokyo",
    "country": "Japan",
    "category": "cultural",
    "duration": 5,
    "season": "spring",
    "budget": "high",
    "image_url": "https://example.com/tokyo.jpg",
    "rating": 4.8,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### PUT /api/v1/tours/:id
Update a tour destination.

**Request Body:** (all fields optional)
```json
{
  "name": "Updated Tour Name",
  "rating": 4.9
}
```

**Response:** (200 OK)
```json
{
  "success": true,
  "message": "Tour updated successfully",
  "data": {
    "id": 1,
    "name": "Updated Tour Name",
    "description": "Experience the beautiful cherry blossoms in Tokyo",
    "location": "Tokyo",
    "country": "Japan",
    "category": "cultural",
    "duration": 5,
    "season": "spring",
    "budget": "high",
    "image_url": "https://example.com/tokyo.jpg",
    "rating": 4.9,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  }
}
```

#### DELETE /api/v1/tours/:id
Delete a tour destination.

**Response:** (200 OK)
```json
{
  "success": true,
  "message": "Tour deleted successfully"
}
```

#### POST /api/v1/tours/suggest
Get tour suggestions based on preferences.

**Request Body:**
```json
{
  "category": "beach",
  "budget": "medium",
  "season": "summer",
  "min_rating": 4.5
}
```

**Response:** (200 OK)
```json
{
  "success": true,
  "message": "Tour suggestions generated",
  "data": [
    {
      "id": 2,
      "name": "Bali Beach Paradise",
      "description": "Relax on pristine beaches and explore traditional Balinese culture",
      "location": "Bali",
      "country": "Indonesia",
      "category": "beach",
      "duration": 7,
      "season": "summer",
      "budget": "medium",
      "image_url": "https://example.com/bali.jpg",
      "rating": 4.7,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "success": false,
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Service not ready

## Data Models

### Tour Categories
- `beach`: Beach destinations
- `mountain`: Mountain and hiking destinations
- `cultural`: Cultural and historical sites
- `city`: Urban destinations

### Budget Levels
- `low`: Budget-friendly options
- `medium`: Moderate pricing
- `high`: Premium experiences

### Seasons
- `spring`: March - May
- `summer`: June - August
- `fall`: September - November
- `winter`: December - February
- `all`: Year-round destinations
