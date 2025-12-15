# Military Index Backend - API Documentation

**Base URL:** `http://localhost:8080` (Development)

**API Version:** v1

**Content-Type:** `application/json`

## Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Rate Limiting](#rate-limiting)
- [Response Format](#response-format)
- [Error Codes](#error-codes)
- [Endpoints](#endpoints)
  - [Health Check](#health-check)
  - [Countries API](#countries-api)
  - [Technologies API](#technologies-api)
  - [Categories API](#categories-api)

---

## Overview

The Military Index Backend API provides access to a comprehensive database of military technologies from around the world. The API follows RESTful principles and returns JSON responses.

### Key Features

- **Pagination**: All list endpoints support pagination
- **Search**: Full-text search capabilities
- **Filtering**: Filter by country, category, status, and year range
- **Sorting**: Customizable sort order
- **Rate Limiting**: 60 requests per minute per IP

---

## Authentication

Currently, the API is **publicly accessible** without authentication. Write operations (POST, PUT, DELETE) are available but may require authentication in future versions.

---

## Rate Limiting

The API implements rate limiting to ensure fair usage:

- **Limit**: 60 requests per minute per IP address
- **Rate Limit Headers** (not yet implemented):
  - `X-RateLimit-Limit`: Maximum requests per minute
  - `X-RateLimit-Remaining`: Remaining requests in current window
  - `X-RateLimit-Reset`: Time when the rate limit resets

When rate limit is exceeded, you'll receive a `429 Too Many Requests` response.

---

## Response Format

### Success Response

All successful responses follow this format:

```json
{
  "success": true,
  "data": { /* response data */ }
}
```

### Paginated Response

List endpoints return paginated results:

```json
{
  "success": true,
  "data": [ /* array of items */ ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 42,
    "total_pages": 5
  }
}
```

### Error Response

Error responses include descriptive messages:

```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| `200 OK` | Request successful |
| `201 Created` | Resource created successfully |
| `400 Bad Request` | Invalid request parameters |
| `404 Not Found` | Resource not found |
| `429 Too Many Requests` | Rate limit exceeded |
| `500 Internal Server Error` | Server error occurred |
| `503 Service Unavailable` | Downstream service unavailable |

---

## Endpoints

### Health Check

#### GET /health

Check the API Gateway health status.

**Request:**
```bash
curl http://localhost:8080/health
```

**Response:** `200 OK`
```json
{
  "status": "healthy",
  "service": "api-gateway",
  "version": "1.0.0"
}
```

---

## Countries API

### List All Countries

#### GET /api/v1/countries

Retrieve a paginated list of all countries.

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `limit` | integer | 10 | Items per page |
| `q` | string | - | Search query (name or code) |

**Request:**
```bash
curl "http://localhost:8080/api/v1/countries?page=1&limit=5"
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Turkey",
      "code": "TUR",
      "flag_url": "/flags/turkey.svg",
      "created_at": "2025-11-18T12:51:49.150303Z",
      "updated_at": "2025-11-18T12:51:49.150303Z"
    },
    {
      "id": 2,
      "name": "United States",
      "code": "USA",
      "flag_url": "/flags/usa.svg",
      "created_at": "2025-11-18T12:51:49.150303Z",
      "updated_at": "2025-11-18T12:51:49.150303Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 5,
    "total": 12,
    "total_pages": 3
  }
}
```

---

### Get Country by ID

#### GET /api/v1/countries/{id}

Retrieve a single country by its ID.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Country ID |

**Request:**
```bash
curl http://localhost:8080/api/v1/countries/1
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Turkey",
    "code": "TUR",
    "flag_url": "/flags/turkey.svg",
    "created_at": "2025-11-18T12:51:49.150303Z",
    "updated_at": "2025-11-18T12:51:49.150303Z"
  }
}
```

**Error Response:** `404 Not Found`
```json
{
  "success": false,
  "error": "Country not found"
}
```

---

### Get Country by Code

#### GET /api/v1/countries/code/{code}

Retrieve a country by its ISO country code.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `code` | string | ISO country code (2-3 letters) |

**Request:**
```bash
curl http://localhost:8080/api/v1/countries/code/TUR
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Turkey",
    "code": "TUR",
    "flag_url": "/flags/turkey.svg",
    "created_at": "2025-11-18T12:51:49.150303Z",
    "updated_at": "2025-11-18T12:51:49.150303Z"
  }
}
```

---

### Search Countries

#### GET /api/v1/countries?q={query}

Search countries by name or code.

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | Yes | Search query |
| `page` | integer | No | Page number (default: 1) |
| `limit` | integer | No | Items per page (default: 10) |

**Request:**
```bash
curl "http://localhost:8080/api/v1/countries?q=Turkey"
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Turkey",
      "code": "TUR",
      "flag_url": "/flags/turkey.svg",
      "created_at": "2025-11-18T12:51:49.150303Z",
      "updated_at": "2025-11-18T12:51:49.150303Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### Create Country

#### POST /api/v1/countries

Create a new country. (Admin only - not yet enforced)

**Request Body:**
```json
{
  "name": "South Korea",
  "code": "KOR",
  "flag_url": "/flags/korea.svg"
}
```

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/countries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "South Korea",
    "code": "KOR",
    "flag_url": "/flags/korea.svg"
  }'
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": 13,
    "name": "South Korea",
    "code": "KOR",
    "flag_url": "/flags/korea.svg",
    "created_at": "2025-12-02T14:30:00Z",
    "updated_at": "2025-12-02T14:30:00Z"
  }
}
```

---

### Update Country

#### PUT /api/v1/countries/{id}

Update an existing country. (Admin only - not yet enforced)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Country ID |

**Request Body:**
```json
{
  "name": "Republic of Turkey",
  "code": "TUR",
  "flag_url": "/flags/turkey.svg"
}
```

**Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/countries/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Republic of Turkey",
    "code": "TUR",
    "flag_url": "/flags/turkey.svg"
  }'
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Republic of Turkey",
    "code": "TUR",
    "flag_url": "/flags/turkey.svg",
    "created_at": "2025-11-18T12:51:49.150303Z",
    "updated_at": "2025-12-02T14:35:00Z"
  }
}
```

