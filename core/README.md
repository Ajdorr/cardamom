
For local environment variables required see config/init.go

# Install
```bash
go run bin/migrate/migrate.go
```

# Run Locally
```bash
go run bin/server/server.go
```

# OAuth2 setup
## Local environment variables
```bash
OAUTH_GITHUB_CLIENT_ID="41ff89a3f1011fec0fae"
OAUTH_GITHUB_CLIENT_SECRET="a8f58276de43ad1fc69152a3a96fd21a4ba6f660"
OAUTH_FACEBOOK_CLIENT_ID="643465770672262"
OAUTH_FACEBOOK_CLIENT_SECRET="eae6fd76403ea87c39f75d6de12279fe"
OAUTH_MICROSOFT_CLIENT_ID="478c5aab-132d-4da2-9781-22286b5a2bdc"
OAUTH_MICROSOFT_CLIENT_SECRET="y1u8Q~hwyOOAMaIKdBLzCg8QrCL52~lwr0uM5cem"
OAUTH_GOOGLE_CREDS_FILEPATH="/path/to/credentials/file/google-creds.json"
```

### Google credentials file
```json
{
  "web": {
    "client_id": "357567058246-dtms8mvh75urpg6hd1c45sdv2c2k9a4a.apps.googleusercontent.com",
    "project_id": "jaguar-cardamom-develop",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_secret": "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",
    "redirect_uris": [
      "http://localhost:8080/auth/oauth-return/google"
    ]
  }
}
```