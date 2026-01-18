package auth

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

var randomSecret string = gotp.RandomSecret(16)

func generateTOTPWithSecret(randomSecret string) {
	uri := gotp.NewDefaultTOTP(randomSecret).ProvisioningUri("user@email.com", "myApp")
	fmt.Println("Secret Key URI:", uri)

	qrcode.WriteFile(uri, qrcode.Medium, 256, "qr.png")

	// Generate and display QR code in the terminal
	qrterminal.GenerateWithConfig(uri, qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    os.Stdout,
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
	})

	fmt.Println("\nScan the QR code with your authenticator app")
}

func verifyOTP(randomSecret string) {
	totp := gotp.NewDefaultTOTP(randomSecret)

	// Wait for user input of the OTP
	fmt.Print("Enter the OTP from your authenticator app: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userInput := scanner.Text()

	// Validate the provided OTP
	if totp.Verify(userInput, time.Now().Unix()) {
		fmt.Println("Authentication successful! Access granted.")
	} else {
		fmt.Println("Authentication failed! Invalid OTP.")
	}
}
