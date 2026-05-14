# 用我们自己的稀疏 secret（H=192），只跑 ours
go run ./4-29/. -logN=16 -layers=1 -run=ours -paramSet=ours

# 用和 AnyuWang 一样的 dense secret（P=2/3, H≈43691），跑 baseline + ours
go run ./4-29/. -logN=16 -layers=1 -run=all -paramSet=anyu25

# 最公平的三方对比（参数对齐 AnyuWang，同时跑他们的二进制）
go run ./4-29/. -logN=16 -layers=1 -run=all -paramSet=anyu25 -reps=1