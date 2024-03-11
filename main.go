package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tidwall/gjson"
)

type ContextKey string

var templates = map[string]string{
	"thank":                 "Dear %s, Your unwavering support and trust in our products/services mean the world to us. We are truly grateful for the opportunity to serve you and for the strong partnership we have built.",
	"transaction completed": "We sincerely appreciate your business and the trust you have placed in us. If you have any further questions or need assistance in the future, please don't hesitate to reach out. We value your satisfaction and look forward to serving you again.",
}

// 获取正在运行的函数名
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// Handler 处理
type Handler interface {
	// 自身的业务
	Do(ctx context.Context) (context.Context, error)
	// 设置下一个对象
	SetNext(h Handler) Handler
	// 执行
	Run(ctx context.Context) error
}

// Next 抽象出来的 可被合成复用的结构体
type Next struct {
	// 下一个对象
	nextHandler Handler
}

// SetNext 实现好的 可被复用的SetNext方法
// 返回值是下一个对象 方便写成链式代码优雅
// 例如 nullHandler.SetNext(argumentsHandler).SetNext(signHandler).SetNext(frequentHandler)
func (n *Next) SetNext(h Handler) Handler {
	n.nextHandler = h
	return h
}

// Run 执行
func (n *Next) Run(ctx context.Context) (err error) {
	// 由于go无继承的概念 这里无法执行当前handler的Do
	// n.Do(c)
	if n.nextHandler != nil {
		// 合成复用下的变种
		// 执行下一个handler的Do
		var nctx context.Context
		if nctx, err = (n.nextHandler).Do(ctx); err != nil {
			return
		}
		// 执行下一个handler的Run
		return (n.nextHandler).Run(nctx)
	}
	return
}

// NullHandler 空Handler
// 由于go无继承的概念 作为链式调用的第一个载体 设置实际的下一个对象
type NullHandler struct {
	// 合成复用Next的`nextHandler`成员属性、`SetNext`成员方法、`Run`成员方法
	Next
}

// Do 空Handler的Do
func (n *NullHandler) Do(ctx context.Context) (err error) {
	// 空Handler 这里什么也不做 只是载体 do nothing...
	return
}

// ArgumentsHandler 校验参数的handler
type ArgumentsHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (a *ArgumentsHandler) Do(ctx context.Context) (context.Context, error) {
	log.Println("校验参数")
	event := new(FacebookEvent)
	fbReq := ctx.Value("fbReq").(events.APIGatewayProxyRequest)

	log.Println(fbReq.Body)

	if gjson.Get(fbReq.Body, "entry.0.messaging.0.message.is_echo").Exists() {
		return ctx, errors.New("ignore echo messages")
	}

	event.Field = "messages"
	event.Id = fbReq.RequestContext.RequestID
	event.Message = gjson.Get(fbReq.Body, "entry.0.messaging.0.message.text").String()
	event.CustomerId = gjson.Get(fbReq.Body, "entry.0.messaging.0.sender.id").String()

	return context.WithValue(ctx, ContextKey("fbEvent"), *event), nil
}

// TemplateHandler 模版化生成回复的handler
type TemplateHandler struct {
	// 合成复用Next
	Next
}

// Do 处理事件的逻辑
func (t *TemplateHandler) Do(ctx context.Context) (context.Context, error) {
	log.Println("模版化生成回复")
	fbEvent := ctx.Value("fbEvent").(FacebookEvent)
	message := ""
	if fbEvent.Field == "messages" {
		if strings.Contains(fbEvent.Message, "thank") {
			message = fmt.Sprintf(templates["thank"], fbEvent.CustomerId)
		}
		if strings.Contains(fbEvent.Message, "transaction completed") {
			message = templates["transaction completed"]
		}
	}
	if message == "" || len(message) == 0 {
		return ctx, errors.New("no need to reply")
	}

	return context.WithValue(ctx, ContextKey("fbReply"), message), nil
}

// MessageSender 发送消息的handler
type MessageSender struct {
	// 合成复用Next
	Next
}

// Do 发送消息的逻辑
func (m *MessageSender) Do(ctx context.Context) (context.Context, error) {
	log.Println("发送回复")
	fbReply := ctx.Value("fbReply").(string)
	fbEvent := ctx.Value("fbEvent").(FacebookEvent)
	if err := sendMessage(fbReply, fbEvent.CustomerId, "RESPONSE"); err != nil {
		return ctx, err
	}

	return ctx, nil
}

// StorageHandler 存储内容的handler
type StorageHandler struct {
	// 合成复用Next
	Next
}

// Do 存储内容的逻辑
func (s *StorageHandler) Do(ctx context.Context) (context.Context, error) {
	log.Println("存储内容")
	fbReply := ctx.Value("fbReply").(string)
	fbEvent := ctx.Value("fbEvent").(FacebookEvent)
	var stor Storage = new(Ddb)
	if err := stor.store(fbEvent, fbReply); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func botHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if event.HTTPMethod == "GET" {
		return verifyHandler(ctx, event)
	}

	// 初始化空handler
	nullHandler := &NullHandler{}
	nullHandler.SetNext(&ArgumentsHandler{}).
		SetNext(&TemplateHandler{}).
		SetNext(&MessageSender{}).
		SetNext(&StorageHandler{})
	// 开始执行业务
	rootCtx := context.Background()
	if err := nullHandler.Run(context.WithValue(rootCtx, ContextKey("fbReq"), event)); err != nil {
		// 异常
		log.Println("Fail | Error:" + err.Error())
	}

	var response events.APIGatewayProxyResponse
	response = events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       event.Body,
	}
	return response, nil
}

func verifyHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var queryParameters = event.QueryStringParameters
	var response events.APIGatewayProxyResponse
	if queryParameters["hub.verify_token"] == "CONGYAO_VERIFIY_TOKEN" && queryParameters["hub.mode"] == "subscribe" {
		response = events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       queryParameters["hub.challenge"],
		}
	} else {
		response = events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "Missing/invalid token",
		}
	}

	return response, nil
}

func main() {
	lambda.Start(botHandler)
}
