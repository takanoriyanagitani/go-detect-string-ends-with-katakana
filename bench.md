# Simple Benchmark

## Text

| size   | tool      | user | sys | cpu | total | rate        | ratio |
|:------:|:---------:|:----:|:---:|:---:|:-----:|:-----------:|:-----:|
| 29 MiB | grep      | 0.7  | 0.0 | 99% | 0.8   |    38 MiB/s | (1.0) |
| 29 MiB | go/wazero | 0.1  | 0.0 | 98% | 0.1   |   216 MiB/s |  6x   |

## Random input

| size      | tool      | user | sys | cpu | total | rate        | ratio |
|:---------:|:---------:|:----:|:---:|:---:|:-----:|:-----------:|:-----:|
|   256 MiB | grep      | 2.4  | 0.0 | 99% | 2.5   |   102 MiB/s | (1.0) |
|   256 MiB | go/wazero | 0.0  | 0.1 | 45% | 0.2   | 1,280 MiB/s | 13x   |
| 2,048 MiB | go/wazero | 0.3  | 0.3 | 38% | 1.3   | 1,575 MiB/s | 15x   |
