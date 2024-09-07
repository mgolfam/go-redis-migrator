#!/bin/bash

# Switch to the desired database
redis-cli -h localhost -p 6379 select 5

# Insert 100 random keys into the Redis database
for i in $(seq 1 100); do
  redis-cli -h localhost -p 6379 -n 1 set key$i "value$i"
done
