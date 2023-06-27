package version

import (
	"io"
	"time"

	"github.com/drewstinnett/gout/v2"
)

var (
	// VersionCLI is the current version of this CLI
	VersionCLI = "dev"
	// VersionServer is the version of the related server
	VersionServer = "n/a"
	// Commit is the commit hash this build was created from
	Commit = "n/a"
	// RawDate is the time when this build was created in raw string
	RawDate = "n/a"
)

// version represents a version
type version struct {
	VersionServer string `yaml:"versionServer" json:"version_server"`
	VersionCLI    string `yaml:"versionCLI" json:"version_cli"`
	Commit        string `yaml:"commit" json:"commit"`
	Date          string `yaml:"date" json:"date"`
}

// Date returns the version's date
func Date() (time.Time, error) {
	t, err := time.Parse(time.RFC3339, RawDate)
	if err != nil {
		return t, &DateParseError{Date: RawDate, Err: err}
	}

	return t, nil
}

// Print prints the version
func Print(w io.Writer, clientOnly bool, outputFormat string) {
	// TODO: Start: Move formatter business logic
	g, _ := gout.New()
	g.SetWriter(w)
	formatter, ok := gout.BuiltInFormatters[outputFormat]
	if ok {
		g.SetFormatter(formatter)
	}
	// TODO: End: Move formatter business logic

	var v = version{
		VersionCLI:    VersionCLI,
		VersionServer: VersionServer,
		Commit:        Commit,
		Date:          RawDate,
	}
	if !clientOnly {
		// TODO: Fetch server version
		v.VersionServer = VersionServer
	}
	g.MustPrint(v)
}
