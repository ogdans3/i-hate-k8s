package docker

import (
	"archive/tar"
	"bytes"
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	engine_models "github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/client/engine-interface/engine-models"
)

const (
	nginxConfigDst = "/etc/nginx/"
)

func AddNewNginxConfigurationToContainer(nginxConfig string, container engine_models.Container) error {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	ctx := context.Background()

	if err := copyStringToContainer(ctx, apiClient, container.Id, nginxConfig, nginxConfigDst); err != nil {
		return err
	}

	if err := reloadNginx(ctx, apiClient, container.Id); err != nil {
		return err
	}
	return nil
}

func copyStringToContainer(ctx context.Context, cli *client.Client, containerID, nginxConfig, dstPath string) error {
	// Create a buffer to hold the tar archive
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// Add a file header for nginx.conf in the tar archive
	tarHeader := &tar.Header{
		Name: "nginx.conf",
		Mode: 0600,
		Size: int64(len(nginxConfig)),
	}

	// Write the header and file content to the tar writer
	if err := tw.WriteHeader(tarHeader); err != nil {
		return err
	}
	if _, err := tw.Write([]byte(nginxConfig)); err != nil {
		return err
	}

	// Close the tar writer to finalize the archive
	if err := tw.Close(); err != nil {
		return err
	}

	return cli.CopyToContainer(ctx, containerID, dstPath, buf, container.CopyToContainerOptions{AllowOverwriteDirWithFile: true})
}

// Function to reload Nginx by executing the reload command inside the container
func reloadNginx(ctx context.Context, cli *client.Client, containerID string) error {
	execConfig := container.ExecOptions{
		Cmd:          []string{"nginx", "-s", "reload"},
		AttachStdout: true,
		AttachStderr: true,
	}

	// Create an exec instance for the reload command
	execID, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return err
	}

	// Attach to the exec instance to run the command
	response, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	defer response.Close()

	// Read the response for any output (if necessary)
	//_, err = io.Copy(os.Stdout, response.Reader)
	return err
}
