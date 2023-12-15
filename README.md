# GO CAR GARAGE MANAGEMENT

I have developed a mini-project using the GoFr framework to build a simple HTTP API for a Car Garage Management service. The API includes CRUD operations for managing entries of cars in the garage and integrates with a MySQL database using docker for persistent data storage. The key functionalities implemented include adding a new entry to the database when a car enters the garage (capturing owner name, car number, and status), retrieving the list of currently parked cars, fetching details of a specific car based on its number, updating the entry when a car undergoes repair with relevant status parameters, and finally, deleting the entry from the database when the car leaves the garage using its number. Furthermore, activity diagrams and UML diagrams are also created to provide a visual representation of the system's interactions and structure, understanding of the project architecture.

## Requirements
- GoFr Framework
- MySQL Database

## Project Structure
/car-management
|-- main.go
|-- main_test.go
|-- README.md
|-- go.mod
|-- go.sum
|-- /configs
| |-- .env

- **main.go**: Entry point of the application and Implementation of the HTTP handlers.
- **main_test.go**: Unit tests for the main functionalities.
- **config/.env**: Database connection and port settings.

## Getting Started

To download gofr, use the command: 
```
go get gofr.dev
```

To download and sync the required modules, use the command:
```
go mod tidy
```

run the mysql server and create a database locally using the following docker command:
```
docker run --name gofr-mysql -e MYSQL_ROOT_PASSWORD=root123 -e MYSQL_DATABASE=test_db -p 3306:3306 -d mysql:8.0.30
```

Access test_db database and create table car:
```
docker exec -it gofr-mysql mysql -uroot -proot123 test_db -e "CREATE TABLE car (id INT AUTO_INCREMENT PRIMARY KEY, Owner VARCHAR(255) NOT NULL, CarNo VARCHAR(255) NOT NULL, Status VARCHAR(255) NOT NULL);"
```

To run the server, use the command:
```
go run main.go
```

Note: Access the server at port 8080 which can be modified in .env file.

## API Endpoints

- **POST /caradd**: Create a new entry for a car in the garage.
- **GET /carlist**: Retrieve the list of cars currently in the garage.
- **GET /carinfo**: Retrieve details of a specific car by its number.
- **PUT /carupdate**: Update the status of a car.
- **DELETE /cardelete**: Delete the entry of a car.

## Postman Collection
[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://god.gw.postman.com/run-collection/31394230-227e3a50-cd31-4313-a608-9fdaa2b4f647?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D31394230-227e3a50-cd31-4313-a608-9fdaa2b4f647%26entityType%3Dcollection%26workspaceId%3D593f0adc-3224-4644-ad3e-caec4d690d43)

### UML CLASS DIAGRAM
![uml](https://github.com/MANYAJAIN195/car-management/assets/71972339/f6e6c940-7609-45d7-aa6c-66a288bfa9dd)


### UML ACTIVITY DIAGRAM
![activity diagram](https://github.com/MANYAJAIN195/car-management/assets/71972339/0fb76b4c-7c51-4ece-886f-cc3d3fed6075)
