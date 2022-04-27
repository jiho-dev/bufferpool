all:
	#go test -v -run TestPool -gcflags="-m"
	go test  -run TestUsage #-gcflags="-m"

bench:
	#go test -v -run TestUsage -gcflags="-m"
	#go test -bench -run TestUsage -gcflags -m=2 -memprofile p.out
	#go test -bench=. -gcflags="-m=2" -count=1 -benchmem
	#go test -cover -count 3  -benchmem  -bench=.
	#go test -bench=. -count=2 -benchmem
	#go test -bench=^BenchmarkPoolMyString . -run=^Nothing . -v -benchtime 5s -benchmem -count=2
	go test -count=3 -benchmem -bench=^BenchmarkPoolMyString . -run=^Nothing .


