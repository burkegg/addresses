package cmd

import (
	"github.com/burkegg/addresses/pkg/addresses"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var dbHost string
var dbPort int
var dbUser string
var dbPassword string
var dbName string
var address string
var port int


// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the server.",
	Long: ` 
Start the server.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if dbHost == "" || dbPassword == "" {
			log.Fatalf("Missing required environment variables. Expected database host and database password but got: \n database host: %s\n database password: %s\n", dbHost, dbPassword)
			return
		}
		err := addresses.RunServer(dbHost, dbPort, dbUser, dbPassword, dbName, address, port, VERSION)
		if err != nil {
			log.Fatalf("Server failed to run: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&dbHost, "databaseHost", "o", os.Getenv("DATABASE_HOST"), "database host")
	runCmd.Flags().IntVarP(&dbPort, "databasePort", "b", 5432, "database port")
	runCmd.Flags().StringVarP(&dbUser, "databaseUser", "u", "postgres", "database user login")
	runCmd.Flags().StringVarP(&dbPassword, "dbPassword", "d", os.Getenv("POSTGRES_PASSWORD"), "database password")
	runCmd.Flags().StringVarP(&dbName, "dbName", "n", "coding_challenge", "database name")
	runCmd.Flags().StringVarP(&address, "address", "a", "0.0.0.0", "Address to run upon")
	runCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server upon.")
}
