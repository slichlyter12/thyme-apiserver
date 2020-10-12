# Thyme API Server

This is the API server for [thyme](https://github.com/slichlyter12/thyme) and [thyme-iOS](https://github.com/slichlyter12/thyme-iOS).

## Installation

Requirements:

- Docker
- Docker Compose

Start service: `docker-compose up`

(Add `--build` flag when making changes to rebuild the Docker container)

This will start a local instance of AWS DynamoDB and the API server.

## API

The API server runs at `:8080` and the DynamoDB backend runs at `:8000`

### `/init` (POST)

Creates the 'Recipe' table

### `/table` (GET)

Returns a list of all tables in the DynamoDB

### `/recipe`

- (POST) Creates a recipe
- (GET) Lists all recipes

#### Input

- (POST) Requires a JSON Body with valid Recipe types (see data structure below)
- (GET) N/A

#### Output

- (POST) Success message
- (GET) List of recipes

## Data Structure

Recipe:

```golang
type Recipe struct {
    ID          string
    Name        string
    Author      string
    Description string
    Cuisine     string
    imageName   string
    incredients map[string]string
    steps       []string
}
```
