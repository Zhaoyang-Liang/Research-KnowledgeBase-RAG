## Overview
This is a proof-of-concept implementation for subring secret encapsulation in CKKS/BGV/BFV.
The CKKS part is implemented on Lattigo, while the BGV/BFV part is on HElib.

## Usage
### CKKS part
In `main/go.sum`, replace the path for `github.com/tuneinsight/lattigo/v6` with the path of this folder. Then run the following commands to build.
```shell
cd main
go build
```

The data in the paper can be reproduced by running the following commands in `main`
```shell
mkdir ../logs
go run main -sparse=1 -change=0 -repeat=20 &> ../logs/sparse.log
go run main -BMTH=1 -eps=-15 -repeat=20 &> ../logs/BMTH.log
go run main -guide=1 -eps=-38 -repeat=20 &> ../logs/guide.log
# ours
go run main -eps=-16 -scale=35 -degree=127 -r=3 -scalestc=33 -scalects=52 -repeat=20 &> ../logs/16_2.log
go run main -eps=-16 -scale=35 -degree=127 -r=3 -scalestc=33 -scalects=52 -repeat=20 -lvls=10 &> ../logs/16_1.log

go run main -eps=-38 -scale=35 -degree=127 -r=4 -scalestc=32 -scalects=50 -repeat=20 &> ../logs/38_2.log
go run main -eps=-38 -scale=35 -degree=127 -r=4 -scalestc=32 -scalects=50 -repeat=20 -lvls=10 &> ../logs/38_1.log

go run main -eps=-64 -scale=35 -degree=127 -r=4 -scalestc=31 -scalects=51 -repeat=20 &> ../logs/64_2.log
go run main -eps=-64 -scale=35 -degree=127 -r=4 -scalestc=31 -scalects=51 -repeat=20 -lvls=10 &> ../logs/64_1.log

go run main -eps=-128 -scale=35 -degree=255 -r=3 -scalestc=30 -scalects=52 -repeat=20 &> ../logs/128_2.log
go run main -eps=-128 -scale=35 -degree=255 -r=3 -scalestc=30 -scalects=52 -repeat=20 -lvls=10 &> ../logs/128_1.log
```

### BGV/BFV part
Our implementation is based on [Ma et al.'s implementation](https://github.com/msh086/bgv-bootstrapping-with-homomorphic-NTT).
After building and installing HElib with our updated patch `HElib_code/BGV_bts_subring.patch`, use the following commands in this folder to build the main program.
```shell
cd HElib_code
mkdir build
cmake -S . -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build
```
The data in the paper can be reproduced by running the following commands in `HElib_code`
```shell
mkdir logs
./build/fatboot newbts=1 extdeg=16 scale=6.3 &> ./logs/16.log
./build/fatboot newbts=1 extdeg=16 scale=8 &> ./logs/32.log
./build/fatboot newbts=1 extdeg=16 scale=10.3 &> ./logs/64.log
./build/fatboot newbts=1 extdeg=16 scale=14 &> ./logs/128.log

./build/fatboot newbts=1 extdeg=1 scale=6.3 &> ./logs/16_base.log
./build/fatboot newbts=1 extdeg=1 scale=8 &> ./logs/32_base.log
./build/fatboot newbts=1 extdeg=1 scale=10.3 &> ./logs/64_base.log
./build/fatboot newbts=1 extdeg=1 scale=14 &> ./logs/128_base.log

./build/fatboot newbts=1 extdeg=1 scale=8.8 h=32 &> ./logs/128_sp.log
```