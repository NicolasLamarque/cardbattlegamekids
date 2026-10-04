package lanserver

import (
	"encoding/base64"

	qrcode "github.com/skip2/go-qrcode"
)

// QRCodeDataURI encode une URL en QR code PNG, prêt à mettre dans un <img src="...">.
func QRCodeDataURI(content string) (string, error) {
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), nil
}
