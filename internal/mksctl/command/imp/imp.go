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

	subArgNum = 2
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
	if _, err := filesFromArgs(args); err != nil {
		return err
	}
	return nil
}

// run returns the command
func run(_ *cobra.Command, args []string) error {
	// Get file from args
	files, err := filesFromArgs(args)
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
	ctx := context.Background()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(f.path))
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, f.file); err != nil {
		return err
	}

	// Add the mainArtifact flag to request.
	if err := writer.WriteField("mainArtifact", strconv.FormatBool(f.mainArtifact)); err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	resp, err := c.Actions().UploadArtifactWithBodyWithResponse(
		ctx,
		nil,
		writer.FormDataContentType(),
		body,
	)
	if err != nil {
		fmt.Printf("Import error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Microcks has discovered \"%s\"\n", resp.Body)

	return nil
}

func filesFromArgs(args []string) ([]file, error) {
	files := []file{}
	for _, arg := range args {
		s := strings.Split(arg, ":")
		path := s[0]
		primary := true
		if len(s) == subArgNum {
			var err error
			primary, err = convertPrimarySubArg(s[1], path)
			if err != nil {
				return files, err
			}
		}

		// Read file
		f, err := os.Open(path)
		if err != nil {
			return files, err
		}
		if os.IsNotExist(err) {
			if f != nil {
				_ = f.Close()
			}
			return files, fmt.Errorf("file \"%s\" does not exists", path)
		}

		files = append(
			files,
			file{path: path, mainArtifact: primary, file: f},
		)
	}
	return files, nil
}

func convertPrimarySubArg(primary, path string) (bool, error) {
	v, err := strconv.ParseBool(primary)
	if err != nil {
		return v, fmt.Errorf("file \"%s\", primary value must be \"true\" or \"false\", actual: %s", path, primary)
	}
	return v, nil
}
