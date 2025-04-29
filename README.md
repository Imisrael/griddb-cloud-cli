A simple CLI tool wrapper for making HTTP Requests to your GridDB Cloud instance.

## Getting Started
To start, first gather your GridDB Cloud credentials and stick them in $HOME/.griddb.yaml (or, you can simply use the `--config` flag and point to your file when using the cli tool)

## Examples



$ ./griddb-cloud-cli checkConnection

    200 OK



$ ./griddb-cloud-cli read -p device1 --limit 1

    [ { "name": "device2", "stmt": "select * limit 2", "columns": null, "hasPartialExecution": true }]

[

  [

    {

      "Name": "ts",

      "Type": "STRING",

      "Value": "1.5945120943859746E9"

    },

    {

      "Name": "device",

      "Type": "STRING",

      "Value": "b8:27:eb:bf:9d:51"

    },

    {

      "Name": "co",

      "Type": "FLOAT",

      "Value": 0.0049559386

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

      "Type": "FLOAT",

      "Value": 0.0076508224

    },

    {

      "Name": "motion",

      "Type": "BOOL",

      "Value": false

    },

    {

      "Name": "smoke",

      "Type": "FLOAT",

      "Value": 0.02041127

    },

    {

      "Name": "temperature",

      "Type": "FLOAT",

      "Value": 22.7

    }

  ]

]



$ ./griddb-cloud-cli read graph device2 -l 10



[ { "name": "device2", "stmt": "select * limit 10", "columns": null, "hasPartialExecution": true }]

Column ts (of type STRING ) is not a `number` type. Omitting

Column device (of type STRING ) is not a `number` type. Omitting

Column light (of type BOOL ) is not a `number` type. Omitting

Column motion (of type BOOL ) is not a `number` type. Omitting

 77.90 ┤                                                             ╭╮                           ╭────────

 75.30 ┤           ╭─╮                     ╭──╮                     ╭╯╰╮                     ╭────╯

 72.71 ┤          ╭╯ ╰╮                   ╭╯  ╰╮                  ╭─╯  ╰╮                   ╭╯

 70.11 ┤        ╭─╯   ╰╮                 ╭╯    ╰╮                ╭╯     ╰─╮                ╭╯

 67.51 ┤       ╭╯      ╰─╮              ╭╯      ╰─╮             ╭╯        ╰╮             ╭─╯

 64.92 ┤      ╭╯         ╰╮           ╭─╯         ╰╮           ╭╯          ╰╮           ╭╯

 62.32 ┤    ╭─╯           ╰╮         ╭╯            ╰╮         ╭╯            ╰╮         ╭╯

 59.72 ┤   ╭╯              ╰─╮      ╭╯              ╰╮      ╭─╯              ╰╮      ╭─╯

 57.13 ┤  ╭╯                 ╰╮    ╭╯                ╰─╮   ╭╯                 ╰╮    ╭╯

 54.53 ┤ ╭╯                   ╰╮ ╭─╯                   ╰╮ ╭╯                   ╰─╮ ╭╯

 51.93 ┼─╯                     ╰─╯                      ╰─╯                      ╰─╯

 49.34 ┤

 46.74 ┤

 44.14 ┤

 41.55 ┤

 38.95 ┤

 36.35 ┤

 33.76 ┤

 31.16 ┤

 28.57 ┤

 25.97 ┤                              ╭────────────╮           ╭────────────╮                          ╭───

 23.37 ┼──╮                   ╭───────╯            ╰───────────╯            ╰───────╮             ╭────╯

 20.78 ┤  ╰───────────────────╯                                                     ╰─────────────╯

 18.18 ┤

 15.58 ┤

 12.99 ┤

 10.39 ┤

  7.79 ┤

  5.20 ┤

  2.60 ┤

  0.00 ┼───────────────────────────────────────────────────────────────────────────────────────────────────

                                          Col names from container device2



                                ■ co   ■ humidity   ■ lpg   ■ smoke   ■ temperature



\# Interactive mode with create and ingest

$ ./griddb-cloud-cli create

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



$ ./griddb-cloud-cli ingest iot_telemetry_data.csv

✔ Does this container already exist? … NO

Use CSV Header names as your GridDB Container Col names? 

ts,device,co,humidity,light,lpg,motion,smoke,temp

✔ Y/n … YES

✔ Container Name: … device6

✔ Choose: … TIME_SERIES

✔ Col ts(TIMESTAMP CONTAINERS ARE LOCKED TO TIMESTAMP FOR THEIR ROWKEY) … TIMESTAMP

✔ (device) Column Type … STRING

✔ (co) Column Type … DOUBLE

✔ (humidity) Column Type … DOUBLE

✔ (light) Column Type … BOOL

✔ (lpg) Column Type … DOUBLE

✔ (motion) Column Type … BOOL

✔ (smoke) Column Type … DOUBLE

✔ (temp) Column Type … DOUBLE

        },

        {

            "name": "device",

            "type": "STRING",

            "index": null

        },

        {

            "name": "co",

            "type": "DOUBLE",

            "index": null

        },

        {

            "name": "humidity",

            "type": "DOUBLE",

            "index": null

        },

        {

            "name": "light",

            "type": "BOOL",

            "index": null

        },

        {

            "name": "lpg",

            "type": "DOUBLE",

            "index": null

        },

        {

            "name": "motion",

            "type": "BOOL",

            "index": null

        },

        {

            "name": "smoke",

            "type": "DOUBLE",

            "index": null

        },

        {

            "name": "temp",

            "type": "DOUBLE",

            "index": null

        }

    ]

} … YES

{"container_name":"device6","container_type":"TIME_SERIES","rowkey":true,"columns":[{"name":"ts","type":"TIMESTAMP","index":null},{"name":"device","type":"STRING","index":null},{"name":"co","type":"DOUBLE","index":null},{"name":"humidity","type":"DOUBLE","index":null},{"name":"light","type":"BOOL","index":null},{"name":"lpg","type":"DOUBLE","index":null},{"name":"motion","type":"BOOL","index":null},{"name":"smoke","type":"DOUBLE","index":null},{"name":"temp","type":"DOUBLE","index":null}]}

201 Created

Container Created. Starting Ingest

0 ts ts

1 device device

2 co co

3 humidity humidity

4 light light

5 lpg lpg

6 motion motion

7 smoke smoke

8 temp temp

✔ Is the above mapping correct? … YES

Ingesting. Please wait...

Inserting 1000 rows

200 OK

Inserting 1000 rows

200 OK

Inserting 1000 rows

