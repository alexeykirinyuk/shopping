package api

import (
	"context"
	"github.com/alexeykirinyuk/shopping/grpc-libs/products"
	"github.com/alexeykirinyuk/shopping/products/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	databaseName     = "products"
	collectionName   = "products"
)

type ProductService struct {
	UnimplementedProductServiceServer
	client *mongo.Client
}

func CreateProductService(client *mongo.Client) *ProductService {
	return &ProductService{
		client: client,
	}
}

func (c *ProductService) getCollection() *mongo.Collection  {
	return c.client.Database(databaseName).Collection(collectionName)
}

func (c *ProductService) Create(ctx context.Context, in *Product) (*ID, error) {
	product := model.Product{
		Name:        in.GetName(),
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		Price:       in.GetPrice(),
		Currency:    in.GetCurrency(),
	}
	insertResult, err := c.getCollection().InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	id := insertResult.InsertedID.(primitive.ObjectID).Hex()
	return &ID{ID: id}, nil
}

func (c *ProductService) Get(ctx context.Context, in *ID) (*Product, error) {
	id, err := primitive.ObjectIDFromHex(in.GetID())
	if err != nil {
		return nil, err
	}

	var result model.Product
	err = c.getCollection().FindOne(ctx, bson.D{{"id", id}}).Decode(&result)
	if err != nil {
		return nil, err
	}

	product := Product{
		ID:          result.ID.Hex(),
		Name:        result.Name,
		Title:       result.Title,
		Description: result.Description,
		Price:       result.Price,
		Currency:    result.Currency,
	}

	return &product, nil
}

func (c *ProductService) Update(ctx context.Context, in *Product) (*Empty, error) {
	id, err := primitive.ObjectIDFromHex(in.GetID())
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id}
	update := bson.M{
		"name":        in.GetName(),
		"title":       in.GetTitle(),
		"description": in.GetDescription(),
		"price":       in.GetPrice(),
		"currency":    in.GetCurrency(),
	}

	_, err = c.getCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &Empty{}, err
}

func (c *ProductService) Delete(ctx context.Context, in *ID) (*Empty, error) {
	id, err := primitive.ObjectIDFromHex(in.GetID())
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id}

	_, err = c.getCollection().DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &Empty{}, nil
}