---

### Delete Country

#### DELETE /api/v1/countries/{id}

Delete a country. (Admin only - not yet enforced)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Country ID |

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/countries/13
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "message": "Country deleted successfully"
  }
}
```

---

## Technologies API

### List All Technologies

#### GET /api/v1/technologies

Retrieve a paginated list of all military technologies.

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `limit` | integer | 10 | Items per page |
| `sort` | string | `id` | Sort field (id, name, year_developed, year_deployed) |
| `order` | string | `asc` | Sort order (asc, desc) |
| `q` | string | - | Search query (name search) |

**Request:**
```bash
curl "http://localhost:8080/api/v1/technologies?page=1&limit=2&sort=name&order=asc"
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 4,
      "country_id": 1,
      "category_id": 3,
      "name": "Altay MBT",
      "description": "Main battle tank with advanced composite armor",
      "designer": "Otokar/FNSS",
      "year_developed": 2013,
      "year_deployed": 2024,
      "manufacturer": "BMC",
      "unit_cost": null,
      "mass": 65000,
      "length": 770,
      "width": 390,
      "height": 250,
      "status": "current",
      "image_url": "/images/altay.jpg",
      "created_at": "2025-11-18T12:51:49.151104Z",
      "updated_at": "2025-11-18T12:51:49.151104Z",
      "country": {
        "id": 1,
        "name": "Turkey",
        "code": "TUR",
        "flag_url": "/flags/turkey.svg",
        "created_at": "2025-11-18T12:51:49.150303Z",
        "updated_at": "2025-11-18T12:51:49.150303Z"
      },
      "category": {
        "id": 3,
        "name": "Ground Vehicles",
        "description": "Tanks, armored vehicles, and artillery",
        "icon_url": "/icons/ground.svg",
        "created_at": "2025-11-18T12:51:49.149239Z"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 2,
    "total": 6,
    "total_pages": 3
  }
}
```

---

### Get Technology by ID

#### GET /api/v1/technologies/{id}

Retrieve a single technology by its ID, including related country and category information.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Technology ID |

**Request:**
```bash
curl http://localhost:8080/api/v1/technologies/1
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "country_id": 1,
    "category_id": 1,
    "name": "KAAN (TF-X)",
    "description": "Fifth-generation stealth multirole fighter aircraft",
    "designer": "TAI",
    "year_developed": 2023,
    "year_deployed": 2028,
    "manufacturer": "Turkish Aerospace Industries",
    "unit_cost": 100000000,
    "mass": 27000,
    "length": 2100,
    "width": 1400,
    "height": 550,
    "status": "prototype",
    "image_url": "/images/kaan.jpg",
    "created_at": "2025-11-18T12:51:49.151104Z",
    "updated_at": "2025-11-18T12:51:49.151104Z",
    "country": {
      "id": 1,
      "name": "Turkey",
      "code": "TUR",
      "flag_url": "/flags/turkey.svg",
      "created_at": "2025-11-18T12:51:49.150303Z",
      "updated_at": "2025-11-18T12:51:49.150303Z"
    },
    "category": {
      "id": 1,
      "name": "Aircraft",
      "description": "Military aircraft including fighters, bombers, and transport",
      "icon_url": "/icons/aircraft.svg",
      "created_at": "2025-11-18T12:51:49.149239Z"
    }
  }
}
```

**Error Response:** `404 Not Found`
```json
{
  "success": false,
  "error": "technology not found"
}
```

---

### Get Technologies by Country

#### GET /api/v1/technologies/country/{code}

Retrieve all technologies from a specific country.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `code` | string | ISO country code (e.g., TUR, USA, RUS) |

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `limit` | integer | 10 | Items per page |

**Request:**
```bash
curl "http://localhost:8080/api/v1/technologies/country/TUR?page=1&limit=10"
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "country_id": 1,
      "category_id": 1,
      "name": "KAAN (TF-X)",
      "description": "Fifth-generation stealth multirole fighter aircraft",
      "designer": "TAI",
      "year_developed": 2023,
      "year_deployed": 2028,
      "manufacturer": "Turkish Aerospace Industries",
      "unit_cost": 100000000,
      "mass": 27000,
      "length": 2100,
      "width": 1400,
      "height": 550,
      "status": "prototype",
      "image_url": "/images/kaan.jpg",
      "created_at": "2025-11-18T12:51:49.151104Z",
      "updated_at": "2025-11-18T12:51:49.151104Z",
      "country": {
        "id": 1,
        "name": "Turkey",
        "code": "TUR",
        "flag_url": "/flags/turkey.svg",
        "created_at": "2025-11-18T12:51:49.150303Z",
        "updated_at": "2025-11-18T12:51:49.150303Z"
      },
      "category": {
        "id": 1,
        "name": "Aircraft",
        "description": "Military aircraft including fighters, bombers, and transport",
        "icon_url": "/icons/aircraft.svg",
        "created_at": "2025-11-18T12:51:49.149239Z"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 6,
    "total_pages": 1
  }
}
```

---

### Get Technologies by Category

#### GET /api/v1/technologies/category/{name}

Retrieve all technologies in a specific category.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `name` | string | Category name (Aircraft, Naval, Ground Vehicles, Missiles, Electronics, Drones, Space & Satellite) |

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `limit` | integer | 10 | Items per page |

**Request:**
```bash
curl "http://localhost:8080/api/v1/technologies/category/Drones"
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "country_id": 1,
      "category_id": 6,
      "name": "Bayraktar TB2",
      "description": "Medium-altitude long-endurance tactical unmanned combat aerial vehicle",
      "designer": "Baykar",
      "year_developed": 2014,
      "year_deployed": 2014,
      "manufacturer": "Baykar Makina",
      "unit_cost": 5000000,
      "mass": 700,
      "length": 640,
      "width": 1200,
      "height": 220,
      "status": "current",
      "image_url": "/images/tb2.jpg",
      "created_at": "2025-11-18T12:51:49.151104Z",
      "updated_at": "2025-11-18T12:51:49.151104Z",
      "country": {
        "id": 1,
        "name": "Turkey",
        "code": "TUR",
        "flag_url": "/flags/turkey.svg",
        "created_at": "2025-11-18T12:51:49.150303Z",
        "updated_at": "2025-11-18T12:51:49.150303Z"
      },
      "category": {
        "id": 6,
        "name": "Drones",
        "description": "Unmanned aerial vehicles and systems",
        "icon_url": "/icons/drones.svg",
        "created_at": "2025-11-18T12:51:49.149239Z"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 2,
    "total_pages": 1
  }
}
```

---

### Get Technologies by Status

#### GET /api/v1/technologies/status/{status}

Retrieve all technologies with a specific status.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `status` | string | Technology status (historical, current, future, concept, prototype) |

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `limit` | integer | 10 | Items per page |

**Request:**
```bash
curl "http://localhost:8080/api/v1/technologies/status/current"
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "country_id": 1,
      "category_id": 6,
      "name": "Bayraktar TB2",
      "description": "Medium-altitude long-endurance tactical unmanned combat aerial vehicle",
      "designer": "Baykar",
      "year_developed": 2014,
      "year_deployed": 2014,
      "manufacturer": "Baykar Makina",
      "unit_cost": 5000000,
      "mass": 700,
      "length": 640,
      "width": 1200,
      "height": 220,
      "status": "current",
      "image_url": "/images/tb2.jpg",
      "created_at": "2025-11-18T12:51:49.151104Z",
      "updated_at": "2025-11-18T12:51:49.151104Z",
      "country": {
        "id": 1,
        "name": "Turkey",
        "code": "TUR",
        "flag_url": "/flags/turkey.svg",
        "created_at": "2025-11-18T12:51:49.150303Z",
        "updated_at": "2025-11-18T12:51:49.150303Z"
      },
      "category": {
        "id": 6,
        "name": "Drones",
        "description": "Unmanned aerial vehicles and systems",
        "icon_url": "/icons/drones.svg",
        "created_at": "2025-11-18T12:51:49.149239Z"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 4,
    "total_pages": 1
  }
}
```

**Valid Status Values:**
- `historical` - Retired or no longer in service
- `current` - Currently in active service
- `future` - Planned for future deployment
- `concept` - Conceptual design phase
- `prototype` - Prototype/testing phase

---

### Get Technologies by Year Range

#### GET /api/v1/technologies/year-range

Retrieve technologies developed or deployed within a specific year range.

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `start` | integer | Yes | Start year |
| `end` | integer | Yes | End year |
| `page` | integer | No | Page number (default: 1) |
| `limit` | integer | No | Items per page (default: 10) |

**Request:**
```bash
curl "http://localhost:8080/api/v1/technologies/year-range?start=2014&end=2024"
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "country_id": 1,
      "category_id": 6,
      "name": "Bayraktar TB2",
      "description": "Medium-altitude long-endurance tactical unmanned combat aerial vehicle",
      "designer": "Baykar",
      "year_developed": 2014,
      "year_deployed": 2014,
      "manufacturer": "Baykar Makina",
      "unit_cost": 5000000,
      "mass": 700,
      "length": 640,
      "width": 1200,
      "height": 220,
      "status": "current",
      "image_url": "/images/tb2.jpg",
      "created_at": "2025-11-18T12:51:49.151104Z",
      "updated_at": "2025-11-18T12:51:49.151104Z",
      "country": {
        "id": 1,
        "name": "Turkey",
        "code": "TUR",
        "flag_url": "/flags/turkey.svg",
        "created_at": "2025-11-18T12:51:49.150303Z",
        "updated_at": "2025-11-18T12:51:49.150303Z"
      },
      "category": {
        "id": 6,
        "name": "Drones",
        "description": "Unmanned aerial vehicles and systems",
        "icon_url": "/icons/drones.svg",
        "created_at": "2025-11-18T12:51:49.149239Z"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 5,
    "total_pages": 1
  }
}
```

---

### Search Technologies

#### GET /api/v1/technologies?q={query}

Search technologies by name using full-text search.

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | Yes | Search query |
| `page` | integer | No | Page number (default: 1) |
| `limit` | integer | No | Items per page (default: 10) |

**Request:**
```bash
curl "http://localhost:8080/api/v1/technologies?q=bayraktar"
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "country_id": 1,
      "category_id": 6,
      "name": "Bayraktar TB2",
      "description": "Medium-altitude long-endurance tactical unmanned combat aerial vehicle",
      "designer": "Baykar",
      "year_developed": 2014,
      "year_deployed": 2014,
      "manufacturer": "Baykar Makina",
      "unit_cost": 5000000,
      "mass": 700,
      "length": 640,
      "width": 1200,
      "height": 220,
      "status": "current",
      "image_url": "/images/tb2.jpg",
      "created_at": "2025-11-18T12:51:49.151104Z",
      "updated_at": "2025-11-18T12:51:49.151104Z",
      "country": {
        "id": 1,
        "name": "Turkey",
        "code": "TUR",
        "flag_url": "/flags/turkey.svg",
        "created_at": "2025-11-18T12:51:49.150303Z",
        "updated_at": "2025-11-18T12:51:49.150303Z"
      },
      "category": {
        "id": 6,
        "name": "Drones",
        "description": "Unmanned aerial vehicles and systems",
        "icon_url": "/icons/drones.svg",
        "created_at": "2025-11-18T12:51:49.149239Z"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 2,
    "total_pages": 1
  }
}
```

---

### Create Technology

#### POST /api/v1/technologies

Create a new military technology. (Admin only - not yet enforced)

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `country_id` | integer | Yes | Country ID |
| `category_id` | integer | Yes | Category ID |
| `name` | string | Yes | Technology name |
| `description` | string | No | Detailed description |
| `designer` | string | No | Design company/organization |
| `year_developed` | integer | No | Year of development |
| `year_deployed` | integer | No | Year of deployment |
| `manufacturer` | string | No | Manufacturing company |
| `unit_cost` | integer | No | Unit cost in USD |
| `mass` | integer | No | Mass in kg |
| `length` | integer | No | Length in cm |
| `width` | integer | No | Width in cm |
| `height` | integer | No | Height in cm |
| `status` | string | Yes | Status (historical/current/future/concept/prototype) |
| `image_url` | string | No | Image URL |

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/technologies \
  -H "Content-Type: application/json" \
  -d '{
    "country_id": 1,
    "category_id": 6,
    "name": "Bayraktar K1z1lelma",
    "description": "Unmanned fighter jet with stealth capabilities",
    "designer": "Baykar",
    "year_developed": 2023,
    "year_deployed": 2025,
    "manufacturer": "Baykar Makina",
    "unit_cost": 20000000,
    "mass": 6000,
    "length": 1450,
    "width": 1400,
    "height": 380,
    "status": "prototype",
    "image_url": "/images/kizilelma.jpg"
  }'
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": 7,
    "country_id": 1,
    "category_id": 6,
    "name": "Bayraktar K1z1lelma",
    "description": "Unmanned fighter jet with stealth capabilities",
    "designer": "Baykar",
    "year_developed": 2023,
    "year_deployed": 2025,
    "manufacturer": "Baykar Makina",
    "unit_cost": 20000000,
    "mass": 6000,
    "length": 1450,
    "width": 1400,
    "height": 380,
    "status": "prototype",
    "image_url": "/images/kizilelma.jpg",
    "created_at": "2025-12-02T14:45:00Z",
    "updated_at": "2025-12-02T14:45:00Z",
    "country": {
      "id": 1,
      "name": "Turkey",
      "code": "TUR",
      "flag_url": "/flags/turkey.svg",
      "created_at": "2025-11-18T12:51:49.150303Z",
      "updated_at": "2025-11-18T12:51:49.150303Z"
    },
    "category": {
      "id": 6,
      "name": "Drones",
      "description": "Unmanned aerial vehicles and systems",
      "icon_url": "/icons/drones.svg",
      "created_at": "2025-11-18T12:51:49.149239Z"
    }
  }
}
```

