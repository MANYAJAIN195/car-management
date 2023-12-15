package main

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"gofr.dev/pkg/datastore"
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/request"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	go main()
	time.Sleep(3 * time.Second)

	tests := []struct {
		desc       string
		method     string
		endpoint   string
		statusCode int
		body       []byte
	}{
		{"get cars", http.MethodGet, "carlist", http.StatusOK, nil},
		{"post cars", http.MethodPost, "caradd", http.StatusCreated, []byte(`{
			"Owner": "Manya Jain",
			"CarNo": "abc xyz 3300",
			"Status": "in process"
		}`),
		},
		{"get a car", http.MethodGet, "carinfo", http.StatusOK, []byte(`{
			"CarNo": "abc xyz 3300"
		}`),
		},
		{"update cars", http.MethodPut, "carupdate", http.StatusOK, []byte(`{
			"CarNo": "abc xyz 3300",
			"Status": "done"
		}`),
		},
		{"delete cars", http.MethodDelete, "cardelete", http.StatusNoContent, []byte(`{
			"CarNo": "abc xyz 3300"
		}`)},
	}

	for i, tc := range tests {
		req, _ := request.NewMock(tc.method, "http://localhost:8080/"+tc.endpoint, bytes.NewBuffer(tc.body))

		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			t.Errorf("TEST[%v] Failed.\tHTTP request encountered Err: %v\n%s", i, err, tc.desc)
			continue
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("TEST[%v] Failed.\tExpected %v\tGot %v\n%s", i, tc.statusCode, resp.StatusCode, tc.desc)
		}

		_ = resp.Body.Close()
	}
}

func TestCoreLayer(*testing.T) {
	app := gofr.New()

	// initializing the seeder
	seeder := datastore.NewSeeder(&app.DataStore, "../db")
	seeder.ResetCounter = true

	createTable(app)
}

func createTable(app *gofr.Gofr) {
	// drop table to clean previously added id's
	_, err := app.DB().Exec("DROP TABLE IF EXISTS car;")

	if err != nil {
		return
	}

	_, err = app.DB().Exec("CREATE TABLE IF NOT EXISTS car " +
		"(id INT AUTO_INCREMENT PRIMARY KEY, Owner VARCHAR(255) NOT NULL, CarNo VARCHAR(255) NOT NULL, Status VARCHAR(255) NOT NULL);")
	if err != nil {
		return
	}
}

