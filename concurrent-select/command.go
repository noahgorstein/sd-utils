package concurrentselect

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/noahgorstein/go-stardog/stardog"
	"github.com/rodaine/table"
)

// Run runs the bench-select command.
func (o Options) Run() error {
	defer o.QueryFile.Close()
	q, err := getQuery(o.QueryFile, o.Query)
	if err != nil {
		return err
	}

	client, err := getStardogClient(o.Username, o.Password, o.JWTToken, o.Server)
	if err != nil {
		return err
	}

	results := make(chan run)
	var wg sync.WaitGroup

	for i := 0; i < o.NumConcurrentQueries; i++ {
		wg.Add(1)
		go selectQuery(client, results, &wg, q, o.Database, o.Timeout, o.Reasoning)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	err = printResults(results, o.Format)
	if err != nil {
		return err
	}

	return nil
}

// getStardogClient returns an authenticated *stardog.Client with the provided token
// or username and password for the server.
func getStardogClient(username, password, token, server string) (*stardog.Client, error) {
	if token != "" {
		t := &stardog.BearerAuthTransport{
			BearerToken: token,
		}
		return stardog.NewClient(server, t.Client())
	}
	t := &stardog.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	return stardog.NewClient(server, t.Client())
}

// getQuery gets the query string to execute from the provided file
// containing the query or the query string passed as a CLI flag.
func getQuery(queryFile *os.File, query string) (string, error) {

	if query != "" {
		return query, nil
	}

	if queryFile != nil {
		fileInfo, err := queryFile.Stat()
		if err != nil {
			log.Fatal(err)
		}
		fileSize := fileInfo.Size()
		buffer := make([]byte, fileSize)

		_, err = queryFile.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		return string(buffer), nil
	}

	return "", errors.New("No query provided in a file or as a string.")

}

// run represents one API request executing a query against the Stardog server.
type run struct {
	// HTTP Status Code
	statusCode int

	// error from API request
	error error

	// time at which API request is made
	startTime time.Time

	// time taken to execute the API request
	elapsedTime time.Duration
}

// printResults prints results from the results chan in the specified format.
func printResults(results chan run, format string) error {
	if format == "table" {
		printTable(results)
		return nil
	}
	err := printCSV(results)
	if err != nil {
		return err
	}
	return nil
}

func printTable(results chan run) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	tbl := table.New("Start Time", "Elapsed Time", "Status Code", "Error")
	tbl.WithHeaderFormatter(headerFmt)
	for msg := range results {
		startTime := msg.startTime.Format(time.RFC3339Nano)
		elapsedTime := msg.elapsedTime.String()
		statusCode := msg.statusCode

		errorMsg := "-"
		if msg.error != nil {
			errorMsg = msg.error.Error()
		}

		tbl.AddRow(startTime, elapsedTime, statusCode, errorMsg)
	}
	tbl.Print()
}

func printCSV(results chan run) error {
	w := csv.NewWriter(os.Stdout)

	var b []byte
	buf := bytes.NewBuffer(b)

	records := [][]string{}
	header := []string{
		"Start Time",
		"Elapsed",
		"Status Code",
		"Error",
	}
	records = append(records, header)

	for msg := range results {
		startTime := msg.startTime.Format(time.RFC3339Nano)
		elapsedTime := msg.elapsedTime.String()
		statusCode := msg.statusCode

		errorMsg := "-"
		if msg.error != nil {
			errorMsg = msg.error.Error()
		}

		record := []string{
			startTime,
			elapsedTime,
			strconv.Itoa(statusCode),
			errorMsg,
		}
		records = append(records, record)
	}

	for _, r := range records {
		if err := w.Write(r); err != nil {
			return fmt.Errorf("error writing record to csv: %v", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	fmt.Println(string(buf.Bytes()))
	return nil
}

// TODO: refactor many params
func selectQuery(client *stardog.Client, results chan run, wg *sync.WaitGroup, query string, database string, timeout int, reasoning bool) {
	queryOpts := stardog.SelectOptions{
		ResultFormat: stardog.QueryResultFormatCSV,
	}

	if timeout != 0 {
		queryOpts.Timeout = timeout
	}

	if reasoning {
		queryOpts.Reasoning = true
	}

	startTime := time.Now()
	_, resp, err := client.Sparql.Select(context.Background(), database, query, &queryOpts)

	run := run{
		startTime:   startTime,
		elapsedTime: time.Since(startTime),
		statusCode:  resp.StatusCode,
	}

	// failure
	if err != nil {
		run.error = err
		results <- run
		wg.Done()
		return
	}

	//success
	results <- run
	wg.Done()
}
