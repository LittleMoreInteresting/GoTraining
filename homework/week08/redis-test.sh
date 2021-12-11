#/bin/bash

echo "Start redis-benchmark "

for i in 10 20 50 100 200 1024 5120
do
echo "----$i----"

redis-benchmark -t set,get -d $i --csv

echo "*** sleep 10s ***"
sleep 10s

done
