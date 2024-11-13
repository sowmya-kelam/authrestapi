## Name
Auth Rest API service

## Description
This Service Below API for user Signup, Signin , Authenticate access token, Revoke access token and Refresh token

1. Sign up (creation of user) using email and password  - /signup   
    Sample Request 

   {   
    "email" :"test2904@gmail.com",
    "password": "test"
    }

2. Sign in - /login
     Authentication of user credentials
     JWT tokens (access token, refresh token) are returned in response 
     access token expire after 2 hrs , refresh after 7 days

    Sample Request 

   {   
    "email" :"test2904@gmail.com",
    "password": "test"
    }

3. Authorization of token  - /authorize-token 
     - Mechanism of sending token along with a request from client to service
     - checks for valid signature and expiry 
    sample request payload :
      {
    "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyOTA0QGdtYWlsLmNvbSIsImV4cCI6MTczMTQxOTM4OSwidXNlcklkIjoxfQ.D9Yk5JhnupeeiH52iRv2dyCsZ9mdr1O4nhHYT_yTAbQ"
     }
4. Revocation of token - /revoke-token 
      revoking a token 
      sample request payload:
      {
    "token" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyOTA0QGdtYWlsLmNvbSIsImV4cCI6MTczMTQxOTM4OSwidXNlcklkIjoxfQ.D9Yk5JhnupeeiH52iRv2dyCsZ9mdr1O4nhHYT_yTAbQ"
    }
5. Refresh a token  -  /refresh-token
    Renew Access token using refresh token
    sample request payload
    {
    "refreshtoken" : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyOTA0QGdtYWlsLmNvbSIsImV4cCI6MTczMjAxNzQzMCwidXNlcklkIjoxfQ.aW0BM82DwcVx2Qht3BT9iuhyqqyLbqCRdtOhvN-ot0k"
}


##RUN the Application 
go run cmd/main.go 

## DB 
Used file store DB "api.db" (inside cmd folder )

## tests
go test -v ./...


## Swagger Documentation

once the server is up , documentation is available at  http://localhost:8080/swagger/index.html#/
