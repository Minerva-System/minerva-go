module github.com/Minerva-System/minerva-go

go 1.20

require google.golang.org/protobuf v1.30.0
require minervarpc v1.0.0
replace minervarpc => ./internal/rpc
