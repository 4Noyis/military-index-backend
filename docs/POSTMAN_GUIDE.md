# Postman Collection Guide

This guide explains how to use the Military Index API Postman collection for testing and development.

## Quick Start

### 1. Import the Collection

1. Open Postman
2. Click **Import** in the top left
3. Select **File** tab
4. Choose `docs/Military_Index_API.postman_collection.json`
5. Click **Import**

### 2. Set Environment Variables

The collection uses a variable for the base URL:

- **Variable**: `base_url`
- **Default Value**: `http://localhost:8080`

To change the environment:

1. Click the **Environment** dropdown (top right)
2. Select **Environments** from sidebar
3. Create a new environment or edit existing
4. Add variable:
   - Key: `base_url`
   - Value: `http://localhost:8080` (or your API URL)
5. Save and activate the environment

### 3. Authentication Setup

The API requires JWT authentication for write operations (POST, PUT, DELETE).

#### Getting a Token

1. Open the **Authentication** > **Login** request
2. Update the username/password in the request body if needed
3. Click **Send**
4. Copy the token from the response

#### Using the Token

**Option A: Manual (for single requests)**
1. Open any protected request (POST, PUT, DELETE)
2. Go to the **Authorization** tab
3. Select **Type**: `Bearer Token`
4. Paste your token in the **Token** field

**Option B: Environment Variable (recommended)**
1. After logging in, save the token to an environment variable:
   - Click **Tests** tab in the Login request
   - Add this script:
   ```javascript
   if (pm.response.code === 200) {
       const response = pm.response.json();
       pm.environment.set("auth_token", response.data.token);
   }
   ```
2. In protected requests, use `{{auth_token}}` in the Authorization header

**Note:** Tokens expire after 24 hours. Re-run the Login request to get a new token.

## Collection Structure

The collection is organized into logical folders:

### 1. Health Check
- **Get API Health**: Verify the API is running

### 2. Authentication
- **Login**: Get JWT token for protected operations

### 3. Countries
- **Get All Countries**: Paginated list of countries
- **Get Country by ID**: Single country by ID
- **Get Country by Code**: Single country by ISO code
- **Search Countries**: Full-text search
- **Create Country**: Add new country (Admin)
- **Update Country**: Modify existing country (Admin)
- **Delete Country**: Remove country (Admin)

### 4. Technologies
- **Get All Technologies**: Paginated list with sorting
- **Get Technology by ID**: Single technology with relations
- **Get Technologies by Country**: Filter by country code
- **Get Technologies by Category**: Filter by category name
- **Get Technologies by Status**: Filter by status
- **Get Technologies by Year Range**: Filter by year range
- **Search Technologies**: Full-text search
- **Create Technology**: Add new technology (Admin, requires auth)
- **Update Technology**: Modify existing technology (Admin, requires auth)
- **Delete Technology**: Remove technology (Admin, requires auth)

### 5. Categories
- **Get All Categories**: List all technology categories
- **Get Category by ID**: Single category by ID
- **Create Category**: Create a new category (requires auth)
- **Update Category**: Update an existing category (requires auth)
- **Delete Category**: Delete a category (requires auth)

### 6. Examples - Common Workflows
Pre-configured requests for common use cases:
- Get Turkish Drones
- Get Recent Technologies (2020-2025)
- Get All Aircraft
- Get Prototype Technologies

## Query Parameters Guide

### Pagination
Most list endpoints support:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)

Example:
```
GET /api/v1/technologies?page=2&limit=20
```

### Sorting (Technologies only)
- `sort` - Sort field (id, name, year_developed, year_deployed)
- `order` - Sort order (asc, desc)

Example:
```
GET /api/v1/technologies?sort=year_developed&order=desc
```

### Search
- `q` - Search query

Example:
```
GET /api/v1/technologies?q=drone
```

### Filtering
Use specific endpoints:
- `/api/v1/technologies/country/{code}` - By country
- `/api/v1/technologies/category/{name}` - By category
- `/api/v1/technologies/status/{status}` - By status
- `/api/v1/technologies/year-range?start=2020&end=2024` - By year

## Testing Workflow

### Basic Testing Flow

1. **Health Check**
   - Run "Get API Health" to verify the API is running
   - Expected: `200 OK` with status "healthy"

2. **Explore Data**
   - Run "Get All Countries" to see available countries
   - Run "Get All Technologies" to see sample technologies

3. **Test Filtering**
   - Run "Get Technologies by Country" (e.g., TUR)
   - Run "Get Technologies by Category" (e.g., Drones)
   - Run "Get Technologies by Status" (e.g., current)

4. **Test Search**
   - Run "Search Technologies" with query "bayraktar"
   - Run "Search Countries" with query "Turkey"

5. **Test CRUD Operations** (Optional)
   - Create a test technology
   - Update the created technology
   - Delete the test technology

### Advanced Testing

#### Test Pagination
1. Get page 1 with limit 2: `?page=1&limit=2`
2. Get page 2 with limit 2: `?page=2&limit=2`
3. Verify total count matches across pages

#### Test Sorting
1. Sort by name ascending: `?sort=name&order=asc`
2. Sort by year descending: `?sort=year_developed&order=desc`
3. Verify results are correctly sorted

#### Test Error Handling
1. Request invalid ID: `/api/v1/technologies/99999`
   - Expected: `404 Not Found`
2. Request invalid country code: `/api/v1/countries/code/INVALID`
   - Expected: `404 Not Found`

## Common Response Formats

### Success Response
```json
{
  "success": true,
  "data": { /* response data */ }
}
```

### Paginated Response
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
```json
{
  "success": false,
  "error": "Error message"
}
```

## Tips & Best Practices

1. **Start Simple**: Begin with GET requests before trying POST/PUT/DELETE
2. **Check Responses**: Review response bodies to understand data structure
3. **Use Variables**: Leverage Postman variables for IDs and repeated values
4. **Save Examples**: Save successful responses as examples in Postman
5. **Test Environments**: Create separate environments for dev/staging/prod
6. **Rate Limiting**: Be aware of the 60 requests/minute rate limit

## Troubleshooting

### Connection Refused
- Ensure Docker services are running: `make up-dev`
- Verify API Gateway is accessible: `curl http://localhost:8080/health`
- Check the correct port (default: 8080)

### 404 Not Found
- Check the endpoint path matches the API documentation
- Verify trailing slashes (both with/without should work)
- Ensure the resource ID exists

### 500 Internal Server Error
- Check API Gateway logs: `docker logs military_index_api_gateway`
- Check service logs: `docker logs military_index_technology_service`
- Verify database is running: `docker logs military_index_db`

### Rate Limit Exceeded (429)
- Wait one minute before making more requests
- Current limit: 60 requests per minute per IP
- Consider adding delays between requests in automated tests

## Environment Setup Examples

### Local Development
```
base_url: http://localhost:8080
```

### Docker Internal Network
```
base_url: http://api-gateway:8080
```

### Staging (Future)
```
base_url: https://api-staging.military-index.com
```

### Production (Future)
```
base_url: https://api.military-index.com
```

## Additional Resources

- [API Documentation](API.md) - Complete endpoint reference
- [Setup Guide](SETUP.md) - Local development setup
- [Contributing Guide](CONTRIBUTING.md) - Contribution guidelines

---

**Need Help?**
- Check the [API Documentation](API.md) for detailed endpoint information
- Review logs: `make logs`
- Report issues: [GitHub Issues](https://github.com/yourusername/military-index-backend/issues)
