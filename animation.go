// animation.go

package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Interface for animation
type Animation interface {
	Start(wg *sync.WaitGroup, stopChan chan struct{})
}

// Structure and method for "typing" animation
type TypingAnimation struct{}

func (a *TypingAnimation) Start(wg *sync.WaitGroup, stopChan chan struct{}) {
	defer wg.Done()
	steps := []string{"#", "##", "###", "####", "#####"}
	for {
		select {
		case <-stopChan:
			return
		default:
			for _, step := range steps {
				fmt.Printf("\r%s", step)
				time.Sleep(500 * time.Millisecond)
			}
			fmt.Print("\r     \r")
		}
	}
}

// Structure and method for Docker-style animation
type DockerStyleAnimation struct{}

func (a *DockerStyleAnimation) Start(wg *sync.WaitGroup, stopChan chan struct{}) {
	defer wg.Done()
	maxWidth := 50
	for {
		select {
		case <-stopChan:
			return
		default:
			for i := 0; i <= maxWidth; i++ {
				fmt.Printf("\r|%s>%s|", strings.Repeat("=", i), strings.Repeat(" ", maxWidth-i))
				time.Sleep(100 * time.Millisecond)
			}
			fmt.Print("\r      \r")
		}
	}
}

// Structure and method for "speedtest" animation
type SpeedtestAnimation struct{}

func (a *SpeedtestAnimation) Start(wg *sync.WaitGroup, stopChan chan struct{}) {
	defer wg.Done()
	frames := []string{`\`, `|`, `/`, `-`} // Rotating indicator
	width := 0                             // Initial length of the string "="
	for {
		select {
		case <-stopChan:
			return
		default:
			// Iterate through each frame
			for _, frame := range frames {
				fmt.Printf("\r%s%s", strings.Repeat("=", width), frame)
				time.Sleep(100 * time.Millisecond)
			}
			width++
			if width > 50 {
				width = 5 // Reset length after 50
			}
		}
	}
}

// Structure and method for Telegram-style "typing" animation
type TelegramTypingAnimation struct{}

func (a *TelegramTypingAnimation) Start(wg *sync.WaitGroup, stopChan chan struct{}) {
	defer wg.Done()
	frames := []string{
		"●     ", // Dot 1
		"● ●   ", // Dots 1 and 2
		"● ● ●",  // All three dots
		"  ● ●",  // Dots 2 and 3
		"    ●",  // Only dot 3
		"  ● ●",  // Dots 2 and 3
		"● ● ●",  // All three dots
		"● ●   ", // Dots 1 and 2
	}
	// Move cursor down by 3 lines
	fmt.Print("\n\n\n")
	for {
		select {
		case <-stopChan:
			// Clear lines and return cursor to its position
			fmt.Print("\033[3A\r\033[K") // Move up 3 lines and clear the current line
			return
		default:
			// Iterate through each frame
			for _, frame := range frames {
				fmt.Printf("\033[3A\rTyping: %s\n\n\n", frame) // Move cursor up 3 lines for animation
				time.Sleep(250 * time.Millisecond)             // Delay between frames
			}
		}
	}
}

// Factory function to select the type of animation
func GetAnimation(style string) Animation {
	switch style {
	case "docker":
		return &DockerStyleAnimation{}
	case "speedtest":
		return &SpeedtestAnimation{}
	case "telegram":
		return &TelegramTypingAnimation{}
	default:
		return &TypingAnimation{}
	}
}
