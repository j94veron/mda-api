package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	_findPendingOrders = "select a.nrointerno, a.token, a.devolucion, a.nro_aremito, a.cuit_destino, a.fecha_emision, b.nro_registro, b.cantidad, b.contenido from  test_db.cab_aremito a inner join test_db.cue_aremito b on a.nrointerno = b.nrointerno and a.status = 'P'"
	_updateOrder       = "update test_db.cab_aremito set status =? , error_message =?  where nrointerno =? "
)

type OrderRepository interface {
	FindPendingOrders() ([]Order, error)
	UpdateOrder(orderId int, newError string, newMessage string) error
}

type orderRepositoryImpl struct {
	dbConnection *sql.DB
}

func (p *orderRepositoryImpl) UpdateOrder(orderId int, newError string, newMessage string) error {
	row, err := p.dbConnection.Prepare(_updateOrder)
	if err != nil {
		return err
	}

	defer func(row *sql.Stmt) {
		err := row.Close()
		if err != nil {

		}
	}(row)

	// Execute the update query
	_, err = row.Exec(newError, newMessage, orderId)
	if err != nil {
		log.Printf("Error al ejecutar la consulta de actualización: %v", err)
		return fmt.Errorf("error al ejecutar consulta %v", err)
	}
	return nil
}

func (p *OrderRepository) FindPendingOrders() ([]Order, error) {
	conn, err :=
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
		var result Row
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
