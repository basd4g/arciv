package commands

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Run:   initCommand,
		Short: "Initialize a repository",
		Long: `Initialize a repository.
The repository's root directory specifies the current directory by generating '.arciv' directory on the current directory.`,
		Args: cobra.NoArgs,
	}
)

func initCommand(cmd *cobra.Command, args []string) {
	if err := initAction(); err != nil {
		Exit(err, 1)
	}
}

func initAction() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	return Repository{Name: "self", Path: currentDir, PathType: PATH_FILE}.Init()
}

func (repository Repository) Init() error {
	if repository.PathType == PATH_FILE {
		// Create dirs if there does not exist
		err := fileOp.mkdirAll(repository.Path + "/.arciv/list")
		if err != nil {
			return err
		}
		err = fileOp.mkdirAll(repository.Path + "/.arciv/blob")
		if err != nil {
			return err
		}

		// Create files if there does not exist
		paths, err := fileOp.findFilePaths(repository.Path + "/.arciv")
		if err != nil {
			return err
		}
		if !isIncluded(paths, "repositories") {
			fileOp.writeLines(repository.Path+"/.arciv/repositories", []string{})
		}
		if !isIncluded(paths, "timeline") {
			fileOp.writeLines(repository.Path+"/.arciv/timeline", []string{})
		}
		return nil
	}
	if repository.PathType == PATH_S3 {
		bucketName = repository.Path
		prepareS3BucketClient()
		// FIXME: add a func, s3Op.Exist("  key string  ")
		_, err := s3Op.loadLines(".arciv/repositories")
		if err != nil {
			err = s3Op.writeLines(".arciv/repositories", []string{})
			if err != nil {
				return err
			}
		}
		_, err = s3Op.loadLines(".arciv/timeline")
		if err != nil {
			err = s3Op.writeLines(".arciv/timeline", []string{})
			if err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("Repository's PathType must be PATH_FILE or PATH_S3")
}

func init() {
	RootCmd.AddCommand(initCmd)
}
