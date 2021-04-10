module go-swe

go 1.15

require (
	github.com/araddon/dateparse v0.0.0-20210207001429-0eec95c9db7e
	github.com/gin-gonic/gin v1.7.1
	github.com/spf13/cobra v1.1.3
	go-common v0.0.0
	go-common-web v0.0.0
)

replace go-common => ../go-common
replace go-common-web => ../go-common/web
