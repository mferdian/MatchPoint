package utils

import (
	"bytes"
	"fieldreserve/dto"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

func GenerateInvoicePDF(booking dto.BookingFullResponse) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set margins
	pdf.SetMargins(20, 20, 20)

	// Generate QR code for booking ID
	qrCode, err := qrcode.Encode(booking.BookingID.String(), qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Register QR code image
	qrReader := bytes.NewReader(qrCode)
	pdf.RegisterImageOptionsReader("qr", gofpdf.ImageOptions{ImageType: "PNG"}, qrReader)

	// Header Section
	drawHeader(pdf)

	// Invoice Title
	pdf.SetY(35)
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(41, 128, 185)
	pdf.CellFormat(0, 15, "INVOICE BOOKING LAPANGAN", "", 1, "C", false, 0, "")

	// Separator line
	pdf.SetDrawColor(41, 128, 185)
	pdf.SetLineWidth(0.5)
	pdf.Line(20, 55, 190, 55)

	pdf.SetY(65)

	// Customer Information Section
	drawSectionHeader(pdf, "INFORMASI BOOKING")

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 11)

	if booking.User.Name != "" {
		drawInfoRow(pdf, "Nama Pemesan:", booking.User.Name)
	}

	drawInfoRow(pdf, "Nama Lapangan:", booking.Field.FieldName)
	drawInfoRow(pdf, "Tanggal Booking:", booking.BookingDate.Format("Monday, 02 January 2006"))

	timeStr := fmt.Sprintf("%s - %s WIB",
		booking.StartTime.Format("15:04"),
		booking.EndTime.Format("15:04"))
	drawInfoRow(pdf, "Waktu:", timeStr)

	duration := booking.EndTime.Sub(booking.StartTime)
	durationStr := fmt.Sprintf("%.0f jam", duration.Hours())
	drawInfoRow(pdf, "Durasi:", durationStr)

	pdf.Ln(10)

	// Payment Information Section
	drawSectionHeader(pdf, "INFORMASI PEMBAYARAN")

	pdf.SetFont("Arial", "", 11)

	drawInfoRow(pdf, "Metode Pembayaran:", booking.PaymentMethod)

	statusColor := getStatusColor(booking.Status)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(50, 8, "Status:", "", 0, "L", false, 0, "")
	pdf.SetTextColor(statusColor[0], statusColor[1], statusColor[2])
	pdf.CellFormat(0, 8, booking.Status, "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	if booking.PaymentVerifiedAt != nil {
		pdf.SetFont("Arial", "", 11)
		drawInfoRow(pdf, "Diverifikasi pada:", booking.PaymentVerifiedAt.Format("02 Jan 2006 15:04 WIB"))
	}

	pdf.Ln(10)

	// Total Payment Section
	drawTotalPaymentSection(pdf, booking.TotalPayment)

	pdf.Ln(10)

	// Draw QR code (big, centered)
	x := (210 - 60) / 2 // 210mm width, center QR code (60mm)
	pdf.ImageOptions("qr", float64(x), pdf.GetY(), 60, 60, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	pdf.Ln(65)

	// Footer
	drawFooter(pdf)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func drawHeader(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(52, 73, 94)
	pdf.CellFormat(0, 8, "FIELD RESERVE", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 6, "Sistem Reservasi Lapangan Olahraga", "", 1, "L", false, 0, "")
}

func drawSectionHeader(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(52, 73, 94)
	pdf.SetFillColor(236, 240, 241)
	pdf.CellFormat(0, 10, title, "", 1, "L", true, 0, "")
	pdf.Ln(3)
}

func drawInfoRow(pdf *gofpdf.Fpdf, label, value string) {
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(50, 8, label, "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 8, value, "", 1, "L", false, 0, "")
}

func drawTotalPaymentSection(pdf *gofpdf.Fpdf, totalPayment float64) {
	pdf.SetFillColor(46, 204, 113)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 14)

	pdf.Rect(20, pdf.GetY(), 170, 15, "F")
	pdf.CellFormat(0, 15, fmt.Sprintf("TOTAL PEMBAYARAN: Rp %s", formatCurrency(totalPayment)), "", 1, "C", false, 0, "")
}


func drawFooter(pdf *gofpdf.Fpdf) {
	pdf.SetY(-30)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(127, 140, 141)
	pdf.CellFormat(0, 6, fmt.Sprintf("Invoice dicetak pada: %s", time.Now().Format("Monday, 02 January 2006 15:04 WIB")), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, "Terima kasih telah menggunakan layanan Field Reserve", "", 1, "C", false, 0, "")
}

func getStatusColor(status string) [3]int {
	switch status {
	case "CONFIRMED", "PAID", "COMPLETED":
		return [3]int{46, 204, 113}
	case "PENDING":
		return [3]int{241, 196, 15}
	case "CANCELLED":
		return [3]int{231, 76, 60}
	default:
		return [3]int{52, 73, 94}
	}
}

func formatCurrency(amount float64) string {
	return fmt.Sprintf("%.0f", amount)
}
