# YouTube Focus

YouTube focus is a dockerized service to continuously collect, retrieve and search YouTube videos.
It uses the YouTube Data API v3 to query and store videos in a PostgreSQL database, and then offers this data up
through a REST API.

## Get it running

### Prerequisites

- Docker
- Git

### Steps

1. Clone the repository
2. Fill in the `.env` file with your YouTube API key. A `.env.sample` file is included for reference with relevant
   environment variables.
3. Spin up the `docker-compose` stack with `docker-compose up`
4. Wait for the database to be initialized, and the API to get populated.
5. The API is now available at `localhost:8080`
6. Query the API on `http://localhost:8080/videos` or `http://localhost:8080/videos?search=your+search+query`

## Features
- [x] Polls the YouTube API in background to retrieve new videos.
- [x] REST API to query video data -- cycle through consistently with
      publish-time based pagination.
- [x] Multiple API keys, cycled through to avoid rate limits.
- [x] One-step setup with Docker.
- [ ] Advanced Search