---

### Update Technology

#### PUT /api/v1/technologies/{id}

Update an existing technology. (Admin only - not yet enforced)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Technology ID |

**Request Body:** Same fields as Create Technology

**Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/technologies/1 \
  -H "Content-Type: application/json" \
  -d '{
    "country_id": 1,
    "category_id": 1,
    "name": "KAAN (TF-X)",
    "description": "Fifth-generation stealth multirole fighter aircraft with advanced avionics",
    "designer": "TAI",
    "year_developed": 2023,
    "year_deployed": 2029,
    "manufacturer": "Turkish Aerospace Industries",
    "unit_cost": 120000000,
    "mass": 27000,
    "length": 2100,
    "width": 1400,
    "height": 550,
    "status": "prototype",
    "image_url": "/images/kaan.jpg"
  }'
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "country_id": 1,
    "category_id": 1,
    "name": "KAAN (TF-X)",
    "description": "Fifth-generation stealth multirole fighter aircraft with advanced avionics",
    "designer": "TAI",
    "year_developed": 2023,
    "year_deployed": 2029,
    "manufacturer": "Turkish Aerospace Industries",
    "unit_cost": 120000000,
    "mass": 27000,
    "length": 2100,
    "width": 1400,
    "height": 550,
    "status": "prototype",
    "image_url": "/images/kaan.jpg",
    "created_at": "2025-11-18T12:51:49.151104Z",
    "updated_at": "2025-12-02T14:50:00Z",
    "country": {
      "id": 1,
      "name": "Turkey",
      "code": "TUR",
      "flag_url": "/flags/turkey.svg",
      "created_at": "2025-11-18T12:51:49.150303Z",
      "updated_at": "2025-11-18T12:51:49.150303Z"
    },
    "category": {
      "id": 1,
      "name": "Aircraft",
      "description": "Military aircraft including fighters, bombers, and transport",
      "icon_url": "/icons/aircraft.svg",
      "created_at": "2025-11-18T12:51:49.149239Z"
    }
  }
}
```

---

### Delete Technology

#### DELETE /api/v1/technologies/{id}

Delete a technology. (Admin only - not yet enforced)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Technology ID |

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/technologies/7
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "message": "Technology deleted successfully"
  }
}
```

