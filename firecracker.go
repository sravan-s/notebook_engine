// thnx to https://github.com/codebench-dev
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/rs/zerolog/log"
)

const (
	PATH_TO_NETWORK        = "/etc/cni/conf.d/"
	NETWORK_CONF_EXTENTION = ".conflist"
)

/*
Creates file with name at given PathTo
// should run sudo ip link delete mynet4242 to remove bridges if you cahnge subnet
*/
func placeConfig(name string, pathTo string) error {
	template := `
  {
		"name": "fcnet%s",
		"cniVersion": "0.4.0",
		"plugins": [
			{
				"type": "bridge",
				"ipMasq": true,
				"bridge":"mynet%s",
				"isDefaultGateway": true,
				"ipam": {
					"type": "host-local",
					"subnet": "172.16.0.0/24",
					"resolvConf": "/etc/resolv.conf"
				}
			},
			{
				"type": "firewall"
			},
			{
				"type": "tc-redirect-tap"
			}
		]
	}
  `
	configStr := fmt.Sprintf(template, name, name)
	configBuff := []byte(configStr)
	err := os.WriteFile(pathTo+name+NETWORK_CONF_EXTENTION, configBuff, 0o644)
	if err != nil {
		log.Error().Msgf("create cni config failed %v", err)
		return err
	}
	return nil
}

func startVm(uuid string) (*firecracker.Machine, context.Context, error) {
	log.Info().Msgf("making a vm with ID %v", uuid)
	// maybe make the below configurable
	const PATH_TO_KERNAL = "./linux/assets/vmlinux"
	socketPath := "./linux/assets/firecracker" + uuid + ".sock"
	if err := deleteFileIfExists(socketPath); err != nil {
		log.Error().Msgf("startVm: failed to remove exisiting socket: %v", err)
		// donot need to exit, maybe file didnt exist
	}
	pathToRootfs := "./linux/assets/" + uuid + ".ext4"
	err := copyFile("./linux/assets/rootfs.ext4", pathToRootfs)
	if err != nil {
		log.Error().Msgf("failed to copy filesystem: %v", err)
		return nil, nil, err
	}

	stdoutPath := "./linux/assets/" + uuid + "stdout.log"
	stderrPath := "./linux/assets/" + uuid + "stderror.log"

	networkName := "fcnet" + uuid
	ifName := "veth" + uuid

	err = placeConfig(uuid, PATH_TO_NETWORK)
	if err != nil {
		log.Error().Msgf("failed to create config: %v", err)
		return nil, nil, err
	}

	//--

	cfg := firecracker.Config{
		SocketPath:      socketPath,
		KernelImagePath: PATH_TO_KERNAL,
		Drives:          firecracker.NewDrivesBuilder(pathToRootfs).Build(),
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(4),
			MemSizeMib: firecracker.Int64(512),
		},
		// https://k-jingyang.github.io/firecracker/2024/06/15/firecracker-bridge.html
		NetworkInterfaces: firecracker.NetworkInterfaces{
			firecracker.NetworkInterface{
				AllowMMDS: true,
				CNIConfiguration: &firecracker.CNIConfiguration{
					NetworkName: networkName,
					IfName:      ifName,
				},
			},
		},
	}

	log.Info().Msgf("Finish creating VM config: %v", cfg)

	stdout, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Error().Msgf("failed to create stdout file: %v", err)
		return nil, nil, err
	}

	stderr, err := os.OpenFile(stderrPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Error().Msgf("failed to create stderr file: %v", err)
		return nil, nil, err
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
		return nil, nil, err
	}

	log.Info().Msgf("Finish creating VM : %v", m)

	defer os.Remove(cfg.SocketPath)

	if err := m.Start(ctx); err != nil {
		log.Error().Msgf("failed to initialize machine: %v", err)
		return nil, nil, err
	}

	ip := m.Cfg.NetworkInterfaces[0].StaticConfiguration.IPConfiguration.IPAddr
	log.Info().Msgf("\n\n\nIp--- %v \n", ip)

	return m, ctx, nil
}

func shutdownCleanup(uuid string) error {
	// maybe remove /var/lib/cni on program_shutdown? also - /etc/cni/conf.d/
	// see cleanup-cni.sh
	return os.Remove(PATH_TO_NETWORK + uuid + NETWORK_CONF_EXTENTION)
}
