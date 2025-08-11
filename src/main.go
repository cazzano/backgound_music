package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		handleRun()
	case "stop":
		handleStop()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleRun() {
	if len(os.Args) < 3 {
		fmt.Println("Error: folder path is required")
		printUsage()
		os.Exit(1)
	}

	folderPath := os.Args[2]
	
	err := RunMusicPlayer(folderPath)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func handleStop() {
	err := StopMusicPlayer()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("🛑 Music player stopped")
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  ./main run <folder_path>")
	fmt.Println("  ./main stop")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ./main run ~/songs")
	fmt.Println("  ./main run /path/to/music")
	fmt.Println("  ./main stop")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  - Plays music in background with shuffle and loop")
	fmt.Println("  - Automatically stops any existing player before starting")
	fmt.Println("  - Use 'stop' command to stop the player")
}