---

## Categories API

### List All Categories

#### GET /api/v1/categories

Retrieve all military technology categories.

**Request:**
```bash
curl http://localhost:8080/api/v1/categories
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Aircraft",
      "description": "Military aircraft including fighters, bombers, and transport",
      "icon_url": "https://api.iconify.design/mdi/airplane.svg",
      "created_at": "2025-11-18T12:51:49.149239Z"
    },
    {
      "id": 2,
      "name": "Naval",
      "description": "Naval vessels and submarine technology",
      "icon_url": "https://api.iconify.design/mdi/ferry.svg",
      "created_at": "2025-11-18T12:51:49.149239Z"
    },
    {
      "id": 3,
      "name": "Ground Vehicles",
      "description": "Tanks, armored vehicles, and artillery",
      "icon_url": "https://api.iconify.design/mdi/tank.svg",
      "created_at": "2025-11-18T12:51:49.149239Z"
    },
    {
      "id": 4,
      "name": "Missiles",
      "description": "Missile systems and rockets",
      "icon_url": "https://api.iconify.design/mdi/rocket-launch.svg",
      "created_at": "2025-11-18T12:51:49.149239Z"
    },
    {
      "id": 5,
      "name": "Electronics",
      "description": "Radar, communication, and electronic warfare",
      "icon_url": "https://api.iconify.design/mdi/radar.svg",
      "created_at": "2025-11-18T12:51:49.149239Z"
    },
    {
      "id": 6,
      "name": "Drones",
      "description": "Unmanned aerial vehicles and systems",
      "icon_url": "https://api.iconify.design/mdi/quadcopter.svg",
      "created_at": "2025-11-18T12:51:49.149239Z"
    },
    {
      "id": 7,
      "name": "Space & Satellite",
      "description": "Space-based military technology",
      "icon_url": "https://api.iconify.design/mdi/satellite-variant.svg",
      "created_at": "2025-11-18T12:51:49.149239Z"
    }
  ]
}
```

