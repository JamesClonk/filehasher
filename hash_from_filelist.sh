#!/bin/bash

echo " "

TARGET_NUMBER=
select number in 0 1 2 3 4 5 6; do
        TARGET_NUMBER=$number
        break;
done
echo " "

while read line
do
        SUM=`md5sum "$line"`
        echo "$SUM"
        echo "$SUM" >> hash_${TARGET_NUMBER}.txt
done < listing_${TARGET_NUMBER}.txt
