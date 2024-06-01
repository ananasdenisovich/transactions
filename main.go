package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI       = "mongodb+srv://ananasovich2002:87787276658Aa.@cluster0.80wl48q.mongodb.net/"
	databaseName   = "paymentDB"
	collectionName = "payments"
)

var client *mongo.Client
var paymentsCollection *mongo.Collection

type Payment struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CartID     string             `json:"cartID"`
	CardNumber string             `json:"cardNumber"`
	ExpiryDate string             `json:"expiryDate"`
	CVV        string             `json:"cvv"`
	Status     string             `json:"status"`
	Timestamp  time.Time          `json:"timestamp"`
}

func init() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	paymentsCollection = client.Database(databaseName).Collection(collectionName)
}

func confirmPayment(c *gin.Context) {
	var payment Payment
	if err := c.BindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment details"})
		return
	}

	isPaymentSuccessful := true

	if isPaymentSuccessful {
		payment.ID = primitive.NewObjectID()
		payment.Timestamp = time.Now()
		payment.Status = "paid"

		_, err := paymentsCollection.InsertOne(context.TODO(), payment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Payment successful", "paymentID": payment.ID})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment failed"})
	}
}

func main() {
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/payment.html", func(c *gin.Context) {
		c.File("./static/payment.html")
	})

	r.POST("/confirm-payment", confirmPayment)

	if err := r.Run(":8081"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}
