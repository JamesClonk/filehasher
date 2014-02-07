#!/bin/bash

echo " "

CHECK_DIR=
select target in $(ls -d /mnt/*); do
        CHECK_DIR=${target}/data
        break;
done
echo " "

TARGET_NUMBER=
select number in 0 1 2 3 4 5 6; do
        TARGET_NUMBER=$number
        break;
done
echo " "

rm -f listing_${TARGET_NUMBER}.txt
rm -f hash_${TARGET_NUMBER}.txt

echo "$CHECK_DIR"

find $CHECK_DIR -type f -print >> listing_${TARGET_NUMBER}.txt

while read line
do
        SUM=`md5sum "$line"`
        echo "$SUM"
        echo "$SUM" >> hash_${TARGET_NUMBER}.txt
done < listing_${TARGET_NUMBER}.txt
