/*
 * Copyright 2017-2022 Provide Technologies Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func ListContainers(docker *client.Client, stack string) []types.Container {
	containers, err := docker.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs([]filters.KeyValuePair{
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-api", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-consumer", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-elasticsearch", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-nats", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-postgres", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-redis", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-ident", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-ident-consumer", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-nchain-api", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-nchain-consumer", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-privacy-api", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-privacy-consumer", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-reachabilitydaemon", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-statsdaemon", strings.ReplaceAll(stack, " ", "")),
			},
			{
				Key:   "name",
				Value: fmt.Sprintf("%s-vault", strings.ReplaceAll(stack, " ", "")),
			},
		}...),
	})
	if err != nil {
		log.Printf("failed to list containers; %s", err.Error())
		os.Exit(1)
	}

	return containers
}

func LogContainers(docker *client.Client, wg *sync.WaitGroup, stack string) error {
	for _, container := range ListContainers(docker, stack) {
		if wg != nil {
			wg.Add(1)
		}

		containerID := make([]byte, len(container.ID))
		copy(containerID, container.ID)

		go func() {
			LogContainer(docker, string(containerID))
			if wg != nil {
				wg.Done()
			}
		}()
	}

	return nil
}

func LogContainer(docker *client.Client, containerID string) error {
	out, err := docker.ContainerLogs(context.Background(), containerID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
	})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return nil
}

func PurgeContainers(docker *client.Client, stack string, purgeVolumes bool) {
	for _, container := range ListContainers(docker, stack) {
		if !purgeVolumes {
			timeout := time.Millisecond * 5000
			err := docker.ContainerStop(context.Background(), container.ID, &timeout)

			if err != nil {
				log.Printf("WARNING: failed to stop container: %s; %s", container.Names[0], err.Error())
			}
		} else {
			err := docker.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			})

			if err != nil {
				log.Printf("WARNING: failed to remove container: %s; %s", container.Names[0], err.Error())
			}
		}
	}
}

func PurgeNetwork(docker *client.Client, stack string) {
	networks, _ := docker.NetworkList(context.Background(), types.NetworkListOptions{})
	for _, ntwrk := range networks {
		if ntwrk.Name == stack {
			docker.NetworkRemove(context.Background(), ntwrk.ID)
		}
	}
}

func StopContainers(docker *client.Client, stack string) {
	for _, container := range ListContainers(docker, stack) {
		timeout := time.Millisecond * 5000
		err := docker.ContainerStop(context.Background(), container.ID, &timeout)

		if err != nil {
			log.Printf("WARNING: failed to stop container: %s; %s", container.Names[0], err.Error())
		}
	}
}
