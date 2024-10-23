package main

import (
	"image"
	"image/png"
	_ "image/png"
	"log"
	"os"

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
	file, err := os.Open("./img/gh_logo.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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

	// Aquí deberías implementar el código para enviar `command` y `bitmap` a la impresora
	// Esto puede variar dependiendo de cómo estés enviando datos a la impresora.
	// Aquí se muestra un ejemplo de impresión, deberías adaptarlo según tus necesidades:
	log.Println("Enviando comando a la impresora:", command)

	// Ejemplo de envío de datos a la impresora
	// `SendToPrinter` es una función ficticia que deberías implementar
	// err := SendToPrinter(append(command, bitmap...))
	// if err != nil {
	//     log.Println("Error al enviar a la impresora:", err)
	// }

	// Simulación de impresión de los datos
	log.Println("Datos de la imagen enviados a la impresora:", bitmap)
	print(command)
	print(bitmap)
}
