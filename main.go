// main.go
package main

import (
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
)

type Car struct {
	ID     int    `json:"ID"`
	Owner  string `json:"Owner"`
	CarNo  string `json:"CarNo"`
	Status string `json:"Status"`
}

func GetCarList(ctx *gofr.Context) (interface{}, error) {
	//var cars []Car

	// Getting cars from the database
	rows, err := ctx.DB().QueryContext(ctx, "SELECT * FROM car")
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	defer rows.Close()
	cars := make([]Car, 0)
	// Iterate over the rows of the result set
	for rows.Next() {
		var car Car
		// Scan the values from the current row into the Car struct fields
		if err := rows.Scan(&car.ID, &car.Owner, &car.CarNo, &car.Status); err != nil {
			return nil, errors.DB{Err: err}
		}

		cars = append(cars, car)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return cars, nil
}

func Get(ctx *gofr.Context, newCar Car) (Car, error) {

	var car Car
	// execute the SELECT query and get new car data
	err := ctx.DB().QueryRowContext(ctx, "SELECT * FROM car WHERE CarNo = ?", newCar.CarNo).Scan(&car.ID, &car.Owner, &car.CarNo, &car.Status)

	if err != nil {
		return Car{}, errors.DB{Err: err}
	}

	return car, nil
}

// GetCarInfo retrieves car information based on car number from the database
func GetCarInfo(ctx *gofr.Context) (interface{}, error) {
	// Read the request body
	var car Car
	if err := ctx.Bind(&car); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := Get(ctx, car)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Create(ctx *gofr.Context, newCar Car) (Car, error) {
	//insert new car data
	_, err := ctx.DB().ExecContext(ctx, "INSERT INTO car (Owner, CarNo, Status) VALUES (?, ?, ?)", newCar.Owner, newCar.CarNo, newCar.Status)
	if err != nil {
		return Car{}, errors.DB{Err: err}
	}

	var car Car
	// execute the SELECT query and get new car data
	err = ctx.DB().QueryRowContext(ctx, "SELECT * FROM car WHERE CarNo = ?", newCar.CarNo).Scan(&car.ID, &car.Owner, &car.CarNo, &car.Status)

	if err != nil {
		return Car{}, errors.DB{Err: err}
	}

	return car, nil
}

// AddCar adds a new car to the database
func AddCar(ctx *gofr.Context) (interface{}, error) {
	// Read the request body
	var car Car
	if err := ctx.Bind(&car); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := Create(ctx, car)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func Update(ctx *gofr.Context, newCar Car) (Car, error) {
	// Updating the status of the car in the database
	_, err := ctx.DB().ExecContext(ctx, "UPDATE car SET Status = ? WHERE CarNo = ?", newCar.Status, newCar.CarNo)
	if err != nil {
		return Car{}, errors.DB{Err: err}
	}

	var car Car
	// execute the SELECT query and get new car data
	err = ctx.DB().QueryRowContext(ctx, "SELECT * FROM car WHERE CarNo = ?", newCar.CarNo).Scan(&car.ID, &car.Owner, &car.CarNo, &car.Status)

	if err != nil {
		return Car{}, errors.DB{Err: err}
	}

	return car, nil
}

// UpdateCarStatus updates the status of the car in the database
func UpdateCarStatus(ctx *gofr.Context) (interface{}, error) {
	// Read the request body
	var car Car
	if err := ctx.Bind(&car); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := Update(ctx, car)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func Delete(ctx *gofr.Context, newCar Car) ([]Car, error) {
	// Deleting car information from the database based on car number
	_, err := ctx.DB().ExecContext(ctx, "DELETE FROM car WHERE CarNo = ?", newCar.CarNo)
	if err != nil {
		return nil, errors.DB{Err: err}
	}
	return nil, nil
}

// DeleteCar deletes car information from the database based on car number
func DeleteCar(ctx *gofr.Context) (interface{}, error) {
	//number := ctx.PathParam("number")

	// Read the request body
	var car Car
	if err := ctx.Bind(&car); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := Delete(ctx, car)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func main() {
	app := gofr.New()

	app.Server.ValidateHeaders = false

	app.GET("/carlist", GetCarList)
	app.GET("/carinfo", GetCarInfo)
	app.POST("/caradd", AddCar)
	app.PUT("/carupdate", UpdateCarStatus)
	app.DELETE("/cardelete", DeleteCar)

	app.Start()
}
