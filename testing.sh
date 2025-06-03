#!/bin/bash

go install; sudo cp /home/israel/go/bin/griddb-cloud-cli /opt/fluentd_scripts; sudo systemctl restart fluentd.service