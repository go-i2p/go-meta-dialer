fmt:
	find . -name '*.go' -not -path './vendor/*' -exec gofumpt -s -extra -w {} \;