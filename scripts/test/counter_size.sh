#!/usr/bin/env bash
export PATH=`pwd`/bin:${PATH}
git clone https://github.com/xuperchain/contract-sdk-cpp.git|| true
cd contract-sdk-cpp
xdev build -o counter1.wasm  example/counter.cc
XDEV_ROOT=`pwd` xdev build --using-precompiled-sdk=false -o counter2.wasm example/counter.cc

# for debug
ls -alh
counter1_size=`du -k "counter1.wasm" | cut -f1`
counter2_size=`du -k "counter2.wasm" | cut -f1`

echo "counter1" ${counter1_size}
echo "counter2" ${counter2_size}

if [ "$counter1_size" = "$counter2_size" ];then
  exit 0
fi
exit -1