---

### Get Category by ID

#### GET /api/v1/categories/{id}

Retrieve a single category by its ID.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Category ID |

**Request:**
```bash
curl http://localhost:8080/api/v1/categories/1
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Aircraft",
    "description": "Military aircraft including fighters, bombers, and transport",
    "icon_url": "https://api.iconify.design/mdi/airplane.svg",
    "created_at": "2025-11-18T12:51:49.149239Z"
  }
}
```

**Error Response:** `404 Not Found`
```json
{
  "success": false,
  "error": "Category not found"
}
```

---

### Create Category

#### POST /api/v1/categories

Create a new technology category.

**Request Body:**
```json
{
  "name": "Cyber Warfare",
  "description": "Cyber security and offensive cyber capabilities",
  "icon_url": "https://api.iconify.design/mdi/shield-lock.svg"
}
```

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cyber Warfare",
    "description": "Cyber security and offensive cyber capabilities",
    "icon_url": "https://api.iconify.design/mdi/shield-lock.svg"
  }'
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": 8,
    "name": "Cyber Warfare",
    "description": "Cyber security and offensive cyber capabilities",
    "icon_url": "https://api.iconify.design/mdi/shield-lock.svg",
    "created_at": "2025-12-14T20:30:00Z",
    "updated_at": "2025-12-14T20:30:00Z"
  },
  "message": "Category created successfully"
}
```

**Error Response:** `400 Bad Request`
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Failed to create category",
    "details": "category name is required"
  }
}
```

