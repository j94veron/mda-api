package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-api/config"
	"log"
)

const (
	_findPendingOrders = "select a.nrointerno, a.token, a.devolucion, a.nro_aremito, a.cuit_destino, a.fecha_emision, b.nro_registro, b.cantidad, b.contenido from  test_db.cab_aremito a inner join test_db.cue_aremito b on a.nrointerno = b.nrointerno and a.status = 'P'"
	_updateOrder       = "update test_db.cab_aremito set status =? , message_error =?  where nrointerno =? "
)

type MysqlRepository interface {
	FindPendingOrders() ([]Order, error)
	UpdateOrder(nroInterno int, newError string, newMessage string) error
}

type mysqlRepositoryImpl struct {
	dataSource string
}

func (p *mysqlRepositoryImpl) UpdateOrder(nroInterno int, newError string, newMessage string) error {
	conn, err := p.getConnection()
	if err != nil {
		return err
	}

	row, err := conn.Prepare(_updateOrder)
	if err != nil {
		return err
	}
	defer row.Close()

	// Execute the update query
	_, err = row.Exec(newError, newMessage, nroInterno)
	if err != nil {
		log.Printf("Error al ejecutar la consulta de actualización: %v", err)
		return err
	}
	log.Printf("Error al ejecutar la consulta de actualización: %v", err)
	return nil
}

func (p *mysqlRepositoryImpl) FindPendingOrders() ([]Order, error) {
	conn, err := p.getConnection()
	if err != nil {
		return nil, err
	}

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	rows, err := conn.Query(_findPendingOrders)
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	orderMap := make(map[int]*Order)
	for rows.Next() {
		var result row
		err := rows.Scan(&result.nroInterno,
			&result.token,
			&result.devolucion,
			&result.nroRemito,
			&result.cuitDestino,
			&result.fechaEmision,
			&result.nroRegistro,
			&result.cantidad,
			&result.contenido,
		)
		if err != nil {
			log.Fatal(err)
		}

		// Check if the user is already in the map
		order, ok := orderMap[result.nroInterno]
		if !ok {
			order = &Order{
				NroInterno:   result.nroInterno,
				NroRemito:    result.nroRemito,
				Devolucion:   result.devolucion,
				Token:        result.token,
				FechaEmision: result.fechaEmision,
				CuitDestino:  result.cuitDestino,
			}
			orderMap[result.nroInterno] = order
		}

		// Append the post to the user's posts
		order.Productos = append(order.Productos, Product{
			NroRegistro: result.nroRegistro,
			Cantidad:    result.cantidad,
			Contenido:   result.contenido,
		})
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return mapToArray(orderMap), nil
}

func mapToArray(inputMap map[int]*Order) []Order {
	var resultArray []Order

	for _, value := range inputMap {
		resultArray = append(resultArray, *value)
	}

	return resultArray
}

// GetConnection - Returns mysql connection.
func (p *mysqlRepositoryImpl) getConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", p.dataSource)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewMysqlRepository(config config.MySQLConfig) MysqlRepository {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DBUser,
		config.DBPass,
		config.DBServer,
		config.DBPort,
		config.DBName,
	)
	return &mysqlRepositoryImpl{
		dataSource: dataSource,
	}
}

type row struct {
	nroInterno   int
	token        string
	devolucion   int
	nroRemito    string
	cuitDestino  string
	fechaEmision string
	nroRegistro  string
	cantidad     int
	contenido    int
}

type Order struct {
	NroInterno   int
	Token        string
	Devolucion   int
	NroRemito    string
	CuitDestino  string
	FechaEmision string
	Productos    []Product
}

type Product struct {
	NroRegistro string
	Cantidad    int
	Contenido   int
}
