package concurrentselect

import "os"

type Options struct {
	Username string `short:"u" name:"username" help:"Stardog username. To be used instead in conjunction with password instead of token." env:"SD_USERNAME"`
	Password string `short:"p" name:"password" help:"Stardog password. To be used instead in conjunction with username instead of token." env:"SD_PASSWORD"`
	JWTToken string `name:"token" help:"JWT token to use for authentication. To be used instead of username and password." env:"SD_TOKEN"`

	Database string `short:"d" name:"database" required:"" help:"Database to execute queries against." env:"SD_DATABASE"`
	Server   string `short:"s" name:"server" required:"" help:"URL of the Stardog server." env:"SD_SERVER"`

	Query     string   `short:"q" name:"query" required:"" help:"Name of stored query or query string to execute." group:"Query Flags" xor:"Query Flags"`
	QueryFile *os.File `short:"f" name:"file" required:"" help:"File containing the query to execute." group:"Query Flags" xor:"Query Flags"`
	Timeout   int      `short:"t" name:"timeout" help:"Timeout in milliseconds for the provided query."`
	Reasoning bool     `short:"r" name:"reasoning" default:"false" help:"Enable reasoning for the provided query."`

	NumConcurrentQueries int    `short:"c" name:"concurrent-queries" required:"" help:"Number of concurrent queries to execute."`
	Format               string `name:"format" enum:"table,csv" default:"table" help:"Result format of benchmark results. Valid formats are 'table' or 'csv'."`
	PrintQuery           bool   `name:"print-query" help:"Print the query executed before the results."`
}
