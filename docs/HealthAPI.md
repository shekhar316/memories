# Health Check API

Simple health check endpoint for **memories-backend**.


### Usage

```
curl -X GET https://<your-domain>/api/v1/health
```

### Response

```json
{
  "success": true,
  "message": "Application is running fine.",
  "data": {
    "status": "healthy",
    "timestamp": "2025-08-24T00:17:35.920656944+05:30",
    "service": "memories-backend",
    "version": "1.0.0",
    "uptime": "41.561189271s"
  }
}
```
### Note
- No authentication required

