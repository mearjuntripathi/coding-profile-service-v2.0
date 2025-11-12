# Coding Profile Service

**Coding Profile Service** is a Go-based backend API for fetching and aggregating coding profiles across multiple competitive programming platforms like **HackerRank, CodeChef, LeetCode, and GeeksforGeeks**.

The service uses **HTML scraping** and APIs where available to fetch public statistics, badges, certifications, and other coding achievements. It is designed with **GraphQL** for fast and flexible queries, making it faster than typical JavaScript or Python implementations.

---

## Features

* Fetch coding profile stats for multiple platforms:

  * **HackerRank:** Coding Score, Problems Solved, Badges, Certifications.
  * **CodeChef:** Rating, Max Rating, Contests Participated, Problems Solved, Rankings.
  * **LeetCode:** Problems Solved, Acceptance Rate, Contest Ratings (future-ready).
  * **GeeksforGeeks:** Problems Solved, Certificates, Badges (future-ready).
* Fully modular scraping architecture:

  * Each platform has its own scraper.
  * Easy to extend for additional platforms.
* GraphQL API for flexible and fast queries.
* Minimal dependencies, written entirely in **Go** for performance.

---

## Tech Stack

| Component   | Technology                                                        |
| ----------- | ----------------------------------------------------------------- |
| Backend     | Go (Golang)                                                       |
| GraphQL API | [gqlgen](https://gqlgen.com/) or custom GraphQL Go implementation |
| Scraping    | `goquery`, `net/http`                                             |
| Data Models | Go structs (`model/StatsResponse`)                                |
| Caching     | Optional internal caching module                                  |

---

## Project Structure

```bash
coding-profile-service/
├── cmd/
│   └── server/
│       └── main.go           # Entry point of the server
├── internal/
│   ├── scraper/
│   │   ├── hackerrank.go
│   │   ├── hackerrankHTMLScraper.go
│   │   ├── codechef.go
│   │   ├── codechefHTMLScraper.go
│   │   ├── leetcode.go
│   │   └── gfg.go
│   ├── cache/
│   │   └── cache.go          # Optional caching logic
│   └── handler/
│       └── stats_handler.go  # GraphQL resolvers / API handlers
├── pkg/
│   └── model/
│       └── stats_response.go # Data model for profiles
├── go.mod
├── go.sum
├── dockerfile                # Dockerfile for containerizing the app
├── .dockerignore             # Excludes unnecessary files from Docker image
└── README.md

```

---

## How It Works

1. **Scraper Layer** (`internal/scraper/`)

   * Each platform has a dedicated scraper module.
   * Scraper can fetch:

     * Public profile info
     * Ratings, problems solved, contests
     * Badges and certifications
   * Scrapers are lightweight, using **Go `net/http` and `goquery`** for HTML parsing.

2. **Handler Layer** (`internal/handler/`)

   * GraphQL resolvers process requests.
   * Queries can specify platform, username, and required fields.

3. **Cache Layer** (`internal/cache/`)

   * Optional: cache results locally to reduce repeated scraping.
   * Can be extended to Redis or other stores.

4. **GraphQL API**

   * Flexible queries: clients request only fields they need.
   * Fast and efficient, fully typed using Go structs.

---

## Example Usage

### GraphQL Query

```graphql
query {
  hackerRank(username: "johnDoe") {
    codingScore
    totalSolved
    badges
    certifications
    certificationLinks
  }
}
```

### Go Client Example

```go
stats, err := scraper.FetchHackerRankHTML("johnDoe")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%+v\n", stats)
```

---

## Getting Started

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/coding-profile-service.git
cd coding-profile-service
```

2. **Install dependencies**

```bash
go mod tidy
```

3. **Run the server**

```bash
go run cmd/server/main.go
```

4. **Query GraphQL**

* Access the GraphQL playground (if configured) at `http://localhost:8080/graphql`
* Execute queries to fetch profile stats.

---

## Future Plans

* Add more coding platforms (Codeforces, AtCoder, etc.).
* Implement **Redis caching** for faster repeated requests.
* Add **unit tests** for scrapers.
* Improve badge extraction to include **SVG images** and more metadata.

---

## Contributing

1. Fork the repository.
2. Create a new branch: `git checkout -b feature/new-scraper`.
3. Add your scraper or feature.
4. Commit changes: `git commit -am 'Add new scraper'`.
5. Push to branch: `git push origin feature/new-scraper`.
6. Create a Pull Request.

```graphql
# ============================
# GraphQL Schema for Coding Profile Service
# ============================

# Root Queryt
type Query {
  # Fetch HackerRank profile
  hackerRank(username: String!): HackerRankProfile

  # Fetch CodeChef profile
  codeChef(username: String!): CodeChefProfile

  # Fetch LeetCode profile
  leetCode(username: String!): LeetCodeProfile

  # Fetch GeeksforGeeks profile
  gfg(username: String!): GFGProfile
}

# ----------------------------
# HackerRank Profile
# ----------------------------
type HackerRankProfile {
  username: String!
  platform: String!
  codingScore: Int
  totalSolved: Int
  badges: [String]           # List of badge titles
  certifications: Int        # Number of certificates
  certificationLinks: [String] # Direct URLs to certificates
}

# ----------------------------
# CodeChef Profile
# ----------------------------
type CodeChefProfile {
  username: String!
  platform: String!
  rating: Int
  maxRating: Int
  globalRank: Int
  countryRank: Int
  contestsParticipated: Int
  totalSolved: Int
}

# ----------------------------
# LeetCode Profile
# ----------------------------
type LeetCodeProfile {
  username: String!
  platform: String!
  totalSolved: Int
  easySolved: Int
  mediumSolved: Int
  hardSolved: Int
  acceptanceRate: Float
  contestRating: Int
  badges: [String]
}

# ----------------------------
# GeeksforGeeks Profile
# ----------------------------
type GFGProfile {
  username: String!
  platform: String!
  totalSolved: Int
  certifications: Int
  badges: [String]
}

# ----------------------------
# Example Usage
# ----------------------------
# query {
#   hackerRank(username: "johnDoe") {
#     codingScore
#     totalSolved
#     badges
#     certifications
#     certificationLinks
#   }
# }
```


### 🔹 **API Endpoint**

```
GET https://coding-profile-service.onrender.com/stats
```

| Query Parameter | Description            | Example           |
| --------------- | ---------------------- | ----------------- |
| `leetcode`      | LeetCode username      | `mearjuntripathi` |
| `codechef`      | CodeChef username      | `isthisarjun`     |
| `gfg`           | GeeksforGeeks username | `mearjuntripathi` |
| `hackerrank`    | HackerRank username    | `mearjuntripathi` |

✅ **Example Request:**

```
https://coding-profile-service.onrender.com/stats?leetcode=mearjuntripathi&codechef=isthisarjun&gfg=mearjuntripathi&hackerrank=mearjuntripathi
```

✅ **Example Response:**

```json
{
    "profiles": [
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
            "badges": [
                "Problem Solving",
                "CPP",
                "Java",
                "Python",
                "Days of Code",
                "Sql",
                "C language"
            ],
            "certifications": 4,
            "certificationLinks": [
                "https://www.hackerrank.com/certificates/dd88f94012d9",
                "https://www.hackerrank.com/certificates/ad7f9b3ad2e1",
                "https://www.hackerrank.com/certificates/b78dde45a6f8",
                "https://www.hackerrank.com/certificates/7fdef080e935"
            ]
        },
        {
            "platform": "codeforces",
            "username": "isthisarjun",
            "totalSolved": 7,
            "rating": 860,
            "contestsParticipated": 3,
            "maxRating": 860,
            "questionsByType": {
                "easy": 7
            }
        },
        {
            "platform": "leetcode",
            "username": "mearjuntripathi",
            "totalSolved": 710,
            "easySolved": 253,
            "mediumSolved": 427,
            "hardSolved": 30
        }
    ]
}
```

---

### Notes:

1. Each platform has its **own type** because fields differ slightly.
2. All profiles include `username` and `platform` for clarity.
3. Arrays like `badges` and `certificationLinks` can be extended later.
4. GraphQL allows **partial queries**, so clients only fetch the fields they need.
5. Future platforms like Codeforces or AtCoder can be added by defining a new type and adding a root query.

---

## License

This project is open-source under **MIT License**.

---
