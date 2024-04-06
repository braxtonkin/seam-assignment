# API Usage
This API simulates a blog application where users can create blog posts, and comment on specific blog posts.

Depending on the endpoint, the API can expect either a Blog or a Comment in the body of the HTTP request.

**Blogs** are structured as follows
```json
{
  "title" : "titleValue",
  "author" : "authorValue",
  "content" : "contentValue"
}
```

**Comments** are structured similarly 
```json
{
  "author" : "authorValue",
  "content" : "contentValue"
}
```

A list of all entrypoints and their associated HTTP methods can be found here


```
METHOD    ENDPOINT                    REQUEST BODY       RESPONSE                               DESCRIPTION

GET       /posts                      n/a                Array of Blog posts                    Returns an array of all blog objects
GET       /posts/{postID}             n/a                One Blog post                          Retrieves a specific blog post by its ID (postID)
POST      /posts                      Blog               The created Blog post and its id       Creates a new blog post
PUT       /posts/{postID}             Blog               The updated Blog post                  Updates an existing blog post by its ID (postID)
DELETE    /posts/{postID}             n/a                Success message or status              Deletes a blog post by its ID (postID)
GET       /posts/{postID}/comments    n/a                Array of Comments                      Retrieves all comments for a blog post specified by the blog post's ID (postID)  
POST      /posts/{postID}/comments    Comment            The created Comment and its id         Creates a new comment for a blog post specified by the blog post's ID (postID)

```

# Quick Start
## Prerequisites
All that is needed to start the server is [Docker](https://docs.docker.com/desktop/). The Dockerfile and docker-compose.yml have already been provided.

To run the server, first clone the github repository
```
git clone https://github.com/braxtonkin/seam-assignment.git && cd seam-assignment
```
Once inside the github repository, the server can be ran with docker compose


```
docker compose up
```

Once running, the endpoint can be accessed at
```
localhost:4000
```

By default, the port is bound to 4000, but this can be configured.

## Configuration
The port which the server is running can be configured via `.env`, by default, the server runs on **port 4000**.
Configuring the username, password, and database name of the postgres database can also be done in the `.env` file.


