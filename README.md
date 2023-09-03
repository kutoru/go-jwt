# What is this
Little auth-related project that uses Golang, JWT and MongoDB.
# Setup
Make sure that in the .env file:
- `MONGODB_URI` is the URI to your MongoDB
- `DB_NAME` is a DB name that can be dropped
- `COLL_NAME` is a collection name that can be dropped
- `JWT_TOKEN_KEY` is any string

Next you will need to set up MongoDB. To do that you could either uncomment the `db.Reset()` call in the main function, or enter the following commands in the Mongo shell:
```
use `DB_NAME`
db.createCollection("`COLL_NAME`")
db.`COLL_NAME`.createIndex({"exp": 1}, {expireAfterSeconds: 0})
db.`COLL_NAME`.createIndex({"guid": 1}, {unique: true})
```
# Usage
After starting `main.go`, you can open `client/index.html` in your browser and use it to send requests to the backend. Alternatively, you could manually send the requests, but note that they should be similar to the requests in the `auth()` and `refresh()` functions in the `index.html` script section.\
Also note that the token and cookie expiry time is 30 seconds. It can be changed in the `glb/global.go` file.
