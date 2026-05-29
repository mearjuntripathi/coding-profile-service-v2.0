# Coding Profile Service

**Coding Profile Service** is a high-performance Go-based backend API for fetching and aggregating coding profiles across multiple competitive programming platforms — **LeetCode, HackerRank, CodeChef, GeeksforGeeks, and Codeforces**.

The service uses **HTML scraping** and platform APIs to fetch public statistics, badges, certifications, and coding achievements. It features **Redis caching**, **parallel scraping via goroutines**, and is fully **Dockerized** for easy local development and deployment.

> 🚀 **Live API:** [https://coding-profile-service.onrender.com](https://coding-profile-service.onrender.com)

---

## What's New

- ✔ **Redis Caching** — repeat requests served in ~2ms instead of ~135ms
- ✔ **Parallel Scraping** — all platforms fetched concurrently via goroutines
- ✔ **Docker Compose** — run Go app + Redis together in one command
- ✔ **Codeforces** support added
- ✔ **Upstash Redis** support for production (Render)
- ✔ **Per-platform TTLs** — smart cache expiry based on how often stats change

---

## Features

- Fetch coding profile stats for 5 platforms simultaneously:
  - **LeetCode** — Problems Solved (Easy/Medium/Hard)
  - **HackerRank** — Coding Score, Badges, Certifications
  - **CodeChef** — Rating, Max Rating, Global/Country Rank, Contests
  - **GeeksforGeeks** — Problems Solved, Streak, Difficulty breakdown
  - **Codeforces** — Rating, Max Rating, Contests, Problems by type
- **Redis caching** with per-platform TTLs to reduce scraping load
- **Parallel requests** — all platforms fetched at the same time
- Fully modular scraper architecture — easy to add new platforms
- REST API with multi-platform query support
- Graceful fallback — app works even if Redis is unavailable

---

## Tech Stack

| Component | Technology |
|---|---|
| Backend | Go (Golang) 1.24+ |
| Scraping | `goquery`, `net/http` |
| Caching | Redis (Docker locally, Upstash in production) |
| Containerization | Docker + Docker Compose |
| Deployment | Render |
| Data Models | Go structs (`pkg/model/StatsResponse`) |

---

## Project Structure

```bash
coding-profile-service/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point — inits Redis, starts server
├── internal/
│   ├── scraper/
│   │   ├── hackerrank.go
│   │   ├── hackerrankHTMLScraper.go
│   │   ├── codechef.go
│   │   ├── codechefHTMLScraper.go
│   │   ├── leetcode.go
│   │   ├── gfg.go
│   │   └── codeforces.go
│   ├── cache/
│   │   └── cache.go                 # Redis client — Get/Set with TTL
│   └── handler/
│       ├── stats_handler.go         # API handler — parallel fetch + cache logic
│       └── request_handler.go       # HTML landing page
├── pkg/
│   └── model/
│       └── stats_response.go        # Shared data model
├── go.mod
├── go.sum
├── dockerfile                       # Multi-stage Go build
├── docker-compose.yml               # App + Redis together
├── .env                             # Local env vars (not committed)
├── .dockerignore
└── README.md
```

---

## How It Works

```
Incoming Request
      ↓
  StatsHandler
      ↓
  ┌─────────────────────────────────┐
  │  For each platform (goroutine)  │
  │                                 │
  │  1. Check Redis cache           │
  │     ↓ HIT → return (~2ms)      │
  │     ↓ MISS                      │
  │  2. Scrape live platform        │
  │     ↓ (~40-135ms)              │
  │  3. Store in Redis with TTL     │
  └─────────────────────────────────┘
      ↓
  Aggregate all results
      ↓
  JSON Response
```

### Cache TTLs per platform

| Platform | TTL | Reason |
|---|---|---|
| LeetCode | 30 min | Problem counts change occasionally |
| GeeksforGeeks | 30 min | Streak updates daily |
| Codeforces | 1 hour | Contest ratings update after contests |
| CodeChef | 2 hours | Rating changes after contests |
| HackerRank | 6 hours | Badges/certs rarely change |

---

## API Reference

### Endpoint

```
GET /stats
```

### Query Parameters

| Parameter | Description | Example |
|---|---|---|
| `leetcode` | LeetCode username | `mearjuntripathi` |
| `codechef` | CodeChef username | `isthisarjun` |
| `gfg` | GeeksforGeeks username | `mearjuntripathi` |
| `hackerrank` | HackerRank username | `mearjuntripathi` |
| `codeforces` | Codeforces username | `isthisarjun` |

All parameters are optional — pass only the platforms you need.

### Example Request

```
GET https://coding-profile-service.onrender.com/stats?leetcode=mearjuntripathi&codechef=isthisarjun&gfg=mearjuntripathi&hackerrank=mearjuntripathi&codeforces=isthisarjun
```

### Example Response

```json
{
  "profiles": [
    {
      "platform": "leetcode",
      "username": "mearjuntripathi",
      "totalSolved": 710,
      "easySolved": 253,
      "mediumSolved": 427,
      "hardSolved": 30,
      "cached": true
    },
    {
      "platform": "gfg",
      "username": "mearjuntripathi",
      "totalSolved": 534,
      "streak": 50,
      "easySolved": 224,
      "mediumSolved": 277,
      "hardSolved": 33,
      "maxRating": 1705
    },
    {
      "platform": "codechef",
      "username": "isthisarjun",
      "totalSolved": 488,
      "rating": 1593,
      "contestsParticipated": 32,
      "maxRating": 1624,
      "globalRank": 18506,
      "countryRank": 16669
    },
    {
      "platform": "hackerrank",
      "username": "mearjuntripathi",
      "totalSolved": 756,
      "badges": ["Problem Solving", "CPP", "Java", "Python", "SQL", "C language"],
      "certifications": 4,
      "certificationLinks": [
        "https://www.hackerrank.com/certificates/dd88f94012d9",
        "https://www.hackerrank.com/certificates/ad7f9b3ad2e1"
      ]
    },
    {
      "platform": "codeforces",
      "username": "isthisarjun",
      "totalSolved": 7,
      "rating": 860,
      "contestsParticipated": 3,
      "maxRating": 860,
      "questionsByType": { "easy": 7 }
    }
  ]
}
```

> `"cached": true` means the response was served from Redis, not scraped live.

---

## Getting Started

### Option 1 — Run with Docker Compose (Recommended)

Runs Go app + Redis together. No local Go or Redis install needed.

```bash
# Clone the repo
git clone https://github.com/mearjuntripathi/coding-profile-service.git
cd coding-profile-service

# Start everything
sudo docker compose up --build
```

App runs at `http://localhost:8080`

```bash
# Stop everything
sudo docker compose down
```

---

### Option 2 — Run Locally (Go + Docker Redis)

```bash
# Step 1 — Start Redis in Docker
sudo docker run -d --name redis-local -p 6379:6379 redis:7-alpine

# Step 2 — Create .env file
echo "REDIS_ADDR=localhost:6379" > .env
echo "REDIS_PASSWORD=" >> .env

# Step 3 — Install dependencies
go mod tidy

# Step 4 — Run the server
go run cmd/server/main.go
```

App runs at `http://localhost:8080`

---

## Environment Variables

| Variable | Description | Local | Production |
|---|---|---|---|
| `REDIS_ADDR` | Redis address | `localhost:6379` | Upstash endpoint |
| `REDIS_PASSWORD` | Redis password | _(empty)_ | Upstash token |
| `SERVER_PORT` | Server port | `8080` | `8080` |

Create a `.env` file for local development (never commit this):

```bash
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
SERVER_PORT=8080
```

---

## Deployment (Render + Upstash)

### 1. Create free Redis on Upstash

1. Go to [upstash.com](https://upstash.com) → Create free database
2. Copy **Endpoint** and **Password**

### 2. Set environment variables on Render

In your Render dashboard → Environment tab:

| Key | Value |
|---|---|
| `REDIS_ADDR` | `your-db.upstash.io:6379` |
| `REDIS_PASSWORD` | `your-upstash-token` |

Render redeploys automatically. No code changes needed between local and production.

---

## Performance

| Scenario | Response Time |
|---|---|
| Cold request (cache miss, 1 platform) | ~135ms |
| Cold request (cache miss, 5 platforms parallel) | ~135ms |
| Warm request (cache hit) | ~2ms |

---

## Contributing

1. Fork the repository
2. Create a new branch: `git checkout -b feature/new-scraper`
3. Add your scraper in `internal/scraper/`
4. Register it in `stats_handler.go` → `fetchPlatformStats()`
5. Commit: `git commit -am 'Add new scraper'`
6. Push and open a Pull Request

---

## Future Plans

- Add **AtCoder** and **SPOJ** platforms
- **Redis persistence** with volume mounts in Docker
- **Unit tests** for all scrapers
- **Rate limiting** to protect against abuse
- **Webhook support** — notify when rating changes

---

## License

This project is open-source under the **MIT License**.