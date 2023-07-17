
## s3db

s3db is a lightweight Golang package that allows you to leverage the power of AWS S3 as a key-value database. It provides a simple and efficient solution for storing and retrieving key-value pairs using the highly reliable and cost-effective AWS S3 object store.

### Why s3db?

AWS S3 is renowned for its exceptional features such as cost-effectiveness, strong consistency, and ease of use. It is widely considered one of the top choices for object storage in the market. However, the question arises: Can we utilize this object store as a key-value database? That's where s3db comes in.

### Features

- Utilizes the AWS S3 object store as a key-value database.
- Lightweight and straightforward Golang package.
- Leverages the benefits of AWS S3, including cost-effectiveness and strong consistency.

### Endpoints

s3db includes a simple Gin application with the following endpoints:
1. **`/drop-db`**: This endpoint deletes all key-value pairs stored in the database.
2. **`/get-record`**: Use this endpoint to retrieve the value associated with a specified key.
3. **`/new-record`**: Add a new key-value pair to the database using this endpoint.

### Getting Started

To start using s3db, follow these steps:

1. Set up your AWS S3 credentials as environment variables.
2. Launch the included Gin application using `go run main.go`
3. Use the provided endpoints to manage your key-value pairs in the AWS S3 object store.
4. Profit! 
