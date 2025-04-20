package repositories

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	pdf "github.com/adrg/go-wkhtmltopdf"
	"github.com/joaquinamado/gobank/internal/app/services"
	"github.com/joaquinamado/gobank/internal/app/storage"
	"github.com/joaquinamado/gobank/internal/app/types"
)

type transfer interface {
	CreateTransfer(*types.TransferRequest, int) (*types.Transfer, error)
	GetTransfers(*types.PaginationQuery) ([]*types.Transfer, error)
	CreateInvoice(int) ([]byte, error)
}

type transferRepo struct {
	store storage.Storage
	acc   account
}

func (t *transferRepo) CreateTransfer(transferReq *types.TransferRequest, senderNumber int) (*types.Transfer, error) {

	if senderNumber == transferReq.ToAccount {
		return nil, fmt.Errorf("Cannot send transers to same account number")
	}

	sender, err := t.acc.GetAccountByNumber(senderNumber)
	if err != nil {
		return nil, err
	}

	if sender.Balance-int64(transferReq.Amount) < 0 {
		return nil, fmt.Errorf("Invalid operation check balance")
	}

	receiver, err := t.acc.GetAccountByNumber(transferReq.ToAccount)
	if err != nil {
		return nil, err
	}

	transfer := &types.Transfer{
		SenderId:   sender.ID,
		ReceiverId: receiver.ID,
		Amount:     int64(transferReq.Amount),
		CreatedAt:  time.Now().UTC(),
	}

	err = t.store.Transfer.CreateTransfer(transfer)

	if err != nil {
		return nil, err
	}

	return transfer, nil
}

func (t *transferRepo) GetTransfers(query *types.PaginationQuery) ([]*types.Transfer, error) {
	// Default vaule for ints
	if query.Size == 0 {
		query.Size = 20
	}

	return t.store.Transfer.GetTransfers(query)
}

func (t *transferRepo) CreateInvoice(id int) ([]byte, error) {
	invoice, err := t.store.Transfer.GetInvoiceInformation(id)
	if err != nil {
		return nil, err
	}
	return createPdf(invoice)
}

func createPdf(transfer *types.TransferInvoice) ([]byte, error) {

	file, err := os.ReadFile("./templates/invoice.html")

	if err != nil {
		return nil, err
	}

	imageBytes, err := os.ReadFile("./templates/images/icon.png")

	if err != nil {
		return nil, err
	}

	base64Image := base64.StdEncoding.EncodeToString(imageBytes)
	imgTag := fmt.Sprintf("data:image/png;base64,%s", base64Image)

	content := string(file)

	content = strings.ReplaceAll(content, "./images/icon.png", imgTag)
	content = strings.
		ReplaceAll(content, "{{date}}", transfer.CreatedAt.Format("2006-Jan-02"))
	content = strings.
		ReplaceAll(content, "{{full_name_sender}}", transfer.SenderFullName)
	content = strings.
		ReplaceAll(content, "{{account_sender}}", fmt.Sprintf("%d", transfer.SenderAccount))
	content = strings.
		ReplaceAll(content, "{{full_name_receiver}}", transfer.ReceiverFullName)
	content = strings.
		ReplaceAll(content, "{{account_receiver}}", fmt.Sprintf("%d", transfer.ReceiverAccount))
	content = strings.
		ReplaceAll(content, "{{amount}}", fmt.Sprintf("%d", transfer.Amount))

	tmpFile, err := os.CreateTemp("", "invoice-*.html")

	if err != nil {
		return nil, err
	}

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	if err != nil {
		return nil, err
	}

	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	if err := pdf.Init(); err != nil {
		return nil, err
	}
	defer pdf.Destroy()

	page, err := pdf.NewObjectFromReader(tmpFile)

	page.LoadImages = true
	page.BlockLocalFileAccess = false

	tmpFile.Close()

	if err != nil {
		return nil, err
	}

	converter, err := pdf.NewConverter()

	if err != nil {
		return nil, err
	}

	defer converter.Destroy()

	converter.Add(page)

	title := fmt.Sprintf("invoice_%d_%d_%d", transfer.CreatedAt.UnixMilli(), transfer.SenderAccount, transfer.ReceiverAccount)

	converter.Title = title

	var buff bytes.Buffer

	if err = services.CallFunc(func() error {
		// Run converter. Due to a limitation of the `wkhtmltox` library, the
		// conversion must be performed on the main thread.
		return converter.Run(&buff)
	}); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil

}
