# sd-utils

`sd-utils` is a command line application containing a variety of tools for working with Stardog.

## Installation

### Build from source

Clone this repo, build from source with `cd sd-utils && go build`, then move the binary to somewhere accessible in your `PATH`, e.g. `mv ./sd-utils /usr/local/bin`.

## `bench-select`

Execute a SPARQL select query concurrently to benchmark server response times.

```text
Usage: sd-utils bench-select --database=STRING --server=STRING --query=STRING --file=FILE --concurrent-queries=INT

execute a SPARQL select query concurrently to benchmark server response times

Flags:
  -h, --help                      Show context-sensitive help.

  -u, --username=STRING           Stardog username. To be used instead in conjunction with password instead of token ($SD_USERNAME).
  -p, --password=STRING           Stardog password. To be used instead in conjunction with username instead of token ($SD_PASSWORD).
      --token=STRING              JWT token to use for authentication. To be used instead of username and password ($SD_TOKEN).
  -d, --database=STRING           Database to execute queries against ($SD_DATABASE).
  -s, --server=STRING             URL of the Stardog server ($SD_SERVER).
  -t, --timeout=INT               Timeout in milliseconds for the provided query.
  -r, --reasoning                 Enable reasoning for the provided query.
  -c, --concurrent-queries=INT    Number of concurrent queries to execute.
      --format="table"            Result format of benchmark results. Valid formats are 'table' or 'csv'

Query Flags
  -q, --query=STRING    Name of stored query or query string to execute.
  -f, --file=FILE       File containing the query to execute.
```

```bash
# environment variables can be used to set some defaults
export SD_USERNAME=admin
export SD_PASSWORD=admin
export SD_SERVER=http://localhost:5820
export SD_DATABASE=myDatabase

sd-utils bench-select -c 5 -q "SELECT * { ?s a ?p } LIMIT 10"
```

<img src="https://vhs.charm.sh/vhs-2CWIY1n4CAhe0Q1saNUCTV.gif" alt="Made with VHS">
<a href="https://vhs.charm.sh">
  <img src="https://stuff.charm.sh/vhs/badge.svg">
</a>


