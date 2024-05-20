package main

import (
	"context"
	"errors"
	"fmt"
	"goql/src/configuration"
	"goql/src/database"
	"goql/src/database/dbmodel"

	_ "github.com/lib/pq"
)

type App struct {
	ctx     context.Context
	configs configuration.AppConfiguration
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	confs, err := configuration.LoadFromFile()
	if err != nil {
		panic("failed to load configurations")
	}
	a.configs = *confs
}

func (a *App) GetCurrentConfigurations() []string {
	names := []string{}

	for _, conf := range a.configs.Connections {
		names = append(names, conf.Name)
	}

	return names
}

func (a *App) GetConfiguration(name string) *configuration.DatabaseConnection {
	conf, exists := a.configs.ConfigurationExists(name)
	if exists {
		return conf
	}
	return nil
}

func (a *App) AddConfiguration(conf configuration.DatabaseConnection) error {
	if _, exists := a.configs.ConfigurationExists(conf.Name); exists {
		return errors.New("configuration name already used")
	}
	a.configs.AddConfiguration(conf)
	return nil
}

func (a *App) TestConnection(conf configuration.DatabaseConnection) error {
	conn, err := conf.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (a *App) SetActiveConnection(connectionId int) error {
	fmt.Printf("Updating connection to %d\n", connectionId)
	conf := a.configs.GetConfFromId(connectionId)
	if conf == nil {
		fmt.Printf("Connection %d does not exists\n", connectionId)
		return fmt.Errorf("connection %d does not exists", connectionId)
	}
	conn, err := conf.Connect()
	if err != nil {
		fmt.Printf("Failed to connect: %s\n", err.Error())
		return err
	}
	a.configs.ActiveConnection = &database.PostgresConnection{
		Id:   connectionId,
		Conn: conn,
	}
	fmt.Println("Updated active connection")
	return nil
}

func (a *App) GetAvailableTables() (*[]dbmodel.Table, error) {
	fmt.Println("Fetching available tables")
	conn := a.configs.ActiveConnection
	if conn == nil {
		fmt.Println("No active connections")
		return nil, fmt.Errorf("no active connection")
	}
	return conn.GetTables()
}

func (a *App) GetTableData(schema string, table string) (*database.TableData, error) {
	conn := a.configs.ActiveConnection
	if conn == nil {
		return nil, fmt.Errorf("no connections active")
	}

	data, err := conn.GetTableData(schema, table)
	if err != nil {
		fmt.Printf("failed to get data for table: %s\n", err.Error())
		return nil, err
	}
	return data, nil
}

func (a *App) UpdateValue(request database.UpdateValueRequest) error {
	conn := a.configs.ActiveConnection
	if conn == nil {
		return fmt.Errorf("no connections active")
	}

	return conn.UpdateValue(request)
}
