// thnx to https://github.com/codebench-dev
package main

import (
	"context"
	"os"
	"time"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func copyFile(from string, to string) error {
	data, err := os.ReadFile(from)
	if err != nil {
		log.Error().Msgf("%v", err)
		return err
	}
	err = os.WriteFile(to, data, 0o644)
	if err != nil {
		log.Error().Msgf("%v", err)
		return err
	}
	return nil
}

func startVm() error {
	uuid := uuid.New().String()
	log.Info().Msgf("making a vm with ID %v", uuid)
	// maybe make the below configurable
	const PATH_TO_KERNAL = "./linux/assets/vmlinux"
	socketPath := "./linux/assets/firecracker" + uuid + ".sock"
	pathToRootfs := "./linux/assets/" + uuid + ".ext4"
	err := copyFile("./linux/assets/rootfs.ext4", pathToRootfs)
	if err != nil {
		log.Error().Msgf("failed to copy filesystem: %v", err)
		return err
	}

	stdoutPath := "./linux/assets/" + uuid + "stdout.log"
	stderrPath := "./linux/assets/" + uuid + "stderror.log"
	//--

	cfg := firecracker.Config{
		SocketPath:      socketPath,
		KernelImagePath: PATH_TO_KERNAL,
		Drives:          firecracker.NewDrivesBuilder(pathToRootfs).Build(),
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(4),
			MemSizeMib: firecracker.Int64(256),
		},
	}

	log.Info().Msgf("Finish creating VM config: %v", cfg)

	stdout, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Error().Msgf("failed to create stdout file: %v", err)
		return err
	}

	stderr, err := os.OpenFile(stderrPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Error().Msgf("failed to create stderr file: %v", err)
		return err
	}

	ctx := context.Background()
	// build our custom command that contains our two files to
	// write to during process execution
	cmd := firecracker.VMCommandBuilder{}.
		WithBin("firecracker").
		WithSocketPath(socketPath).
		WithStdout(stdout).
		WithStderr(stderr).
		Build(ctx)

	log.Info().Msgf("Finish creating VM cmd: %v", cmd)
	m, err := firecracker.NewMachine(ctx, cfg, firecracker.WithProcessRunner(cmd))
	if err != nil {
		log.Error().Msgf("failed to create new machine: %v", err)
		return err
	}

	log.Info().Msgf("Finish creating VM : %v", m)

	defer os.Remove(cfg.SocketPath)

	if err := m.Start(ctx); err != nil {
		log.Error().Msgf("failed to initialize machine: %v", err)
		return err
	}

	log.Info().Msgf("Start execute VM: %v", m)
	time.Sleep(2000 * time.Millisecond)

	// wait for VMM to execute
	// if err := m.Wait(ctx); err != nil {
	//	log.Error().Msgf("wait for VMM to execute: %v", err)
	//	return err
	// }

	log.Info().Msgf("Shutting down VM: %v", m)

	m.Shutdown(ctx)

	return nil
}
