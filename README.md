
## s3db

s3db is a lightweight Golang package that allows you to leverage the power of AWS S3 as a database. It provides a simple and efficient solution for storing and retrieving records using the highly reliable and cost-effective AWS S3 object store.

### Why s3db?

AWS S3 is renowned for its exceptional features such as cost-effectiveness, strong consistency, and ease of use. It is widely considered one of the top choices for object storage in the market. However, the question arises: Can we utilize this object store as a database? That's where s3db comes in.

### Features

- Utilizes the AWS S3 object store as a database.
- Lightweight and straightforward Golang package.
- Leverages the benefits of AWS S3, including cost-effectiveness and strong consistency.

### Endpoints

s3db includes a simple Gin application with the following endpoints:
1. **`POST /records/:id`**: Add a new record to the database using this endpoint.
2. **`GET /records/:id`**: Retrieve the value associated with a specific id.
3. **`GET /records`**: Get all IDs of the records stored in the database.
4. **`POST /drop-db`**: Delete all stored records in the database.


### Getting Started

To start using s3db, follow these steps:

1. Set up your AWS S3 credentials as environment variables.
2. Launch the included Gin application using `go run main.go`
3. Use the provided endpoints to manage your records in the AWS S3 object store.
4. Profit! 


### Scalability

Theoretical scalability of s3db is as follows:
- **`Read TPS`**: 5,500 GET requests per second.
- **`Write TPS`**: 3,500 PUT/POST/DELETE requests per second.
- **`Monthly Reads`**: 5,500 * 60 * 60 * 24 * 30 = 475,200,000 requests.
- **`Monthly Writes`**: 3,500 * 60 * 60 * 24 * 30 = 302,400,000 requests.
