# Thyme API Server

This is the API server for [thyme](https://github.com/slichlyter12/thyme) and [thyme-iOS](https://github.com/slichlyter12/thyme-iOS).

## Installation

Requirements: 
* Docker
* Docker Compose

Start service: `docker-compose up`

(Add `--build` flag when making changes to rebuild the Docker container)

This will start a local instance of AWS DynamoDB and the API server.

## API
The API server runs at `:8080` and the DynamoDB backend runs at `:8000`

### `/init` (POST)
Creates the 'Recipe' table

### `/table` (GET)
Returns a list of all tables in the DynamoDB

### `/recipe` (POST)
Creates a recipe

#### Input
Requires a JSON Body with valid Recipe types (see data structure below)

#### Output
Success message

## Data Structure
`Recipe`
```
    {
        "name":     string,
        "author":   string
    }
```
