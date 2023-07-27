package imp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/JulienBreux/mksctl/internal/mksctl/api/client"
	"github.com/JulienBreux/mksctl/internal/mksctl/config"
	"github.com/spf13/cobra"
)

const (
	cmdName      = "import FILE:PRIMARY ..."
	cmdShortDesc = "Import API artifacts into Mickrocks."

	exactNSubArgs = 2
)

var (
	cmd = &cobra.Command{
		Use:     cmdName,
		Short:   cmdShortDesc,
		Aliases: []string{"i"},
		Args: cobra.MatchAll(
			cobra.MinimumNArgs(1),
		),
		PreRunE: preRun,
		RunE:    run,
	}
)

type file struct {
	mainArtifact bool

	path string
	file *os.File
}

// New returns a command to import
func New() *cobra.Command {
	return cmd
}

func preRun(_ *cobra.Command, args []string) error {
	if _, err := argsToFiles(args); err != nil {
		return err
	}
	return nil
}

// run returns the command
func run(_ *cobra.Command, args []string) error {
	// Get file from args
	files, err := argsToFiles(args)
	if err != nil {
		return err
	}

	// Connect to client
	cli, err := client.New(config.Config.APIURL)
	if err != nil {
		return err
	}

	// Upload files
	for _, file := range files {
		if err := upload(cli, file); err != nil {
			return err
		}
	}

	return nil
}

////////// TODO: MOVE!

func upload(c client.Client, f file) error {
	contentType, body, err := prepareBody(f)
	if err != nil {
		return err
	}

	resp, err := c.Actions().UploadArtifactWithBodyWithResponse(
		context.Background(),
		nil,
		contentType,
		body,
	)

	if err != nil {
		fmt.Printf("Import error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Microcks has discovered \"%s\"\n", resp.Body)

	return nil
}

func prepareBody(f file) (string, *bytes.Buffer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	contentType := writer.FormDataContentType()

	part, err := writer.CreateFormFile("file", filepath.Base(f.path))
	if err != nil {
		return contentType, nil, err
	}

	if _, err = io.Copy(part, f.file); err != nil {
		return contentType, nil, err
	}

	if err := writer.WriteField("mainArtifact", strconv.FormatBool(f.mainArtifact)); err != nil {
		return contentType, nil, err
	}

	if err = writer.Close(); err != nil {
		return contentType, nil, err
	}

	return contentType, body, nil
}

func argsToFiles(args []string) ([]file, error) {
	files := []file{}

	for _, arg := range args {
		file, err := argToFile(arg)
		if err != nil {
			return files, err
		}
		file, err = readFileContent(*file)
		if err != nil {
			return files, err
		}
		files = append(files, *file)
	}

	return files, nil
}

func argToFile(arg string) (*file, error) {
	a := strings.Split(arg, ":")
	// Check sub arguments length
	if len(a) != exactNSubArgs {
		return nil, fmt.Errorf("unable to decode argument")
	}
	// Check primary argument value
	mainArtifact, err := strconv.ParseBool(a[1])
	if err != nil {
		return nil, fmt.Errorf("file \"%s\", primary value must be \"true\" or \"false\", actual: %s", a[0], a[1])
	}
	return &file{path: a[0], mainArtifact: mainArtifact}, nil
}

func readFileContent(f file) (*file, error) {
	of, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	if os.IsNotExist(err) {
		if of != nil {
			_ = of.Close()
		}
		return nil, fmt.Errorf("file \"%s\" does not exists", f.path)
	}
	f.file = of
	return &f, nil
}
