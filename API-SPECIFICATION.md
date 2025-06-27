# Workflow Builder Backend API Specification

This document outlines the REST API endpoints that your backend should implement to work with the Workflow Builder frontend.

## Base URL
```
http://localhost:8080/api
```

## Authentication
The frontend is configured to send requests with `Content-Type: application/json`. If you need authentication, you can add headers in the API client.

## Endpoints

### Workflows

#### GET /workflows
Get all workflows
**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "nodes": [],
      "connections": [],
      "isActive": boolean,
      "status": "active" | "draft" | "paused" | "error",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ]
}
```

#### POST /workflows
Create a new workflow
**Request:**
```json
{
  "name": "string",
  "description": "string",
  "nodes": [],
  "connections": []
}
```
**Response:**
```json
{
  "success": true,
  "data": {
    "id": "string",
    "name": "string",
    "description": "string",
    "nodes": [],
    "connections": [],
    "isActive": false,
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
}
```

#### GET /workflows/{id}
Get a specific workflow
**Response:**
```json
{
  "success": true,
  "data": {
    "id": "string",
    "name": "string",
    "description": "string",
    "nodes": [
      {
        "id": "string",
        "type": "schedule|webhook|postgres|mysql|...",
        "triggerType": "schedule|webhook",
        "position": {
          "x": number,
          "y": number
        },
        "data": {},
        "connections": ["string"]
      }
    ],
    "connections": [
      {
        "id": "string",
        "sourceId": "string",
        "targetId": "string"
      }
    ],
    "isActive": boolean,
    "status": "active" | "draft" | "paused" | "error",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
}
```

#### PUT /workflows/{id}
Update a workflow
**Request:**
```json
{
  "name": "string",
  "description": "string",
  "nodes": [],
  "connections": [],
  "updatedAt": "2023-01-01T00:00:00Z"
}
```
**Response:** Same as GET /workflows/{id}

#### DELETE /workflows/{id}
Delete a workflow
**Response:**
```json
{
  "success": true,
  "data": null,
  "message": "Workflow deleted successfully"
}
```

#### POST /workflows/{id}/execute
Execute a workflow
**Request:**
```json
{
  "input": {
    "key": "value"
  }
}
```
**Response:**
```json
{
  "success": true,
  "data": {
    "executionId": "string"
  }
}
```

#### GET /workflows/{workflowId}/executions/{executionId}
Get workflow execution status
**Response:**
```json
{
  "success": true,
  "data": {
    "status": "running|completed|failed|cancelled",
    "progress": number,
    "result": {},
    "error": "string"
  }
}
```

#### PUT /workflows/{id}/status
Toggle workflow active status
**Request:**
```json
{
  "isActive": boolean
}
```
**Response:** Same as GET /workflows/{id}

### Templates

#### GET /workflows/templates
Get workflow templates
**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "category": "string",
      "nodes": [],
      "connections": []
    }
  ]
}
```

#### POST /workflows/from-template
Create workflow from template
**Request:**
```json
{
  "templateId": "string",
  "name": "string"
}
```
**Response:** Same as POST /workflows

### Webhooks

#### POST /webhooks/test
Test webhook endpoint
**Request:**
```json
{
  "url": "string",
  "payload": {}
}
```
**Response:**
```json
{
  "success": true,
  "data": {
    "success": boolean,
    "response": {},
    "error": "string"
  }
}
```

## Error Responses

All error responses should follow this format:
```json
{
  "success": false,
  "error": "Error message",
  "message": "Optional detailed message"
}
```

HTTP Status Codes:
- 200: Success
- 201: Created
- 400: Bad Request
- 404: Not Found
- 500: Internal Server Error

## Node Types

The following node types are supported:
- `schedule` - Scheduled triggers
- `webhook` - Webhook triggers
- `postgres`, `mysql`, `mariadb`, `tidb`, `neon`, `mongodb`, `snowflake`, `supabase`, `clickhouse`, `hydra` - Database nodes
- `rest-api`, `graphql` - API nodes
- `ai-agent`, `transformer`, `condition`, `loop`, `response`, `error-handler` - Action nodes

## Example Implementation Notes

1. Store workflows in your preferred database (PostgreSQL, MongoDB, etc.)
2. Implement proper validation for node types and connections
3. Handle workflow execution asynchronously
4. Implement proper error handling and logging
5. Consider rate limiting for API endpoints
6. Add authentication/authorization as needed