---

### Update Category

#### PUT /api/v1/categories/{id}

Update an existing category.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Category ID |

**Request Body:**
```json
{
  "name": "Cyber Warfare",
  "description": "Cyber security, offensive and defensive cyber capabilities",
  "icon_url": "https://api.iconify.design/mdi/shield-lock.svg"
}
```

**Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/categories/8 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cyber Warfare",
    "description": "Cyber security, offensive and defensive cyber capabilities",
    "icon_url": "https://api.iconify.design/mdi/shield-lock.svg"
  }'
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 8,
    "name": "Cyber Warfare",
    "description": "Cyber security, offensive and defensive cyber capabilities",
    "icon_url": "https://api.iconify.design/mdi/shield-lock.svg",
    "created_at": "2025-12-14T20:30:00Z",
    "updated_at": "2025-12-14T20:35:00Z"
  },
  "message": "Category updated successfully"
}
```

**Error Response:** `404 Not Found`
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Category not found"
  }
}
```

---

### Delete Category

#### DELETE /api/v1/categories/{id}

Delete a category by ID.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Category ID |

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/categories/8
```

**Response:** `204 No Content`

**Error Response:** `404 Not Found`
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Category not found"
  }
}
```

**Note:** Categories that are referenced by technologies cannot be deleted due to foreign key constraints.

