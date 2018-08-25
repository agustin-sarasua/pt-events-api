package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var ginLambda *ginadapter.GinLambda

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var infoLogger = log.New(os.Stdout, "INFO ", log.Llongfile)

func getClaimsSub(ctx events.APIGatewayProxyRequestContext) string {
	jc, _ := json.Marshal(ctx.Authorizer)
	r := make(map[string]map[string]interface{})
	err := json.Unmarshal([]byte(jc), &r)
	if err != nil {
		fmt.Printf("Something went wrong %v", err)
	}
	fmt.Printf("Printing sub: %s ", r["claims"]["sub"].(string))
	return r["claims"]["sub"].(string)
}

func getEvents(c *gin.Context) {
	//apiGwContext, _ := ginLambda.GetAPIGatewayContext(c.Request)
	//sub := getClaimsSub(apiGwContext)
	placeID := c.Query("placeID")
	startTime := c.DefaultQuery("startTime", time.Now().Format("2006-01-02T15:04:05"))
	endTime := c.DefaultQuery("endTime", time.Now().Add(time.Hour*time.Duration(1)).Format("2006-01-02T15:04:05"))

	updateThingsShadow("deeplens_44a03494-1555-4a6f-871b-205b310b2dba", "", nil)

	items, err := getItems(placeID, startTime, endTime, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(200, items)
}

func createEvent(c *gin.Context) {
	// the methods are available in your instance of the GinLambda
	// object and receive the http.Request object
	apiGwContext, _ := ginLambda.GetAPIGatewayContext(c.Request)
	e := Event{}
	uid, _ := uuid.NewV4()
	e.ID = uid.String()
	//apiGwStageVars, _ := ginLambda.GetAPIGatewayStageVars(c.Request)

	// stage variables are stored in a map[string]string
	// stageVarValue := apiGwStageVars["MyStageVar"]

	err := c.BindJSON(&e)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	_, err = time.Parse("2006-01-02T15:04:05", e.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	_, err = time.Parse("2006-01-02T15:04:05", e.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	if e.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Name not found",
		})
		return
	}
	e.Sub = getClaimsSub(apiGwContext)
	err = putItem(&e)
	if err != nil {
		fmt.Printf("Error saving item in db %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusAccepted, e)
}

func updateThingsShadow(thingName string, topic string, event *Event) {
	sess := session.New(aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))
	iotSvc := iot.New(sess)
	res, err := iotSvc.DescribeEndpoint(&iot.DescribeEndpointInput{})
	if err != nil {
		fmt.Println("describe endpoint err", err)
		return
	}
	fmt.Println("Endpoint:", res)

	svc := iotdataplane.New(sess, &aws.Config{Endpoint: res.EndpointAddress})
	// jc, _ := json.Marshal(event)
	// result, err := svc.Publish(&iotdataplane.PublishInput{
	// 	Payload: jc,
	// 	Topic:   &topic,
	// })
	// fmt.Println("publish", result, err)

	shadow, err := svc.GetThingShadow(&iotdataplane.GetThingShadowInput{
		ThingName: &thingName,
	})
	fmt.Println("shadow", shadow, err)
}

// Handler is the main entry point for Lambda. Receives a proxy request and
// returns a proxy response
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("Gin cold start")
		r := gin.Default()
		// r.GET("/events", getPersons)
		r.POST("/events", createEvent)
		r.GET("/events", getEvents)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}
