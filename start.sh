#/bin/bash
make

pid=`lsof -i tcp:3000|grep "*:"|awk '{print $2}'|uniq`

if [ $pid ]; then
  echo $pid
 kill -1 $pid
 else
  ./bin/fastApi_mac server &
fi


# lsof -i tcp:3000|grep "*:"|awk '{print $2}'|xargs kill -1
