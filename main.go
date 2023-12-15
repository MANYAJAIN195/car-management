package main

import (
	"gofr.dev/pkg/gofr"
)

type Car struct {
	ID     int    `json:"ID"`
	Owner  string `json:"Owner"`
	CarNo  string `json:"CarNo"`
	Status string `json:"Status"`
}

func main() {
	app := gofr.New()

	app.POST("/caradd/{name}/{number}/{status}", func(ctx *gofr.Context) (interface{}, error) {
		name := ctx.PathParam("name")
		number := ctx.PathParam("number")
		status := ctx.PathParam("status")

		// Inserting a car row in the database using SQL
		_, err := ctx.DB().ExecContext(ctx, "INSERT INTO car (Owner, CarNo, Status) VALUES (?, ?, ?)", name, number, status)

		return nil, err
	})

	app.GET("/carlist", func(ctx *gofr.Context) (interface{}, error) {
		var cars []Car

		// Getting cars from the database using SQL
		rows, err := ctx.DB().QueryContext(ctx, "SELECT * FROM car")
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var car Car
			if err := rows.Scan(&car.ID, &car.Owner, &car.CarNo, &car.Status); err != nil {
				return nil, err
			}

			cars = append(cars, car)
		}

		// return the cars
		return cars, nil
	})

	app.PUT("/carupdate/{number}/{status}", func(ctx *gofr.Context) (interface{}, error) {
		number := ctx.PathParam("number")
		status := ctx.PathParam("status")

		// Updating the status of the car in the database using SQL
		_, err := ctx.DB().ExecContext(ctx, "UPDATE car SET Status = ? WHERE CarNo = ?", status, number)

		return nil, err
	})

	app.GET("/carinfo/{number}", func(ctx *gofr.Context) (interface{}, error) {
		number := ctx.PathParam("number")

		// Retrieving car information based on car number from the database using SQL
		var car Car
		err := ctx.DB().QueryRowContext(ctx, "SELECT * FROM car WHERE CarNo = ?", number).Scan(&car.ID, &car.Owner, &car.CarNo, &car.Status)

		if err != nil {
			return nil, err
		}

		return car, nil
	})

	app.DELETE("/cardelete/{number}", func(ctx *gofr.Context) (interface{}, error) {
		number := ctx.PathParam("number")

		// Deleting car information from the database based on car number using SQL
		_, err := ctx.DB().ExecContext(ctx, "DELETE FROM car WHERE CarNo = ?", number)

		return nil, err
	})

	app.Start()
}