func TestAddCar(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		ctx.Logger.Error("mock connection failed")
	}

	ctx.DataStore = datastore.DataStore{ORM: db}
	ctx.Context = context.Background()
	tests := []struct {
		desc    string
		car     Car
		mockErr error
		err     error
	}{
		{"Valid case", Car{Owner: "Test123", CarNo: "BH xyz 2208", Status: "done"}, nil, nil},
		{"DB error", Car{Owner: "Test234", CarNo: "BH xyz 2201"}, errors.DB{}, errors.DB{Err: errors.DB{}}},
	}

	for i, tc := range tests {
		// Set up the expectations for the INSERT query
		mock.ExpectExec("INSERT INTO car (Owner, CarNo, Status) VALUES (?, ?, ?)").
			WithArgs(tc.car.Owner, tc.car.CarNo, tc.car.Status).
			WillReturnResult(sqlmock.NewResult(2, 1)).
			WillReturnError(tc.mockErr)

		// Set up the expectations for the SELECT query
		rows := sqlmock.NewRows([]string{"id", "Owner", "CarNo", "Status"}).
			AddRow(tc.car.ID, tc.car.Owner, tc.car.CarNo, tc.car.Status)
		mock.ExpectQuery("SELECT * FROM car WHERE CarNo = ?").
			WithArgs(tc.car.CarNo).
			WillReturnRows(rows).
			WillReturnError(tc.mockErr)

		resp, err := Create(ctx, tc.car)

		ctx.Logger.Log(resp)
		assert.IsType(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestGetCarList(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		ctx.Logger.Error("mock connection failed")
	}

	ctx.DataStore = datastore.DataStore{ORM: db}
	ctx.Context = context.Background()

	tests := []struct {
		desc    string
		car     []Car
		mockErr error
		err     error
	}{
		{"Valid case with cars", []Car{
			{ID: 1, Owner: "Manya jain", CarNo: "UP HGS 7845", Status: "done"},
			{ID: 2, Owner: "Ayushi Jain", CarNo: "DL SDV 4815", Status: "in progress"},
		}, nil, nil},
		{"Valid case with no car", []Car{}, nil, nil},
		{"Error case", nil, errors.Error("database error"), errors.DB{Err: errors.Error("database error")}},
	}

	for i, tc := range tests {
		rows := sqlmock.NewRows([]string{"id", "Owner", "CarNo", "Status"})
		for _, c := range tc.car {
			rows.AddRow(c.ID, c.Owner, c.CarNo, c.Status)
		}

		mock.ExpectQuery("SELECT * FROM car").WillReturnRows(rows).WillReturnError(tc.mockErr)

		resp, err := GetCarList(ctx)
		ctx.Logger.Log(resp)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestGetCarInfo(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		ctx.Logger.Error("mock connection failed")
	}

	ctx.DataStore = datastore.DataStore{ORM: db}
	ctx.Context = context.Background()

	tests := []struct {
		desc    string
		car     Car
		mockErr error
		err     error
	}{
		{"Valid case", Car{CarNo: "BH xyz 2208"}, nil, nil},
		{"DB error", Car{CarNo: "BH xyz 2200"}, errors.DB{}, errors.DB{Err: errors.DB{}}},
	}
	for i, tc := range tests {
		rows := sqlmock.NewRows([]string{"id", "Owner", "CarNo", "Status"}).
			AddRow(tc.car.ID, tc.car.Owner, tc.car.CarNo, tc.car.Status)
		mock.ExpectQuery("SELECT * FROM car WHERE CarNo = ?").
			WithArgs(tc.car.CarNo).
			WillReturnRows(rows).
			WillReturnError(tc.mockErr)

		resp, err := Get(ctx, tc.car)
		ctx.Logger.Log(resp)
		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestUpdateCar(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		ctx.Logger.Error("mock connection failed")
	}

	ctx.DataStore = datastore.DataStore{ORM: db}
	ctx.Context = context.Background()

	tests := []struct {
		desc    string
		car     Car
		mockErr error
		err     error
	}{
		{"Valid case", Car{CarNo: "BH xyz 2208", Status: "processing"}, nil, nil},
		{"DB error", Car{ID: 16}, errors.DB{}, errors.DB{Err: errors.DB{}}},
	}

	for i, tc := range tests {
		// Set up the expectations for the UPDATE query
		mock.ExpectExec("UPDATE car SET Status = ? WHERE CarNo = ?").
			WithArgs(tc.car.Status, tc.car.CarNo).
			WillReturnResult(sqlmock.NewResult(2, 1)).
			WillReturnError(tc.mockErr)

		// Set up the expectations for the SELECT query
		rows := sqlmock.NewRows([]string{"id", "Owner", "CarNo", "Status"}).
			AddRow(tc.car.ID, tc.car.Owner, tc.car.CarNo, tc.car.Status)
		mock.ExpectQuery("SELECT * FROM car WHERE CarNo = ?").
			WithArgs(tc.car.CarNo).
			WillReturnRows(rows).
			WillReturnError(tc.mockErr)

		resp, err := Update(ctx, tc.car)

		ctx.Logger.Log(resp)
		assert.IsType(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func TestDeleteCar(t *testing.T) {
	ctx := gofr.NewContext(nil, nil, gofr.New())
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		ctx.Logger.Error("mock connection failed")
	}

	ctx.DataStore = datastore.DataStore{ORM: db}
	ctx.Context = context.Background()

	tests := []struct {
		desc    string
		car     Car
		mockErr error
		err     error
	}{
		{"Valid case", Car{CarNo: "BH xyz 2208"}, nil, nil},
		{"DB error", Car{CarNo: "BH xyz 2200"}, errors.DB{}, errors.DB{Err: errors.DB{}}},
	}
	for i, tc := range tests {
		// Set up the expectations for the DELETE query
		mock.ExpectExec("DELETE FROM car WHERE CarNo = ?").
			WithArgs(tc.car.CarNo).
			WillReturnResult(sqlmock.NewResult(2, 1)).
			WillReturnError(tc.mockErr)

		resp, err := Delete(ctx, tc.car)

		ctx.Logger.Log(resp)
		assert.IsType(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}
