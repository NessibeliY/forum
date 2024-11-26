# Forum

## Description

This project is a web forum where users can:
- Communicate with each other by creating posts and comments.
- Associate posts with categories.
- Like or dislike posts and comments.
- Filter posts by categories, created posts, or liked posts.

## Features

1. **User Authentication:**
    - Users can register by providing an email, username, and password.
    - Email must be unique, and password will be encrypted (bonus).
    - Users can log in, and their session will be tracked using cookies.
    - Each session has an expiration date.

2. **Posts and Comments:**
    - Registered users can create posts and comments.
    - Posts can be associated with one or more categories.
    - Both posts and comments are visible to all users (registered or not).

3. **Likes and Dislikes:**
    - Registered users can like or dislike posts and comments.
    - The number of likes and dislikes is visible to all users.

4. **Filtering:**
    - Users can filter posts by:
        - Categories (acting as subforums).
        - Posts created by the logged-in user.
        - Posts liked by the logged-in user.

## Technology Stack

- **Backend:** Go (Golang)
- **Database:** SQLite
- **Encryption:** Bcrypt (Bonus)
- **Sessions & Cookies:** HTTP cookies for session management
- **Containerization:** Docker
- **UUID:** Used to uniquely identify user sessions (Bonus)

## Setup

### Prerequisites
- Go installed (version 1.16 or later)
- Docker installed
- SQLite3 installed
- `github.com/mattn/go-sqlite3` for SQLite integration
- `golang.org/x/crypto/bcrypt` for password encryption (Bonus)
- `github.com/google/uuid` for session management (Bonus)

### Installation

1. Clone the repository:
   ```bash
   git clone git@git.01.alem.school:nyeltay/forum.git
   cd forum
   ```

2. Run the program:
   ```bash
   go run .
   ```
   
Or run the project with Docker:
   ```bash
    make build
    make run
   ```

## Authors

- [@algaliyev](https://01.alem.school/git/algaliyev)
- [@nyeltay](https://01.alem.school/git/nyeltay)