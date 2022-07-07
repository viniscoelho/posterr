# posterr

## How to run the tests
To check the test coverage, run the following command from the root folder:
```
go test $(go list ./... | grep -v /vendor/) -coverprofile cover.out && go tool cover -html=cover.out
```

This will create a `cover.out` file which will be opened in browser tab, showing detailed coverage status for each package.
This information is also available in the command line.

## How to use it
For this project, I've used PostgresSQL as a database. You can install it from [here](https://www.postgresql.org/download/) and choose the appropriate version for your OS.

Run `go build -o posterr .` in the `src` folder. The following arguments are available:
- `--init-db` - Initializes the database. Required at first run.
- `--port` - Sets the application port.

For example:
```bash
go build -o posterr .
./posterr --init-db --port 4000
```

## Planning

### Questions
 - Is there a limit on how many replies could be done within a day?
 - May replies be reposted and quote-reposted?
 - May replies be replied?  
 - In which order should the replies being shown on this new feed? Ascending or Descending?

### Implementation
In the database layer, I could treat each reply as a post as well. However, it would be necessary to update the `posts` database schema, adding an extra column `reply_id`. `reply_id` contains the id of the post which the reply is from. In this sense, the content of a post, where reply is not null, will be the message used to reply a post. Also, the same logic used to repost and quote-repost could be used for a reply, as it is a post as well. The table would be something like this:

```sql
CREATE TABLE posts(
        post_id VARCHAR (36) PRIMARY KEY,
        username VARCHAR (14) NOT NULL REFERENCES users (username),
        content VARCHAR (777) NULL,
        reposted_id VARCHAR (36) NULL,
        reply_id VARCHAR (36) NULL,
        created_at TIMESTAMPTZ DEFAULT NOW(),
        FOREIGN KEY (reposted_id) REFERENCES posts (post_id)
        FOREIGN KEY (reply_id) REFERENCES posts (post_id))
```

For the actual code, it would be necessary to update the `Posterr` interface to include a function to create replies. Maybe the same `WriteContent` function could be used, by adding an extra parameter, but the logic would get more complex. With some refactor, maybe not too complex.

The API implementation would quite similar to `CreateContent`, or, depending on the changes above, could be reused to handle it as well.

Finally, no changes would be necessary in the search mechanism with this implementation, since replies are treated as posts.

## Critique

### Testing
The code is partially covered by unit tests. Some functions like `ListHomePageContent`, `ListProfileContent`, `SearchContent` and `GetUserProfile` do not have unit testing simply because the setup was quite big. Thus, I used some unit tests to populate the database and tested each manually. I managed to find some bugs and fixed them. Also, I created mocks for the interfaces to use for testing the APIs. However, I could not complete it.

### Extra feature
The first phase was straight-forward to implement. I only had a few issues with the language when I was trying to use the `LIKE` key from Postgres. Thus, I took some time to figure out how to make this work in Go.

### Remarks
I`ve added a caching mechanism to store no. of followers and following users, assuming that a user will not run many of these operations in a short period of time. However, I created a small load unit tests and the connection pool breaks for many requests at once. The connection pool could be improved to be long living ones, thus reused, but some changes would be necessary to handle this, but I'm not entirely sure this would solve the problem.

### Scaling
For scaling this project, the first thing would be to replicate the service. However, to make this possible, it would be necessary a load balancer to coordinate the incoming traffic and redirect the requests to each server. I could ship it to a cloud provider, such as AWS or Azure, which already has such functionality and make usage of it. The cloud provider would be also responsible for provisioning servers, which could also de-provision if the load is below some defined threshold. For example, I would containerize the project into docker images and use Azure Kubernetes Service (AKS) as orchestrator of the containers.

The next phase, I would replicate the database, as the traffic load would increase with the replicated servers. The main database could responsible for accepting only write/update traffic and the additional ones would handle the read operations. The additional databases would read from the main one to keep their data up-to-date. In case of a failure in the main database, one of the additionals would be promoted to main to keep the system running. Some caching could be added as well to reduce latency. In this case, a NoSQL database, like MongoDB, could store some operations that do not change as much, such as followers and following, unless there is a celebrity-case which will increase the load of these operations.

For another phase of scaling, the data could be partitioned into shards, but it would require some extra logic to generate the hash function using the username. Also, depending on the traffic of a user or group of users, these could have their own shards to improve the traffic load. Hopefully, it won't be necessary to reshard the data, which adds a problem of recreating the hash function and necessity of moving around data to be in accordance to the new function.
