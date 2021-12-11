#/bin/bash

echo "Redis batch set test"

before=`redis info memory|grep used_memory:|awk 'BEGIN{FS=":"}{print $2}'`
echo "Before Memory:$before"

total=$1
l=$2
input=`printf %0${l}d 1`
klen=${#total}

for ((i=1; i<=total; i++))
do

key=`printf test_%0${klen}d $i`

one_set=`redis set ${key} ${input};`

done

after=`redis info memory|grep used_memory:|awk 'BEGIN{FS=":"}{print $2}'`

echo "After Memory:$after"

add = $((10#${after//$'\r'}-${before//$'\r'}))

echo "Memory Add:$add"

keylen=$(($add/$total-$l))

echo "Avg key length:${keylen}"

echo "remove test data"

redis keys "test*" | xargs redis  del
