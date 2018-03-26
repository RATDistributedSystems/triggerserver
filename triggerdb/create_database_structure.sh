#!/bin/bash
(sleep 65 && cqlsh -f create_trigger_structure.cql) &
cd / && docker-entrypoint.sh