# Envelope Security Documentation

## Overview

The envelope security feature provides encrypted data transmission for sensitive information like login credentials. This prevents exposure of email addresses, phone numbers, and passwords during API communication.

## How It Works

### 1. Encryption Process
- Client sensitive data is encrypted using AES-256-GCM encryption
- Uses PBKDF2 key derivation for additional security
- Includes timestamp validation to prevent replay attacks
- Adds random nonce for each request

### 2. Envelope Structure
```json
{
  "envelope": {
    "data": "base64-encoded-encrypted-data",
    "nonce": "base64-encoded-random-nonce",
    "timestamp": 1699887654
  }
}
```

## API Endpoints

### Create Envelope
**POST** `/api/v1/auth/create-envelope`

Encrypts any data into a secure envelope format.

**Request:**
```json
{
  "identifier": "user@example.com",
  "password": "secret-password"
}
```

**Response:**
```json
{
  "data": {
    "envelope": {
      "data": "encrypted-base64-string",
      "nonce": "random-nonce",
      "timestamp": 1699887654
    }
  },
  "message": "Envelope created successfully"
}
```

### Secure Login
**POST** `/api/v1/auth/login`

Accepts both direct login and envelope-encrypted login.

**Direct Login (Legacy):**
```json
{
  "identifier": "user@example.com",
  "password": "secret-password"
}
```

**Envelope Login (Secure):**
```json
{
  "envelope": {
    "data": "encrypted-base64-string",
    "nonce": "random-nonce", 
    "timestamp": 1699887654
  }
}
```

## Implementation Examples

### JavaScript Client Example
```javascript
// Create envelope for login
const loginData = {
  identifier: "user@example.com",
  password: "secret-password"
};

// Step 1: Create envelope
const envelopeResponse = await fetch('/api/v1/auth/create-envelope', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(loginData)
});

const { data: { envelope } } = await envelopeResponse.json();

// Step 2: Use envelope for secure login
const loginResponse = await fetch('/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ envelope })
});

const result = await loginResponse.json();
console.log('Login successful:', result);
```

### cURL Examples

**Create Envelope:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/create-envelope \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "user@example.com",
    "password": "secret-password"
  }'
```

**Secure Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "envelope": {
      "data": "encrypted-data-here",
      "nonce": "nonce-here",
      "timestamp": 1699887654
    }
  }'
```

## Security Features

### 1. **AES-256-GCM Encryption**
- Industry standard encryption algorithm
- Provides both confidentiality and authenticity
- Resistant to tampering and forgery

### 2. **Key Derivation (PBKDF2)**
- Derives encryption key from secret using PBKDF2
- 10,000 iterations for added security
- Uses SHA-256 as the hash function

### 3. **Timestamp Validation**
- Prevents replay attacks
- 5-minute expiration window
- Allows 1-minute clock skew tolerance

### 4. **Random Nonces**
- Each envelope gets a unique nonce
- Prevents duplicate request attacks
- 16-byte cryptographically secure random values

### 5. **Backward Compatibility**
- Supports both envelope and direct login methods
- Graceful fallback for legacy clients
- No breaking changes to existing API

## Configuration

### Environment Variables
```env
# Envelope encryption secret (change in production!)
ENVELOPE_SECRET=your-super-secret-envelope-key-change-in-production
```

### Security Best Practices
1. **Change default secrets** in production
2. **Use HTTPS** for all communications
3. **Rotate envelope secrets** regularly
4. **Monitor** for failed decryption attempts
5. **Rate limit** envelope creation endpoints

## Error Handling

### Common Errors
- `Invalid or expired envelope` - Envelope decryption failed or expired
- `Envelope timestamp is in the future` - Clock skew too large
- `Invalid envelope data` - Malformed envelope structure
- `Invalid payload format` - Corrupted decrypted data

### Error Responses
```json
{
  "error": "Invalid or expired envelope"
}
```

## Benefits

### 1. **Data Protection**
- Sensitive credentials never transmitted in plain text
- Email addresses and phone numbers are encrypted
- Passwords are double-protected (envelope + hashing)

### 2. **Replay Protection**
- Timestamp validation prevents old requests from being replayed
- Unique nonces ensure each request is fresh

### 3. **Network Security**
- Even if HTTPS is compromised, data remains encrypted
- Additional layer of security for sensitive operations

### 4. **Compliance**
- Helps meet data protection requirements
- Reduces risk of credential exposure in logs
- Supports privacy regulations (GDPR, CCPA, etc.)

## Migration Guide

### Existing Clients
- No immediate action required
- Direct login method remains supported
- Gradual migration to envelope method recommended

### New Clients
- Should implement envelope method from start
- Use create-envelope endpoint for sensitive data
- Follow security best practices

This envelope security system significantly enhances the protection of sensitive user data during authentication while maintaining backward compatibility with existing clients.