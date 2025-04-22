# Structured log syslog

[![Go Reference](https://pkg.go.dev/badge/github.com/mocheryl/slog-syslog.svg)](https://pkg.go.dev/github.com/mocheryl/slog-syslog)

Go structured log syslog handler.

## Requirements

- Go >= 1.21

## Usage

``` go
package main

import (
	"log"
	"log/slog"

	slogsyslog "github.com/mocheryl/slog-syslog"
)

func main() {
	h, err := slogsyslog.New(nil)
	if err != nil {
		log.Fatalf("Could not initialize syslog slog handler: %s\n", err)
		return
	}

	l := slog.New(h)
	slog.SetDefault(l)

	slog.Info("Hello, World!")
	// Output: Jun 29 08:44:16 localhost helloworld[1234]: Hello, World!

	h.Close()
}
```

## Credits

Most of the code and ideas taken from the following projects:
- Structured log text handler from the standard library.
- Syslog client from the standard library.
- [RackSpace Managed Security Development srslog](https://github.com/RackSec/srslog).
- [HashiCorp go-syslog](https://github.com/hashicorp/go-syslog).

## Contributing

Any contributions are welcome - see the
[CONTRIBUTING.md](.github/CONTRIBUTING.md) file for details.

## License

This project is licensed under the 3-clause BSD license - see the
[LICENSE](LICENSE) file for details.
