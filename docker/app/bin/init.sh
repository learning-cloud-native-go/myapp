#!/bin/sh
echo 'Runing migrations...'
/myapp/bin/migrate up > /dev/null 2>&1 &

echo 'Start application...'
/myapp/bin/app