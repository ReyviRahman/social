package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/ReyviRahman/social/internal/store"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=255"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	go func() {
		// Tambahkan recovery supaya kalau ada error fatal, server tidak mati (Exit Code 2)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC di email goroutine: %v", r)
			}
		}()

		// Panggil fungsi kirim email
		err := app.sendHostingerEmail(user.Email, plainToken)

		// JANGAN panggil helper response (seperti badRequestResponse) di sini
		// Cukup cetak ke terminal pakai log
		if err != nil {
			log.Printf("Gagal kirim email aktivasi ke %s: %v", user.Email, err)
			return
		}

		log.Printf("Email aktivasi berhasil dikirim ke %s", user.Email)
	}()

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) sendHostingerEmail(userEmail string, token string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// 1. Data untuk Template (Hanya Link)
	data := struct {
		ActivationLink string
	}{
		ActivationLink: "http://localhost:8080/users/activate/" + token,
	}

	// 2. Parse Template
	tmpl, err := template.ParseFiles(filepath.Join("templates", "user_invitation.html"))
	if err != nil {
		return fmt.Errorf("parse template error: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("execute template error: %w", err)
	}

	// 3. Susun Email Header & Content
	// Menambahkan header To: agar email lebih valid di mata provider email
	subject := "Subject: Konfirmasi Aktivasi Akun UNeeNDA\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	fromHeader := fmt.Sprintf("From: UNeeNDA <%s>\n", from)
	toHeader := fmt.Sprintf("To: %s\n", userEmail)

	msg := append([]byte(fromHeader+toHeader+subject+mime+"\n"), body.Bytes()...)

	// 4. Pengiriman via TLS (Port 465)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Karena sebelumnya berhasil dengan true
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}
	if err = client.Rcpt(userEmail); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write(msg); err != nil {
		return err
	}

	return w.Close()
}
