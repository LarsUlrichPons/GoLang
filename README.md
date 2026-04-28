# Go Virtual Piano Synthesizer 🎹

A realistic, interactive virtual piano and audio synthesizer built entirely in Go. This project uses pure mathematical sine waves to generate sound in real-time, completely eliminating the need for external audio files. 

It features a full chromatic octave (C4 to B4), an interactive graphical interface, and a dynamic LCD screen that tracks frequencies.

## ✨ Features
* **Real-Time Audio Synthesis:** Generates 44.1kHz sine waves on the fly using Go's `math` package.
* **Interactive UI:** Visually responds to key presses with realistic key depressions and highlighting.
* **Live LCD Display:** Shows the currently played musical note and its exact mathematical frequency in Hertz (Hz).
* **Polyphony Support:** Play multiple notes or chords simultaneously.
* **Zero External Assets:** All visuals (shapes, colors, text) and audio are rendered natively by the engine.

## 🚀 Getting Started

### Prerequisites
You will need [Go](https://go.dev/doc/install) installed on your machine. 

### Installation & Running
1. Clone the repository to your local machine:
   ```bash
   git clone [https://github.com/YOUR-USERNAME/YOUR-REPO-NAME.git](https://github.com/YOUR-USERNAME/YOUR-REPO-NAME.git)
