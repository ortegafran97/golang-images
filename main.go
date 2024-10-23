package main

import (
	"image"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"regexp"

	"github.com/godoes/printers"
	"github.com/nfnt/resize"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
		return
	}
}

func print(data []byte) {
	printerName, err := printers.GetDefault()
	check(err)
	printer, err := printers.Open(printerName)
	check(err)
	defer printer.Close()

	log.Printf("Default printer -> %s", printerName)

	printer.StartDocument("test", "RAW")
	printer.Write(data)
	printer.EndDocument()

}

const (
	// Modo de impresión para la impresora térmica
	ModeHighDensity = 33
)

// Función principal
func main() {
	// Abrir imagen PNG
	args := os.Args

	var file *os.File
	var err error

	if len(args) == 2 {
		file, err = checkImage(args[1])
	} else {
		file, err = os.Open("./img/gh_logo.png")
		log.Println("Default file")
	}

	log.Printf("File found -> %s", file.Name())
	defer file.Close()

	check(err)

	// Decodificar la imagen
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// Escalar la imagen a 200x200 usando la función resize.Resize
	newImage := resize.Resize(200, 200, img, resize.Lanczos3)

	// Convertir la imagen a bitmap
	bitmap, _, _ := imageToBitmapBytes(newImage)

	// Enviar bitmap a la impresora
	sendBitmapToPrinter(bitmap, newImage.Bounds().Dx(), newImage.Bounds().Dy(), ModeHighDensity)

}

// Convierte una imagen a bitmap en blanco y negro y luego a []byte
func imageToBitmapBytes(img image.Image) ([]byte, int, int) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Asegúrate de que el ancho sea múltiplo de 8 para el formato de la impresora
	paddedWidth := (width + 7) / 8 * 8
	bitmap := make([]byte, (paddedWidth/8)*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			byteIndex := (y * paddedWidth / 8) + (x / 8)
			bitIndex := uint(7 - (x % 8))

			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8((r + g + b) / 3 >> 8) // Promedio de los colores RGB

			// Umbral para determinar si el pixel es negro o blanco
			if gray < 128 {
				bitmap[byteIndex] |= (1 << bitIndex) // Establece el bit en 1 para negro
			}
		}
	}

	return bitmap, paddedWidth, height
}

// Envía un bitmap a la impresora térmica
func sendBitmapToPrinter(bitmap []byte, width int, height int, mode int) {
	// Comando de inicio para imprimir imagen
	command := []byte{0x1B, 0x28, 0x2E, 0x01, byte(mode), 0x00}

	printable := append(command, bitmap...)

	print(printable)
}

func checkImage(imgPath string) (*os.File, error) {
	r, _ := regexp.Compile("\\W*.png")

	if !r.MatchString(imgPath) {
		log.Fatal("La imagen no es PNG")
	}

	return os.Open(imgPath)
}
