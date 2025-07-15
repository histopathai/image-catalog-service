package adapter

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/histopathai/image-catalog-service/internal/models"
)

type FirestoreImageRepository struct {
	client     *firestore.Client
	collection *firestore.CollectionRef
}

func NewFirestoreCollection(client *firestore.Client, collectionName string) (*FirestoreImageRepository, error) {
	return &FirestoreImageRepository{
		client:     client,
		collection: client.Collection(collectionName),
	}, nil
}

func (r *FirestoreImageRepository) Read(ctx context.Context, imageID string) (*models.Image, error) {
	doc, err := r.collection.Doc(imageID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}
	var image models.Image
	if err := doc.DataTo(&image); err != nil {
		return nil, fmt.Errorf("failed to convert document to image: %w", err)
	}
	return &image, nil
}

func (r *FirestoreImageRepository) Update(ctx context.Context, image *models.Image) error {
	updates := make([]firestore.Update, 0)

	if image.DiseaseType != nil {
		updates = append(updates, firestore.Update{
			Path:  "disease_type",
			Value: image.DiseaseType,
		})
	}

	if image.Classification != nil {
		updates = append(updates, firestore.Update{
			Path:  "classification",
			Value: image.Classification,
		})
	}

	if image.SubType != nil {
		updates = append(updates, firestore.Update{
			Path:  "subtype",
			Value: image.SubType,
		})
	}

	if image.Grade != nil {
		updates = append(updates, firestore.Update{
			Path:  "grade",
			Value: image.Grade,
		})
	}

	if len(updates) == 0 {
		return nil // No updates to apply
	}
	_, err := r.collection.Doc(image.ID).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("failed to update image: %w", err)
	}
	return nil
}

func (r *FirestoreImageRepository) Delete(ctx context.Context, imageID string) error {
	_, err := r.collection.Doc(imageID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	return nil
}

func (r *FirestoreImageRepository) Filter(ctx context.Context, filter *models.ImageFilter) ([]*models.Image, error) {
	query := r.collection.Query

	if filter.DatasetName != nil && *filter.DatasetName != "" {
		query = query.Where("dataset_name", "==", *filter.DatasetName)
	}
	if filter.OrganType != nil && *filter.OrganType != "" {
		query = query.Where("organ_type", "==", *filter.OrganType)
	}
	if filter.DiseaseType != nil && *filter.DiseaseType != "" {
		query = query.Where("disease_type", "==", *filter.DiseaseType)
	}
	if filter.Classification != nil && *filter.Classification != "" {
		query = query.Where("classification", "==", *filter.Classification)
	}
	if filter.SubType != nil && *filter.SubType != "" {
		query = query.Where("subtype", "==", *filter.SubType)
	}
	if filter.Grade != nil && *filter.Grade != "" {
		query = query.Where("grade", "==", *filter.Grade)
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to filter images: %w", err)
	}

	var images []*models.Image
	for _, doc := range docs {
		var image models.Image
		if err := doc.DataTo(&image); err != nil {
			return nil, fmt.Errorf("failed to convert document to image: %w", err)
		}
		image.ID = doc.Ref.ID // Set the ID from the document reference
		images = append(images, &image)
	}

	return images, nil
}
