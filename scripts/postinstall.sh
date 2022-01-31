#!/bin/bash

chown -R 1000:1000 /opt/premiumizearrd/
systemctl enable premiumizearrd.service
systemctl daemon-reload
systemctl start premiumizearrd.service
