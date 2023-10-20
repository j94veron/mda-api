package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"go-api/config"
)

const (
	findOrderByID = "SELECT * from cab_aremito WHERE nrointerno = ?"
)

type MysqlRepository interface {
	getConnection() (*sql.DB, error)
}

type mysqlRepositoryImpl struct {
	dataSource string
}

// GetConnection - Returns mysql connection.
func (p *mysqlRepositoryImpl) getConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", p.dataSource)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p *mysqlRepositoryImpl) FindOrderById(orderID string) (*Order, error) {
	conn, err := p.getConnection()
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	order := &Order{}
	row := conn.QueryRow(findOrderByID, orderID)
	err = row.Scan(&order.NroInterno,
		&order.Token,
		&order.NroRemito,
		&order.Devolucion,
		&order.CuitDestino,
		&order.FechaEmision)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no se encontró el pedido con ID %d: %w", orderID, err)
		}

		return nil, fmt.Errorf("error al escanear fila: %w", err)
	}

	return order, nil
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

type Order struct {
	NroInterno   string
	Token        string
	Devolucion   string
	NroRemito    string
	CuitDestino  string
	FechaEmision string
}
