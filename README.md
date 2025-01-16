# forum


### SQLite
Below is an ER Diagram for the forum database.

accountType: 
- 1 user
- 2 moderator
- 3 administrator

user:
- A user must have one and only one account type
- A user can have mutliple sessions but only 1 active session.
- A user can creates multiple posts.
- A user can provide multiple feedback.

post:
- A post must have one and only one author
- A post must have a title and body
- A post must have at least one category, but can have multiple categories.
- A post keep tracks the number of comments and likes for easy query to get the information to display on homepage.

comment:
- A comment must have a postId.
- A comment can have no or one parentId.
- A comment without a parentId is a comment at root.
- A comment with a parentId is a reply to a comment.

feedback:
- A feedback must have one and only one author.
- A feedback must have one and only one post

```mermaid
erDiagram
    category{
        int         categoryId PK
        string      name
    }

    post_category{
        int         postId FK
        int         categoryId FK
    }

    user{
        int         userId PK
        int         typeId FK
        string      name
        string      email
        string      passwordHash
        date        registerDate
        dateTime    lastLogin
    }

    post{
        int         postId PK
        int         userId FK
        int         commentCount
        int         likeCount
        int         dislikeCount
        string      title
        text        content
        dateTime    createdAt
    }

    comment{
        int         commentId PK
        int         userId FK
        int         postId FK
        int         parentId(nullable) FK
        text        content
        dateTime    createdAt
    }

    postFeedback{
        int         userId PK
        int         postId PK
        int         rate
        dateTime    createdAt
    }

    commentFeedback{
        int         userId PK
        int         commentId PK
        int         rate
        dateTime    createdAt
    }

    accountType{
        int         typeId PK
        string      name
        boolean     canCreatePost
        boolean     canComment
        boolean     canFeedback
        boolean     canModerate
        boolean     canBanUser
    }
    session{
        UUID        sessionId PK
        int         userId FK
        boolean     isActive
        dateTime    startedAt
        dateTime    expiredAt
        dateTime    lastAccessed
    }

    session }o--|| user : "initiates"
    accountType ||--o{ user : "assigned"
    category ||--o{ post_category : "under"
    post_category }|--|| post : "linked to"
    user ||--o{ post : "creates"
    post ||--o{ comment : "receives"
    post ||--o{ postFeedback : "receives"
    user ||--o{ comment : "makes"
    comment |o--o{ comment : "receives"
    comment ||--o{ commentFeedback : "receives"
    user ||--o{ postFeedback : "provides"
    user ||--o{ commentFeedback : "provides"
```

### pages
- home page (guest/ user)
    - posts without parents (show number of comments and feedbacks)
    - filter by category, created, likes
    - sort by dateCreated, most likes, most comments
    - paginate posts?
- login
- register
- create post (must have at least one category)
- view post (show all comments and feedbacks for specified post)
    - paginate comments

### tasks
- create database
    - list down quries/commands used
    - indexing, required, and unique constraints
- create webpages
    - homepage
    - login
    - register
    - create post
    - view post (comments/feedback happens here)
- frontend logic
    - restriction for login and registration
    - restriction for multiple feedback for same post by one user
    - sorting
- backend logic
    - decide structs for data and user? for page data
    - handler for home page (filtering)
    - handler for login (sessioId with UUID)
    - handler for registration (bycrpt)
    - handler for create post. check if valid sessionId
    - handler for view post
    - handler for (commenting/feedback). check if valid sessionId
    - handler for error
    - function to read database (SELECT query)
    - function to write data base (INSERT query)

### considerations
- https
- tracked invalid login for rate-limit
- session expiration logic
- deleting user, post, comment
- logging?
- share
- image
- limit number of replies to replies
- profile page for change of user details/settings
