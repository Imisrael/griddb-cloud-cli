We have already written a [Quick Start Guide](https://griddb.net/en/blog/griddb-cloud-quick-start-guide/) on how to use GridDB Cloud. And though we believe it's simple enough to get started using GridDB Cloud's WebAPI, we wanted to make some of the simple commands usuable from the CLI without having to make CURL requests which include your authentication headers in every command. Enter the GridDB Cloud CLI tool: [GitHub](https://github.com/Imisrael/griddb-cloud-cli)

The GridDB Cloud CLI Tool aims to make managing your GridDB Cloud database a little more manageable from the comfort of your own terminal! Tasks like querying, pushing data, creatinvg containers, etc can all be accomplished now in your CLI with the help of this tool. In this article, we will walkthrough how to install and use this tool and show some examples of what you can accomplish with it. 

## Getting Started (Download & Configuration)

The CLI Tool is distributed via [github as a single binary file](https://github.com/Imisrael/griddb-cloud-cli/releases/tag/v0.1.2). In the release section, you can download the appropriate version for your machine. Once downloaded, you can insert it in a directory in your PATH for your CLI and use from anywhere in the CLI, or alternatively, you can simply use the binary file from within the location it's located (ie. `./griddb-cloud-cli help`).

The tool is written in Go, so you could also clone the repo and build your own binary:

```bash
$ go get
$ go build
```

### Configuration

This tool expects a `.yaml` file to exist in `$HOME/.griddb.yaml` with the following fields: 

```bash
cloud_url: "url"
cloud_username: "example"
cloud_pass: "pass"
```

Alternatively, you save the file elsewhere and include the `--config` flag when running your tool (ie. `griddb-cloud-cli --config /opt/configs/griddb.yaml checkConnection`).

You will also still need to whitelist your IP Address in the GridDB Cloud Portal. Unfortunately this is not something that is achievable through the CLI Tool at this time.


## Features & Commands

This tool was written with the help of the ever-popular [Cobra](https://github.com/spf13/cobra) Library. Because of this, we are able to use the `--help` flag for all the commands in case you forget the functionality of some of the commands and their flags. 

```bash
$ griddb-cloud-cli help

A series of commands to help you manage your cloud-based DB.
Standouts include creating a container and graphing one using 'read graph' and 'create' respectfully

Usage:
  griddb-cloud-cli [command]

Available Commands:
  checkConnection Test your Connection with GridDB Cloud
  completion      Generate the autocompletion script for the specified shell
  create          Interactive walkthrough to create a container
  delete          Test your Connection with GridDB Cloud
  help            Help about any command
  ingest          Ingest a `csv` file to a new or existing container
  list            Get a list of all of the containers
  put             Interactive walkthrough to push a row
  read            Query container with TQL
  show            get container info 
  sql             Run a sql command

Flags:
      --config string   config file (default is $HOME/.griddb.yaml)
  -h, --help            help for griddb-cloud-cli
```

So with that out of the way, let's begin with the commands.

### All GridDB CLI Tool Commands

On your first time around, you should run the `checkConnection` command as a sanity check to ensure that you can connect to your instance. The tool will tell you if you have improper auth or if you're blocked by the firewall: 

#### Check Connection

```bash
$ griddb-cloud-cli checkConnection

[10005:TXN_AUTH_FAILED] (address=172.25.23.68:10001, partitionId=0)
2025/04/30 08:32:33 Authentication Error. Please check your username and password in your config file 
```

```bash
$ griddb-cloud-cli checkConnection

2025/04/30 08:33:48 (403) IP Connection Error. Is this IP Address Whitelisted? Please consider whitelisting Ip Address: X.X.X.116
```

```bash
$ griddb-cloud-cli checkConnection
2025/04/30 08:35:20 Please set a config file with the --config flag or set one in the default location $HOME/.griddb.yaml
```

And if everything is settled correctly:

```bash
$ griddb-cloud-cli checkConnection
200 OK
```

#### List Containers

You can list all containers inside of your Cloud DB Instance:

```bash
$ griddb-cloud-cli list

0: actual_reading_1
1: actual_reading_10
2: boiler_control_10
3: device1
4: device2
5: device3
6: device4
7: device6
```

#### Show Container

You can display the schema and other info about an individual container:


```bash
$ griddb-cloud-cli show device2  

{
    "container_name": "device2",
    "container_type": "TIME_SERIES",
    "rowkey": true,
    "columns": [
        {
            "name": "ts",
            "type": "TIMESTAMP",
            "timePrecision": "MILLISECOND",
            "index": []
        },
        {
            "name": "device",
            "type": "STRING",
            "index": []
        },
        {
            "name": "co",
            "type": "DOUBLE",
            "index": []
        },
        {
            "name": "humidity",
            "type": "FLOAT",
            "index": []
        },
        {
            "name": "light",
            "type": "BOOL",
            "index": []
        },
        {
            "name": "lpg",
            "type": "DOUBLE",
            "index": []
        },
        {
            "name": "motion",
            "type": "BOOL",
            "index": []
        },
        {
            "name": "smoke",
            "type": "DOUBLE",
            "index": []
        },
        {
            "name": "temperature",
            "type": "DOUBLE",
            "index": []
        }
    ]
}
```

#### Querying/Reading a Container

You can run TQL or SQL queries on your containers. TQL is the simpler option: 

```bash
$ griddb-cloud-cli read device2 --limit 1 --pretty

    [ { "name": "device2", "stmt": "select * limit 1", "columns": null, "hasPartialExecution": true }]
[
  [
    {
      "Name": "ts",
      "Type": "TIMESTAMP",
      "Value": "2006-01-02T07:04:05.700Z"
    },
    {
      "Name": "device",
      "Type": "STRING",
      "Value": "b8:27:eb:bf:9d:51"
    },
    {
      "Name": "co",
      "Type": "DOUBLE",
      "Value": 0.004955938648391245
    },
    {
      "Name": "humidity",
      "Type": "FLOAT",
      "Value": 51
    },
    {
      "Name": "light",
      "Type": "BOOL",
      "Value": false
    },
    {
      "Name": "lpg",
      "Type": "DOUBLE",
      "Value": 0.00765082227055719
    },
    {
      "Name": "motion",
      "Type": "BOOL",
      "Value": false
    },
    {
      "Name": "smoke",
      "Type": "DOUBLE",
      "Value": 0.02041127012241292
    },
    {
      "Name": "temperature",
      "Type": "DOUBLE",
      "Value": 22.7
    }
  ]
]
```

The `read` command will run a simple TQL query of your container which you can then specify the folliwng: an offset (--offset), a limit (-l, --limit), pretty print(-p, --pretty), just rows (--rows), which columns you want to see (--columns) or just the straight obj delivered from GridDB Cloud (--raw). 

Normally when you query a container with GridDB Cloud, it will send your results as two arrays, one with your column object, and another with more arrays of just row data. You can query this with --raw, but the default is to make a JSON and send that unstructured. If you use Pretty like above, it will indent and space it out for you. Just printing rows is better if you querying lots of rows: 

```bash
$ griddb-cloud-cli read device1 --limit 25 --rows

[ { "name": "device1", "stmt": "select * limit 25", "columns": null, "hasPartialExecution": true }]

ts,co,humidity,light,lpg,motion,smoke,temp,
[2020-07-12T01:00:25.984Z 0.0041795988 77.5999984741 true 0.006763671 false 0.0178934842 26.8999996185]
[2020-07-12T01:00:53.485Z 0.0048128545 53.5 false 0.0074903843 false 0.0199543908 21.7]
[2020-07-12T01:01:35.020Z 0.0030488793 74.9000015259 true 0.0053836916 false 0.014022829 19.5]
[2020-07-12T01:01:52.751Z 0.0049817187 51.3 false 0.0076795919 false 0.020493267 22.4]
[2020-07-12T01:01:59.191Z 0.003937408 72.9000015259 true 0.006477819 false 0.0170868731 24.7999992371]
[2020-07-12T01:02:01.157Z 0.0050077601 51.1 false 0.0077086115 false 0.0205759974 22.4]
[2020-07-12T01:02:01.445Z 0.0030841269 74.8000030518 true 0.0054286446 false 0.0141479363 19.6000003815]
[2020-07-12T01:02:04.938Z 0.0048169262 53.5 false 0.0074949679 false 0.0199674343 21.7]
[2020-07-12T01:02:05.182Z 0.0025840714 75.5999984741 false 0.0047765452 false 0.0123403139 19.6000003815]
[2020-07-12T01:02:12.428Z 0.0030488793 74.9000015259 true 0.0053836916 false 0.014022829 19.6000003815]
[2020-07-12T01:02:16.506Z 0.0048277855 53.5 false 0.0075071874 false 0.0200022097 21.7]
[2020-07-12T01:02:19.376Z 0.0030401715 74.9000015259 true 0.005372564 false 0.0139918711 19.6000003815]
[2020-07-12T01:02:21.754Z 0.0041428371 77.5999984741 true 0.0067205832 false 0.0177717486 26.8999996185]
[2020-07-12T01:02:29.017Z 0.0048400659 53.5 false 0.0075209965 false 0.0200415141 21.7]
[2020-07-12T01:02:33.443Z 0.0042300404 77.5999984741 true 0.0068226226 false 0.0180601254 26.7999992371]
[2020-07-12T01:02:35.686Z 0.00255591 75.5999984741 false 0.0047388314 false 0.0122362642 19.6000003815]
[2020-07-12T01:02:41.697Z 0.0030488793 75 true 0.0053836916 false 0.014022829 19.6000003815]
[2020-07-12T01:03:03.206Z 0.0042019019 77.5999984741 true 0.006789761 false 0.0179672218 26.7999992371]
[2020-07-12T01:03:04.701Z 0.0049946711 51.3 false 0.0076940309 false 0.0205344276 22.5]
[2020-07-12T01:03:04.768Z 0.0040601528 72.6999969482 true 0.0066232815 false 0.0174970393 24.7999992371]
[2020-07-12T01:03:05.999Z 0.0040886168 77.5 true 0.0066568388 false 0.0175917499 26.7999992371]
[2020-07-12T01:03:08.403Z 0.0048101357 53.7 false 0.0074873232 false 0.0199456799 21.8]
[2020-07-12T01:03:08.942Z 0.0049860142 51.1 false 0.0076843815 false 0.02050692 22.4]
[2020-07-12T01:03:10.023Z 0.0048141805 53.5 false 0.0074918772 false 0.0199586389 21.7]
[2020-07-12T01:03:12.863Z 0.0050019251 51.1 false 0.0077021129 false 0.020557469 22.3]
```

#### Querying Number Data into an ASCII Line Chart

Using a subcommand of `read`, you can also run a TQL query and read the results into a graph. For example: 

```bash
$ griddb-cloud-cli read graph device1 -l 10 --columns temp,humidity

[ { "name": "device1", "stmt": "select * limit 10", "columns": ["temp","humidity"], "hasPartialExecution": true }]
 77.60 ┼╮
 75.66 ┤╰╮                   ╭╮                                          ╭╮                    ╭───────────
 73.73 ┤ ╰╮                 ╭╯╰╮                   ╭╮                   ╭╯╰╮                  ╭╯
 71.79 ┤  ╰╮               ╭╯  │                  ╭╯╰╮                  │  ╰╮                ╭╯
 69.85 ┤   ╰╮             ╭╯   ╰╮                ╭╯  ╰╮                ╭╯   ╰╮              ╭╯
 67.92 ┤    │            ╭╯     ╰╮              ╭╯    ╰╮              ╭╯     ╰╮            ╭╯
 65.98 ┤    ╰╮          ╭╯       ╰╮            ╭╯      ╰╮            ╭╯       ╰╮          ╭╯
 64.04 ┤     ╰╮        ╭╯         ╰╮          ╭╯        ╰╮          ╭╯         ╰╮        ╭╯
 62.11 ┤      ╰╮      ╭╯           ╰╮        ╭╯          ╰╮        ╭╯           ╰╮      ╭╯
 60.17 ┤       ╰╮    ╭╯             ╰╮      ╭╯            │       ╭╯             ╰╮    ╭╯
 58.23 ┤        ╰╮  ╭╯               ╰╮    ╭╯             ╰╮     ╭╯               ╰╮  ╭╯
 56.30 ┤         ╰╮╭╯                 ╰╮  ╭╯               ╰╮   ╭╯                 ╰╮╭╯
 54.36 ┤          ╰╯                   ╰╮╭╯                 ╰╮  │                   ╰╯
 52.42 ┤                                ││                   ╰╮╭╯
 50.49 ┤                                ╰╯                    ╰╯
 48.55 ┤
 46.61 ┤
 44.68 ┤
 42.74 ┤
 40.80 ┤
 38.87 ┤
 36.93 ┤
 34.99 ┤
 33.06 ┤
 31.12 ┤
 29.18 ┤
 27.25 ┼─╮
 25.31 ┤ ╰───╮                                   ╭────╮
 23.37 ┤     ╰───╮                      ╭────────╯    ╰────────╮
 21.44 ┤         ╰───────╮       ╭──────╯                      ╰───────╮     ╭──────────────╮
 19.50 ┤                 ╰───────╯                                     ╰─────╯              ╰──────────────
                                          Col names from container device1

                                                ■ temp   ■ humidity
```

The results are color-coded so that you can accuately see which cols are mapped to which values. It also automatically omits non-number types if you just want to read the entire container a line chart: 

```bash
$ griddb-cloud-cli read graph device1 -l 5

Column ts (of type TIMESTAMP ) is not a `number` type. Omitting
Column light (of type BOOL ) is not a `number` type. Omitting
Column motion (of type BOOL ) is not a `number` type. Omitting
 77.60 ┼─╮
 75.01 ┤ ╰─╮                                            ╭─╮
 72.43 ┤   ╰──╮                                      ╭──╯ ╰──╮                                          ╭──
 69.84 ┤      ╰──╮                                ╭──╯       ╰──╮                                     ╭─╯
 67.25 ┤         ╰─╮                           ╭──╯             ╰─╮                                ╭──╯
 64.67 ┤           ╰──╮                     ╭──╯                  ╰──╮                          ╭──╯
 62.08 ┤              ╰──╮               ╭──╯                        ╰──╮                    ╭──╯
 59.49 ┤                 ╰─╮          ╭──╯                              ╰──╮              ╭──╯
 56.91 ┤                   ╰──╮    ╭──╯                                    ╰─╮         ╭──╯
 54.32 ┤                      ╰────╯                                         ╰──╮   ╭──╯
 51.73 ┤                                                                        ╰───╯
 49.15 ┤
 46.56 ┤
 43.97 ┤
 41.39 ┤
 38.80 ┤
 36.21 ┤
 33.63 ┤
 31.04 ┤
 28.46 ┤
 25.87 ┼───────────╮                                                                                    ╭──
 23.28 ┤           ╰───────────╮                                              ╭─────────────────────────╯
 20.70 ┤                       ╰──────────────────────────────────────────────╯
 18.11 ┤
 15.52 ┤
 12.94 ┤
 10.35 ┤
  7.76 ┤
  5.18 ┤
  2.59 ┤
  0.00 ┼───────────────────────────────────────────────────────────────────────────────────────────────────
                                          Col names from container device1

                                    ■ co   ■ humidity   ■ lpg   ■ smoke   ■ temp
```

#### Creating Containers

You can create containers using an interactive question prompt in the CLI. It will ask for container name, container type, rowkey, and col names and types. 

For example, let's create a new time series container with two columns: 

```bash
$ griddb-cloud-cli create


✔ Container Name: … sample1
✔ Choose: … TIME_SERIES
✔ How Many Columns for this Container? … 2
✔ Col name For col #1 … ts
✔ Col #1(TIMESTAMP CONTAINERS ARE LOCKED TO TIMESTAMP FOR THEIR ROWKEY) … TIMESTAMP
✔ Col name For col #2 … temp
✔ Column Type for col #2 … DOUBLE
✔ Make Container? 
{
    "container_name": "sample1",
    "container_type": "TIME_SERIES",
    "rowkey": true,
    "columns": [
        {
            "name": "ts",
            "type": "TIMESTAMP",
            "index": null
        },
        {
            "name": "temp",
            "type": "DOUBLE",
            "index": null
        }
    ]
} … YES
{"container_name":"sample1","container_type":"TIME_SERIES","rowkey":true,"columns":[{"name":"ts","type":"TIMESTAMP","index":null},{"name":"temp","type":"DOUBLE","index":null}]}
201 Created
```

If you can't easily follow along with the prompt here, please just download the tool and try it for yourself!

And note, as explained in the prompts, if you select to create a TIME_SERIES Container, the rowkey is auto set to true and the first col must have a type of TIMESTAMP. Collection containers have different ruls.

#### Putting Rows to Containers

Similarly, you can follow along with the prompt to push data into your container, 1 by 1. Here we will push to our new container `sample1` and use `NOW()` as our current timestamp: 

```bash 
$ griddb-cloud-cli put sample1

Container Name: sample1
✔ Column 1 of 2
 Column Name: ts
 Column Type: TIMESTAMP … NOW()
✔ Column 2 of 2
 Column Name: temp
 Column Type: DOUBLE … 20.2
[["2025-04-30T07:43:03.700Z",  20.2]]
✔ Add the Following to container sample1? … YES
200 OK
```

#### Ingesting CSV Data

You can also ingest full CSV files with this tool. It too uses an interactive prompt as there is information that needs to be set for each col, such as index position in csv and data type. Once you set those, it will ingest the data in chunks of 1000.