---

## Data Models

### Country Model

```typescript
{
  id: number;                    // Unique identifier
  name: string;                  // Country name
  code: string;                  // ISO country code (2-3 letters)
  flag_url: string;              // Flag image URL
  created_at: string;            // ISO 8601 timestamp
  updated_at: string;            // ISO 8601 timestamp
}
```

### Technology Model

```typescript
{
  id: number;                    // Unique identifier
  country_id: number;            // Foreign key to Country
  category_id: number;           // Foreign key to TechCategory
  name: string;                  // Technology name
  description: string;           // Detailed description
  designer: string;              // Design company/organization
  year_developed: number;        // Year of development
  year_deployed: number;         // Year of deployment
  manufacturer: string;          // Manufacturing company
  unit_cost: number | null;      // Unit cost in USD
  mass: number;                  // Mass in kg
  length: number;                // Length in cm
  width: number;                 // Width in cm
  height: number;                // Height in cm
  status: string;                // Status (historical/current/future/concept/prototype)
  image_url: string;             // Image URL
  created_at: string;            // ISO 8601 timestamp
  updated_at: string;            // ISO 8601 timestamp
  country: Country;              // Related country object
  category: TechCategory;        // Related category object
}
```

### TechCategory Model

```typescript
{
  id: number;                    // Unique identifier
  name: string;                  // Category name
  description: string;           // Category description
  icon_url: string;              // Icon image URL
  created_at: string;            // ISO 8601 timestamp
  updated_at: string;            // ISO 8601 timestamp
}
```

