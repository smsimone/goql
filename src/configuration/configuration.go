package configuration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"goql/src/database"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

const (
	configFile = "config.txt"
)

type DatabaseConnection struct {
	Id       *int   `json:"id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Port     int64  `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type AppConfiguration struct {
	Connections      []DatabaseConnection
	ActiveConnection database.ActiveConnection
}

func (c *DatabaseConnection) GenerateConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", c.Username, c.Password, c.Url, c.Port, c.Database)
}

func (c *DatabaseConnection) Connect() (*sql.Conn, error) {
	connStr := c.GenerateConnectionString()
	fmt.Printf("Connecting to '%s'\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Failed to connect to database: %s\n", err.Error())
		return nil, fmt.Errorf("failed to connect to database: %s", err.Error())
	}

	defer db.Close()
	conn, err := db.Conn(context.Background())
	if err != nil {
		fmt.Printf("Failed to get connection: %s\n", err.Error())
		return nil, fmt.Errorf("failed to open connection: %s", err.Error())
	}

	fmt.Printf("Succesfully connected to database: %p\n", db)
	return conn, nil
}

func (a *AppConfiguration) ConfigurationExists(name string) (*DatabaseConnection, bool) {
	for _, conf := range a.Connections {
		if conf.Name == name {
			return &conf, true
		}
	}
	return nil, false
}

func (a *AppConfiguration) GetConfFromId(id int) *DatabaseConnection {
	for _, conf := range a.Connections {
		if *conf.Id == id {
			return &conf
		}
	}
	return nil
}

func (a *AppConfiguration) AddConfiguration(conf DatabaseConnection) {
	defer a.SaveToFile()

	if conf.Id != nil {
		for i, val := range a.Connections {
			if *conf.Id == *val.Id {
				a.Connections[i] = conf
				return
			}
		}
	}

	id := len(a.Connections) + 1
	conf.Id = &id
	a.Connections = append(a.Connections, conf)
}

func (a *AppConfiguration) SaveToFile() {
	jsonValues := []string{}
	for _, conn := range a.Connections {
		bytes, _ := json.Marshal(conn)
		jsonValues = append(jsonValues, string(bytes))
	}

	content := strings.Join(jsonValues, "\n")
	err := os.WriteFile(configFile, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %s\n", err.Error())
	}
}

func LoadFromFile() (*AppConfiguration, error) {
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Failed to read file: %s\n", err.Error())
		return nil, err
	}

	content := string(bytes)
	configs := []DatabaseConnection{}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		var tmp DatabaseConnection
		if err := json.Unmarshal([]byte(line), &tmp); err != nil {
			fmt.Printf("Failed to parse json: %s\n", err.Error())
			continue
		}
		if tmp.Id != nil {
			configs = append(configs, tmp)
		}
	}
	fmt.Printf("Loaded %d configurations\n", len(configs))

	confs := AppConfiguration{Connections: configs}

	defer confs.SaveToFile()

	return &confs, nil
}
