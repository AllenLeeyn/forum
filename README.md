# Forum

Forum is an application built with Go and SQLite that provides user accounts, posts, comments, and feedback (like/dislike) functionality.

## Table of Contents

- [About](#about)
- [Usage](#usage)
- [Database Schema](#database-schema)
- [Installation](#installation)
- [Dependencies](#dependencies)
- [Pages & Features](#pages--features)
- [Credits](#credits)

---

## About

This project is a **forum system** built in **Go** using the standard `net/http` library and **SQLite** for lightweight data management. The application supports:

- **User Accounts & Sessions:**
  Users register, log in, and receive session cookies to remain authenticated. Only one active session is allowed at a time.
  
- **Posts:**
  Registered users can create posts that include a title and content. Posts can be assigned to one or more categories for easy filtering.
  
- **Comments:**
  Users can comment on posts.
  
- **Feedback:**
  Both posts and comments support feedback in the form of likes and dislikes. Feedback data is stored separately for posts and comments.
  
- **Template-Based Rendering & Static Assets:**
  HTML templates render pages (home, login, signup, profile, view post, error, etc.), and static assets (CSS, JavaScript, images) are served from the `/static` URL path (sourced from the `assets/` folder).

A diagram of the data flow is provided below:

![Forum Program Flow Diagram](program_flow.png)

---

## Usage

To use the forum application:

1. **Register or Log In:**
   Create an account or log in to access posting, liking, disliking, and commenting features.
2. **Create & Categorize Posts:**
   Write posts and assign one or more categories.
3. **Comment:**
   Comment on posts.
4. **Provide Feedback:**
   Like or dislike posts and comments.
5. **Browse & Filter:**
   View posts by category, by your own posts, or by posts you have liked.

---

## Database Schema

Below is an **Entity-Relationship Diagram (ERD)** for the forum database:

![Forum ER Diagram](forum_erd.png)

### **User**

- **Attributes:**
  Each user has an ID, name, email, password hash, and account type (e.g., User, Moderator, Administrator).
- **Relationships:**
  - A user may only have one active session at a time.
  - A user can create multiple posts and comments.
  - A user can provide feedback on posts and comments.

### **Session**

- **Purpose:**
  Manages user authentication with a session ID, expiration time, and last access timestamp.

### **Post**

- **Attributes:**
  A post includes a title, body, creation timestamp, and a list of category IDs.
- **Relationships:**
  - Each post is created by one user.
  - Posts can receive comments and feedback.
  - Posts belong to one or more categories.

### **Comment**

- **Attributes:**
  A comment includes its content and creation timestamp.
- **Relationships:**
  - Each comment is associated with a specific post.
  - Comments are root-level.
  - Feedback on comments is stored separately.

### **Feedback**

- **Purpose:**
  Feedback (like/dislike) can be applied to both posts and comments.
- **Relationships:**
  - Each piece of feedback is associated with one user and one post or comment.

---

## Installation

To set up the forum locally, follow these steps:

```sh
# Clone the repository
git clone https://01.gritlab.ax/git/ylee/forum
cd forum

# Run the application
go run main.go
```

The application will be available at [http://localhost:8080](http://localhost:8080).

---

## Dependencies

This project requires the following Go dependencies:

- `github.com/gofrs/uuid v4.4.0+incompatible`
- `github.com/mattn/go-sqlite3 v1.14.24`
- `golang.org/x/crypto v0.32.0`

These packages will be automatically installed when you run the project (via Go modules).

---

## Pages & Features

The forum application includes the following pages and functionality:

- **Home:**
  Displays a feed of posts along with filtering options (e.g., by category, user posts, or liked posts).

- **Login:**
  A page where existing users authenticate.

- **Sign Up:**
  A registration page for new users, including password strength validation and email uniqueness checks.

- **Profile:**
  Displays user details and lists posts created by the user.

- **New Post:**
  A popup where authenticated users can create a new post and select one or more categories.

- **View Post:**
  Displays the full content of a post, along with its comments and interactive feedback buttons.

- **Feedback:**
  Allows users to like or dislike posts and comments.

- **Logout:**  
  Ends the current user session and redirects to Home.

- **Terms:**  
  Displays the forum’s terms and conditions.

- **Error:**  
  An error page for displaying messages when an error occurs.

---

## Credits

This project was developed by **Allen, Anass, Johannes, Milli, and Richard**.

