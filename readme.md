# statscoll [![Build Status](https://travis-ci.org/frozzare/statscoll.svg?branch=master)](https://travis-ci.org/frozzare/statscoll) [![Go Report Card](https://goreportcard.com/badge/github.com/frozzare/statscoll)](https://goreportcard.com/report/github.com/frozzare/statscoll)

> Work in progress!

Collect numerics stats via http. It's not amed to be complex or advanced.

## Installation

```
go get -u github.com/frozzare/statscoll
```

## Usage

Create a config file `config.yml`

The `dsn` value is the data source name used to connect to the mysql database. Read more about [dsn](https://github.com/go-sql-driver/mysql#dsn-data-source-name).

This is the default values.

```yaml
port: 9300
dsn: root@/statscoll?charset=utf8&parseTime=true
```

Then run it:

```
statscoll [-c config.yml]
```

Then you can add stats via http:

```json
POST /collect
{
    "metric": "metric-name",
    "value": 10.0
}
```

Additional properties:
- custom `timestamp`
- project name `project`

Then you can list stats value:

```json
GET /stats/metric-name
[
    {
        "metric": "metric-name",
        "value": 10.0
    }
]
```

Possible query strings are:
- `start` to filter stat values that starts with timestamp value
- `end` to filter stat values that ends with timestamp value
- `project` to filter projects with same metric names as other projects.

Get total of value with the same query strings as for stats endpoint:

```json
GET /total/metric-name
{
    "total": 10.0
}
```

## License

MIT © [Fredrik Forsmo](https://github.com/frozzare)