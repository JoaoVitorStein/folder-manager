## Instructions to Run
To run this project you are going to need Docker, Docker Compose.
When you have all the dependencies you can execute the command:
```
docker-compose up --build
``` 

With this command the application will be built and started with a PostgreSQL database.

In the application startup the data that is in csv_file.csv will be imported to the database.


## Using the API
The endpoint is located at localhost:8080/folders and it has the following query parameters:

- filter(optional): string to filter by folder name
- sort: field to apply the sort (if you want to sort descending just add a "-" before the field name)
- page: page for pagination
- size: page size

### Examples

Request for all folders that start with test sorted by priority ascending in page 1:
```
curl "localhost:8080/folders?filter=test&sort=priority&page=1&size=10"
```


Request for all folders with no name filter sorted by priority descending in page 1:
``
curl "localhost:8080/folders?filter=&sort=-priority&page=1&size=10"
``

## Notes
If you can't run the application try changing the volume mapping in [docker-compose](docker-compose.yml#32) to a valid path in your machine