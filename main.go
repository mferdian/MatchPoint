package main

import (
	"fieldreserve/config/database"
	"fmt"
)

func main() {
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	fmt.Println("OK")
}
