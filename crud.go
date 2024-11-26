package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	gomail "gopkg.in/mail.v2"
	"math/rand"
	"net/http"
	"time"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcde2331312315897963156fghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type UserCRUD struct {
	userCollection *mongo.Collection
}

func (service *UserCRUD) GetUsers(c *gin.Context) {
	var users []*User
	ctx := c.Request.Context()
	filter := bson.D{{}}
	cursor, err := service.userCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	for cursor.Next(ctx) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
			return
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	cursor.Close(ctx)

	if len(users) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (service *UserCRUD) CreateUser(c *gin.Context) {
	var user User
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var userFetch *User
	query := bson.D{bson.E{Key: "email", Value: user.Email}}
	err := service.userCollection.FindOne(ctx, query).Decode(&userFetch)
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error user already exists",
		})
		return
	}
	user.Password = StringWithCharset(6, charset)
	_, insertErr := service.userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": insertErr.Error()})
		return
	}
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", "hincapied07@email.com")
	message.SetHeader("To", user.Email)
	message.SetHeader("Subject", "Registro en SCC exitoso")

	htmlTemplate := `
    <!DOCTYPE html>
    <html lang="es">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Datos de Inicio de Sesión</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f4f4f4;
                color: #333;
                padding: 20px;
            }
            .container {
                background-color: #ffffff;
                border-radius: 8px;
                padding: 20px;
                width: 80%%;
                max-width: 600px;
                margin: 0 auto;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            }
            h2 {
                color: #333;
                text-align: center;
            }
            p {
                font-size: 16px;
            }
            .button {
                display: inline-block;
                background-color: #4CAF50;
                color: #fff;
                padding: 10px 20px;
                text-decoration: none;
                border-radius: 5px;
                text-align: center;
                margin-top: 20px;
            }
            .footer {
                font-size: 12px;
                color: #777;
                text-align: center;
                margin-top: 20px;
            }
            .footer a {
                color: #777;
                text-decoration: none;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h2>¡Bienvenido! Aquí están tus datos de inicio de sesión</h2>
            <p>Estimado/a %s,</p>
            <p>Gracias por registrarte en SCC. A continuación, encontrarás tus datos de inicio de sesión:</p>
            <table>
                <tr>
                    <td><strong>Nombre de usuario:</strong></td>
                    <td>%s</td>
                </tr>
                <tr>
                    <td><strong>Contraseña:</strong></td>
                    <td>%s</td>
                </tr>
            </table>
        </div>
    </body>
    </html>
    `

	// Personalizar la plantilla con los datos específicos
	body := fmt.Sprintf(htmlTemplate, user.Name, user.Email, user.Password)

	// Set email body
	message.SetBody("text/html", body)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "hincapied07@gmail.com", "anea utkx jrbj ixhi")

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"user":     user.Email,
		"password": user.Password,
	})
}

func (service *UserCRUD) Login(c *gin.Context) {
	var user User
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var userFetch *User
	query := bson.D{bson.E{Key: "email", Value: user.Email}}
	err := service.userCollection.FindOne(ctx, query).Decode(&userFetch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error user not exists",
		})
		return
	}
	if userFetch.Email == user.Email && userFetch.Password == user.Password {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "unauthorized",
	})
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
