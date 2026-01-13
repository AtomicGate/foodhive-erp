package correspondence

import (
	"context"
	"errors"
	"time"
)

type Document struct {
	ID          int       `json:"id"`
	SenderID    int       `json:"sender_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FilePath    string    `json:"file_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DocumentRecipient struct {
	ID             int       `json:"id"`
	DocumentID     int       `json:"document_id"`
	RecipientEmpID *int      `json:"recipient_emp_id,omitempty"`
	RecipientDepID *int      `json:"recipient_dep_id,omitempty"`
	Status         string    `json:"status"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DocumentRead struct {
	ID         int       `json:"id"`
	DocumentID int       `json:"document_id"`
	ReaderID   int       `json:"reader_id"`
	ReadAt     time.Time `json:"read_at"`
}

// Repository interfaces defined in core but implemented outside
type Repository interface {
	// Document methods
	CreateDocument(ctx context.Context, doc *Document) (int, error)
	GetDocumentByID(ctx context.Context, id int) (*Document, error)
	// Recipient methods
	CreateRecipient(ctx context.Context, recipient *DocumentRecipient) error
	UpdateRecipientStatus(ctx context.Context, recipientID int, status string) error
	GetRecipientsByDocumentID(ctx context.Context, documentID int) ([]*DocumentRecipient, error)
	// Read methods
	CreateDocumentRead(ctx context.Context, read *DocumentRead) error
	GetReadsByDocumentID(ctx context.Context, documentID int) ([]*DocumentRead, error)
}

type StorageService interface {
	UploadFile(bucketName, fileName string, data []byte) (string, error)
	DownloadFile(bucketName, filePath string) ([]byte, error)
}

type Service interface {
	SendDocument(ctx context.Context, doc *Document, recipientEmpIDs []int, fileContent []byte, fileName string) (int, error)
	ReadDocument(ctx context.Context, documentID, readerID int) ([]byte, error)
	GetDocumentRecipients(ctx context.Context, documentID int) ([]*DocumentRecipient, error)
	GetDocumentReads(ctx context.Context, documentID int) ([]*DocumentRead, error)
}

type service struct {
	repo           Repository
	storageService StorageService
}

func NewService(repo Repository, storageService StorageService) Service {
	return &service{
		repo:           repo,
		storageService: storageService,
	}
}

func (svc *service) SendDocument(ctx context.Context, doc *Document, recipientEmpIDs []int, fileContent []byte, fileName string) (int, error) {
	// Upload file
	filePath, err := svc.storageService.UploadFile("correspondence", fileName, fileContent)
	if err != nil {
		return 0, err
	}
	doc.FilePath = filePath

	// Create document in DB
	docID, err := svc.repo.CreateDocument(ctx, doc)
	if err != nil {
		return 0, err
	}
	doc.ID = docID

	// Add recipients
	for _, recipientID := range recipientEmpIDs {
		recipient := &DocumentRecipient{
			DocumentID:     docID,
			RecipientEmpID: &recipientID,
			Status:         "SENT",
		}
		err := svc.repo.CreateRecipient(ctx, recipient)
		if err != nil {
			return 0, err
		}
	}

	return docID, nil
}

func (svc *service) ReadDocument(ctx context.Context, documentID, readerID int) ([]byte, error) {
	// Check if the reader is a recipient
	recipients, err := svc.repo.GetRecipientsByDocumentID(ctx, documentID)
	if err != nil {
		return nil, err
	}

	var isRecipient bool
	var recipientID int
	for _, recipient := range recipients {
		if recipient.RecipientEmpID != nil && *recipient.RecipientEmpID == readerID {
			isRecipient = true
			recipientID = recipient.ID
			break
		}
	}

	if !isRecipient {
		return nil, errors.New("user is not a recipient of this document")
	}

	// Update recipient status to 'READ'
	err = svc.repo.UpdateRecipientStatus(ctx, recipientID, "READ")
	if err != nil {
		return nil, err
	}

	// Log the read action
	read := &DocumentRead{
		DocumentID: documentID,
		ReaderID:   readerID,
	}
	err = svc.repo.CreateDocumentRead(ctx, read)
	if err != nil {
		return nil, err
	}

	// Retrieve the document
	doc, err := svc.repo.GetDocumentByID(ctx, documentID)
	if err != nil {
		return nil, err
	}

	// Download the document file
	fileContent, err := svc.storageService.DownloadFile("correspondence", doc.FilePath)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

func (svc *service) GetDocumentRecipients(ctx context.Context, documentID int) ([]*DocumentRecipient, error) {
	return svc.repo.GetRecipientsByDocumentID(ctx, documentID)
}

func (svc *service) GetDocumentReads(ctx context.Context, documentID int) ([]*DocumentRead, error) {
	return svc.repo.GetReadsByDocumentID(ctx, documentID)
}
