package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/moby/moby/pkg/stdcopy"
)

func RunPythonContainer(input string) ([]byte, error) {
	ctx := context.Background()

	// Создаём Docker-клиент
	cli, err := client.NewClientWithOpts(client.WithVersion("1.47"))
	if err != nil {
		log.Println("Failed to create Docker client:", err)
		return nil, err
	}

	// Создаём контейнер
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        "python-worker",
			OpenStdin:    true,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          []string{"python", "/worker.py", input},
		},
		nil,
		nil,
		nil,
		"",
	)
	if err != nil {
		log.Println("Failed to create container:", err)
		return nil, err
	}

	// Запускаем контейнер
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Println("Failed to start container:", err)
		return nil, err
	}

	// Получаем логи контейнера
	out, err := cli.ContainerLogs(
		ctx,
		resp.ID,
		container.LogsOptions{
			ShowStdout: true,
			ShowStderr: false,
			Follow:     true,
			Timestamps: false,
		},
	)
	if err != nil {
		log.Println("Failed to get container logs:", err)
		return nil, err
	}
	defer out.Close()

	// Разделяем потоки stdout и stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, out)
	if err != nil {
		log.Println("Failed to copy logs:", err)
		return nil, err
	}

	// Получаем вывод из stdout
	encodedString := stdoutBuf.String()

	// Декодируем base64-строку в бинарные данные
	decoded, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		log.Printf("Failed to decode base64: %v, string length: %d\n", err, len(encodedString))
		return nil, err
	}

	// Ждём завершения контейнера
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Println("Container wait error:", err)
			return nil, err
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			log.Printf("Container exited with status %d\n", status.StatusCode)
			return nil, fmt.Errorf("container exited with status %d", status.StatusCode)
		}
	}

	// Удаляем контейнер
	if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
		log.Println("Failed to remove container:", err)
	}

	return decoded, nil
}