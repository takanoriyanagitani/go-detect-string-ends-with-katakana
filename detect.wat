(module

  (func $ends_with3zkana128v (param $iv v128) (result i32)
    ;; check 1: 0xE382**
    local.get $iv
    v128.const i32x4 0x00E382A2 0x00E382A1 0x00E382A1 0x00
    i32x4.ge_u

    local.get $iv
    v128.const i32x4 0x00E382BF 0x00E382BF 0x00E382BF 0xFFFFFFFF
    i32x4.le_u

    v128.and

    ;; check 2: 0xE383**
    local.get $iv
    v128.const i32x4 0x00E38380 0x00E38380 0x00E38380 0x00
    i32x4.ge_u

    local.get $iv
    v128.const i32x4 0x00E383BA 0x00E383BC 0x00E383BA 0xFFFFFFFF
    i32x4.le_u

    v128.and


    ;; check 1 or check 2
    v128.or

    i32x4.all_true
  )

  (func $ends_with3zkana32i (export "ends_with3zkana32i")
    (param $first i32)
    (param $mid i32)
    (param $last i32)
    (result i32)

    v128.const i64x2 0 0

    local.get $first
    i32x4.replace_lane 0

    local.get $mid
    i32x4.replace_lane 1

    local.get $last
    i32x4.replace_lane 2

    call $ends_with3zkana128v
  )

  (func $ends_with3zkana64i (export "ends_with3zkana64i")
    (param $first i64)
    (param $last i64)
    (result i32)

    v128.const i64x2 0 0

    local.get $first
    i64x2.replace_lane 0

    local.get $last
    i64x2.replace_lane 1

    v128.const i64x2 0 0
    i8x16.shuffle 2 1 0 15  5 4 3 15  8 7 6 15  15 15 15 15

    v128.const i32x4 0x00FFFFFF 0x00FFFFFF 0x00FFFFFF 0x00
    v128.and

    call $ends_with3zkana128v
  )

)
