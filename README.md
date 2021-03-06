[![Build Status](https://travis-ci.com/MathieuB1/Golang-User-Management.svg?branch=main)](https://travis-ci.com/MathieuB1/Golang-User-Management)

Welcome to Golang-User-Management a Simple GO API for User Management

# Golang-User-Management

    - Home Page
    - Create User
    - Update User (Basic Auth)
    - Get User (Basic Auth)
    - Delete User (Basic Auth)
    - List Users

## Installation & Boot
1. ```apt update && apt install -y docker.io docker-compose```
2. ```git clone https://github.com/MathieuB1/Golang-User-Management && cd Golang-User-Management/```
3. ```docker-compose down && docker volume prune -f && docker-compose build && docker-compose up```
4. ```docker ps```

You should be able to target http://localhost:8000/ url.

## API EndPoints
```
GET STATUS : http://localhost:8000/ (not protected)
GET ALL USERS: http://localhost:8000/users/ (not protected)

CREATE USER: http://localhost:8000/users/ (Basic Auth)
GET USER: http://localhost:8000/users/ID (Basic Auth)
UPDATE USER: http://localhost:8000/users/ID (Basic Auth)
DELETE USER: http://localhost:8000/users/ID (Basic Auth)
```

## API Response
```
    {
        "ID": "1",
        "FirstName": "Steve",
        "LastName": "Rogers"
        "Email": "s.rogers@avenge.rs"
    }
```

## Documentation
https://documenter.getpostman.com/view/3768217/TzeTK9rk
