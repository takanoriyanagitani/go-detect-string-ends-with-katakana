#!/bin/zsh

wpath=./opt.wasm
wpath=./detect.wasm

ibench=./sample.d/heavy.dat

bin=./cmd/grep-string-ends-with-katakana-special/grep-string-ends-with-katakana-special

input1(){
  echo e382a8 e3838a e382b8 | xxd -r -ps
}

input2(){
  echo e382a8 e3838a e382b8 e383bc | xxd -r -ps
}

input3(){
  echo e382a8 e3838a | xxd -r -ps
}

input4(){
  echo e382a8 e382a8 e3838a e382b8 | xxd -r -ps
}

check_wasm(){
    $bin -wasm-path "${wpath}"
}

check_grep(){
    LC_ALL=C \
      grep \
        -E \
        $'\xE3(\x82[\xA2-\xBF]|\x83[\x80-\xBA])\xE3(\x82[\xA1-\xBF]|\x83[\x80-\xBC])\xE3(\x82[\xA1-\xBF]|\x83[\x80-\xBA])$'
}

compare(){
  echo check using the wasm
  input1 | check_wasm
  input2 | check_wasm
  input3 | check_wasm
  input4 | check_wasm
  echo
  
  echo check using grep
  input1 | check_grep
  input2 | check_grep
  input3 | check_grep
  input4 | check_grep
}

bench_input(){
  local size
  local size_default
  size_default=$(( 0xfffffff ))

  size=${1:-${size_default}}
  openssl rand $size
}

bench_grep(){
  time bench_input $1 |
    dd bs=1048576 status=progress |
    LC_ALL=C \
      grep \
        -E \
        $'\xE3(\x82[\xA2-\xBF]|\x83[\x80-\xBA])\xE3(\x82[\xA1-\xBF]|\x83[\x80-\xBC])\xE3(\x82[\xA1-\xBF]|\x83[\x80-\xBA])$' |
        wc -c
}

bench_wasm(){
  time bench_input $1 |
    dd bs=1048576 status=progress |
    $bin -wasm-path "${wpath}" |
        wc -c
}

#bench_grep $((  0xfffffff ))
bench_wasm $((  0xfffffff ))
#bench_wasm $(( 0x37fffffff ))
