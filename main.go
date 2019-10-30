package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
)

var (
	cli     *client.Client
	servers map[string]DdevMySqlServer
)

type DdevMySqlServer struct {
	Name string
	Port uint16
}

func init() {
	servers = make(map[string]DdevMySqlServer)
}

func main() {
	initDockerClient()
	discoverMysqlServers()
	runMysqlServer()
}

func initDockerClient() {
	log.Println("Initializing docker client")
	var err error
	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

func discoverMysqlServers() {
	log.Println("Searching for ddev mysql containers")

	args := filters.NewArgs()
	args.Add("label", "com.docker.compose.service=db")
	args.Add("label", "com.ddev.platform=ddev")
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{Filters: args})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		for _, port := range container.Ports {
			if port.IP == "127.0.0.1" && port.PrivatePort == 3306 {
				dbserver := DdevMySqlServer{}
				dbserver.Name = container.Labels["com.ddev.site-name"]
				dbserver.Port = port.PublicPort
				servers[dbserver.Name] = dbserver
			}
		}
	}
	log.Printf("Found %d db servers", len(servers))
}

func runMysqlServer() {
	log.Println("Starting mysql proxy server")
	l, _ := net.Listen("tcp", "127.0.0.1:32123")
	c, _ := l.Accept()

	conn, _ := server.NewConn(c, "db", "db", DdevHandler{})

	for {
		conn.HandleCommand()
	}
}

type DdevHandler struct {
	server.EmptyHandler
}

func (d DdevHandler) UseDB(dbName string) error {
	// @todo: switch active database to the requested one.
	log.Printf("Switching to database %s", dbName)
	return nil
}

func (d DdevHandler) HandleQuery(query string) (*mysql.Result, error) {
	// if strings.ToLower(query) == "show databases" {
	// 	var dbnames []string
	// 	for k := range servers {
	// 		dbnames = append(dbnames, k)
	// 	}

	// 	rs, err := mysql.BuildSimpleTextResultset(dbnames, dbnames)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("something went wrong")
	// 	}

	// 	r := mysql.Result{
	// 		Status:       mysql.SERVER_STATUS_AUTOCOMMIT,
	// 		AffectedRows: uint64(len(dbnames)),
	// 		Resultset:    rs,
	// 	}
	// 	return &r, nil
	// }

	// @todo: show databases should list all databases found at startup
	// @todo: all other queries should be relayed to the active database

	return nil, fmt.Errorf("not supported now")
}
