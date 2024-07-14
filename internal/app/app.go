package app

import (
	"context"
	"path/filepath"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/core-go/firestore/export"
	f "github.com/core-go/io/formatter"
	w "github.com/core-go/io/writer"
)

type ApplicationContext struct {
	Export func(ctx context.Context) (int64, error)
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	opts := option.WithCredentialsJSON([]byte(cfg.Credentials))
	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	formatter, err := f.NewFixedLengthFormatter[User]()
	if err != nil {
		return nil, err
	}
	writer, err := w.NewFileWriter(GenerateFileName)
	if err != nil {
		return nil, err
	}
	exporter := export.NewExporter[User](client.Collection("userimport"), BuildQuery, formatter.Format, writer.Write, writer.Close, "CreateTime", "UpdateTime")

	return &ApplicationContext{
		Export: exporter.Export,
	}, nil
}

type User struct {
	Id          string     `json:"id" gorm:"column:id;primary_key" bson:"_id" format:"%011s" length:"11" dynamodbav:"id" firestore:"id" validate:"required,max=40"`
	Username    string     `json:"username" gorm:"column:username" bson:"username" length:"10" dynamodbav:"username" firestore:"username" validate:"required,username,max=100"`
	Email       *string    `json:"email" gorm:"column:email" bson:"email" dynamodbav:"email" firestore:"email" length:"31" validate:"email,max=100"`
	Phone       string     `json:"phone" gorm:"column:phone" bson:"phone" dynamodbav:"phone" firestore:"phone" length:"20" validate:"required,phone,max=18"`
	DateOfBirth *time.Time `json:"dateOfBirth" gorm:"column:dateOfBirth" bson:"dateOfBirth" length:"10" format:"dateFormat:2006-01-02" dynamodbav:"dateOfBirth" firestore:"dateOfBirth" avro:"dateOfBirth"`
	CreateTime  *time.Time `json:"createTime" gorm:"column:create_time" length:"10" format:"dateFormat:2006-01-02" bson:"createTime" dynamodbav:"createTime" firestore:"-"`
	UpdateTime  *time.Time `json:"updateTime" gorm:"column:update_time" length:"10" format:"dateFormat:2006-01-02" bson:"updateTime" dynamodbav:"updateTime" firestore:"-"`
}

func BuildQuery(ctx context.Context, collection *firestore.CollectionRef) *firestore.DocumentIterator {
	iter := collection.Documents(ctx)
	return iter
}
func GenerateFileName() string {
	fileName := time.Now().Format("20060102150405") + ".csv"
	fullPath := filepath.Join("export", fileName)
	w.DeleteFile(fullPath)
	return fullPath
}
