A simple CLI tool wrapper for making HTTP Requests to your GridDB Cloud instance.

## Getting Started

To start, first gather your GridDB Cloud credentials and stick them in $HOME/.griddb.yaml (or, you can simply use the `--config` flag and point to your file when using the cli tool) Required fields:

```bash
cloud_url: "url"
cloud_username: "example"
cloud_pass: "pass"
```

## Examples

$ ./griddb-cloud-cli checkConnection

```bash
    200 OK
```

$ ./griddb-cloud-cli list 

```bash
0: actual_reading_1
1: actual_reading_10
2: boiler_control_10
3: device1
4: device2
5: device3
6: device4
7: device6
```

./griddb-cloud-cli show device2        

```bash
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

$ ./griddb-cloud-cli read device1 --limit 1 --pretty

```bash
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
```

$ ./griddb-cloud-cli read device2 --limit 9 --rows

```bash
[ { "name": "device2", "stmt": "select * limit 9", "columns": null, "hasPartialExecution": true }]

ts,device,co,humidity,light,lpg,motion,smoke,temperature,
[1.5945120943859746E9 b8:27:eb:bf:9d:51 0.0049559386 51 false 0.0076508224 false 0.02041127 22.7]
[1.5945120947355676E9 00:0f:00:70:91:0a 0.0028400887 76 false 0.0051143835 false 0.013274836 19.7]
[1.5945120980735729E9 b8:27:eb:bf:9d:51 0.0049760123 50.9 false 0.007673227 false 0.020475125 22.6]
[1.594512099589146E9 1c:bf:ce:15:ec:4d 0.004403027 76.8 true 0.0070233373 false 0.018628225 27]
[1.594512101761235E9 b8:27:eb:bf:9d:51 0.0049673636 50.9 false 0.0076635773 false 0.020447621 22.6]
[1.5945121044684107E9 1c:bf:ce:15:ec:4d 0.004391004 77.9 true 0.0070094587 false 0.018588908 27]
[1.5945121054488637E9 b8:27:eb:bf:9d:51 0.0049760253 50.9 false 0.0076732417 false 0.020475166 22.6]
[1.594512106869076E9 00:0f:00:70:91:0a 0.0029381157 76 false 0.005241482 false 0.013627521 19.7]
[1.5945121082753816E9 1c:bf:ce:15:ec:4d 0.0043454715 77.9 true 0.006956802 false 0.018439783 27]
```



$ ./griddb-cloud-cli read graph device2 -l 10

```bash
[ { "name": "device2", "stmt": "select * limit 10", "columns": null, "hasPartialExecution": true }]

Column ts (of type TIMESTAMP ) is not a `number` type. Omitting
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
```

\# Interactive mode with create and ingest

$ ./griddb-cloud-cli create

```bash
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

$ ./griddb-cloud-cli ingest iot_telemetry_data.csv

```bash
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
```