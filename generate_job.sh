#!/bin/bash

for i in {1..500}
do
  colonies process submit --spec fibonacci_process_spec.json
done
