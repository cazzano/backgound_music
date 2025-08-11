package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type PIDData struct {
	PID    int    `json:"pid"`
	Folder string `json:"folder"`
}

func getPIDFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "mpv_background")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %v", err)
	}

	return filepath.Join(configDir, "pid.json"), nil
}

func savePID(pid int, folder string) error {
	pidFile, err := getPIDFilePath()
	if err != nil {
		return err
	}

	data := PIDData{
		PID:    pid,
		Folder: folder,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal PID data: %v", err)
	}

	return os.WriteFile(pidFile, jsonData, 0644)
}

func loadPID() (*PIDData, error) {
	pidFile, err := getPIDFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return nil, nil // No PID file exists
	}

	data, err := os.ReadFile(pidFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read PID file: %v", err)
	}

	var pidData PIDData
	if err := json.Unmarshal(data, &pidData); err != nil {
		return nil, fmt.Errorf("failed to parse PID file: %v", err)
	}

	return &pidData, nil
}

func deletePIDFile() error {
	pidFile, err := getPIDFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	return os.Remove(pidFile)
}

func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %v", pid, err)
	}

	// Try to kill the process
	if err := process.Signal(syscall.SIGTERM); err != nil {
		// If SIGTERM fails, try SIGKILL
		if err := process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process %d: %v", pid, err)
		}
	}

	return nil
}

func stopExistingPlayer() error {
	pidData, err := loadPID()
	if err != nil {
		return err
	}

	if pidData == nil {
		return nil // No existing player
	}

	fmt.Printf("🛑 Stopping existing player (PID: %d)\n", pidData.PID)
	
	if err := killProcess(pidData.PID); err != nil {
		// Process might already be dead, just clean up the file
		fmt.Printf("⚠️  Could not kill process %d (might already be stopped): %v\n", pidData.PID, err)
	}

	return deletePIDFile()
}

func startMPVBackground(folder string) error {
	// Expand tilde to home directory
	if len(folder) >= 2 && folder[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %v", err)
		}
		folder = filepath.Join(homeDir, folder[2:])
	}

	// Check if folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist: %s", folder)
	}

	// Start mpv in background
	cmd := exec.Command("mpv",
		"--no-video",
		"--loop-playlist",
		"--shuffle",
		"--quiet",
		folder,
	)

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start mpv: %v", err)
	}

	// Save PID
	if err := savePID(cmd.Process.Pid, folder); err != nil {
		// Kill the process if we can't save PID
		cmd.Process.Kill()
		return fmt.Errorf("failed to save PID: %v", err)
	}

	fmt.Printf("🎵 Started background music player\n")
	fmt.Printf("📁 Folder: %s\n", folder)
	fmt.Printf("🆔 PID: %d\n", cmd.Process.Pid)
	fmt.Printf("🔀 Shuffle and loop enabled\n")

	return nil
}

// Public functions for main.go

func RunMusicPlayer(folder string) error {
	// First, stop any existing player
	if err := stopExistingPlayer(); err != nil {
		fmt.Printf("⚠️  Warning: Could not stop existing player: %v\n", err)
	}

	// Start new player
	return startMPVBackground(folder)
}

func StopMusicPlayer() error {
	pidData, err := loadPID()
	if err != nil {
		return fmt.Errorf("failed to load PID: %v", err)
	}

	if pidData == nil {
		return fmt.Errorf("no music player is running")
	}

	if err := killProcess(pidData.PID); err != nil {
		return fmt.Errorf("failed to stop player: %v", err)
	}

	if err := deletePIDFile(); err != nil {
		return fmt.Errorf("failed to clean up PID file: %v", err)
	}

	return nil
}