### Pagination Model

```typescript
{
  page: number;                  // Current page number
  limit: number;                 // Items per page
  total: number;                 // Total number of items
  total_pages: number;           // Total number of pages
}
```

---

## Categories Reference

Available technology categories:

| ID | Name | Description |
|----|------|-------------|
| 1 | Aircraft | Military aircraft including fighters, bombers, and transport |
| 2 | Naval | Naval vessels and submarine technology |
| 3 | Ground Vehicles | Tanks, armored vehicles, and artillery |
| 4 | Missiles | Missile systems and rocket technology |
| 5 | Electronics | Electronic warfare and radar systems |
| 6 | Drones | Unmanned aerial vehicles and systems |
| 7 | Space & Satellite | Space and satellite technology |

---

## Examples

### Complete Workflow Example

Here's a complete example of using the API to explore military technologies:

```bash
# 1. Check API health
curl http://localhost:8080/health

# 2. Get all countries
curl "http://localhost:8080/api/v1/countries"

# 3. Get all categories
curl "http://localhost:8080/api/v1/categories"

# 4. Find Turkey's technologies
curl "http://localhost:8080/api/v1/technologies/country/TUR"

# 5. Get all drones
curl "http://localhost:8080/api/v1/technologies/category/Drones"

# 6. Search for specific technology
curl "http://localhost:8080/api/v1/technologies?q=KAAN"

# 7. Get technologies from 2020-2024
curl "http://localhost:8080/api/v1/technologies/year-range?start=2020&end=2024"

# 8. Get current operational technologies
curl "http://localhost:8080/api/v1/technologies/status/current"
```

---

## Best Practices

1. **Pagination**: Always use pagination for list endpoints to avoid large response sizes
2. **Caching**: Response data changes infrequently - consider caching responses
3. **Rate Limiting**: Respect the rate limit of 60 requests/minute
4. **Error Handling**: Always check the `success` field in responses
5. **Search**: Use the `q` parameter for efficient searching instead of fetching all data
6. **Filtering**: Combine filters (country, category, status) to narrow down results

---

## Support & Feedback

For issues, questions, or feature requests:
- GitHub Issues: [https://github.com/yourusername/military-index-backend/issues](https://github.com/yourusername/military-index-backend/issues)
- Documentation: [docs/](../docs/)

---

**Last Updated:** December 2, 2025
**API Version:** 1.0.0
