package queue

import (
    "context"
    "fmt"
    "time"
    
    amqp "github.com/rabbitmq/amqp091-go"
    "github.com/sirupsen/logrus"
)

type RabbitMQ struct {
    conn         *amqp.Connection
    ch           *amqp.Channel
    uri          string
    consumers    map[string]chan bool
    isConnected  bool
    logger       *logrus.Logger
}

func NewRabbitMQ(user, password, host, port string, logger *logrus.Logger) (*RabbitMQ, error) {
    uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
    rmq := &RabbitMQ{
        uri:       uri,
        consumers: make(map[string]chan bool),
        logger:    logger,
    }
    
    go rmq.reconnectLoop()
    
    return rmq, nil
}

func (r *RabbitMQ) reconnectLoop() {
    for {
        if !r.isConnected {
            if err := r.connect(); err != nil {
                r.logger.Errorf("Failed to connect to RabbitMQ: %v", err)
                time.Sleep(5 * time.Second)
                continue
            }
        }
        time.Sleep(10 * time.Second)
    }
}

func (r *RabbitMQ) connect() error {
    conn, err := amqp.Dial(r.uri)
    if err != nil {
        return err
    }

    ch, err := conn.Channel()
    if err != nil {
        return err
    }

    r.conn = conn
    r.ch = ch
    r.isConnected = true

    go func() {
        <-r.conn.NotifyClose(make(chan *amqp.Error))
        r.logger.Info("RabbitMQ connection lost")
        r.isConnected = false
    }()

    return nil
}

func (r *RabbitMQ) CreateQueue(clientID string) error {
    queueName := fmt.Sprintf("%s.process", clientID)
    
    _, err := r.ch.QueueDeclare(
        queueName,
        true,  
        false, 
        false,
        false,
        nil,  
    )
    
    return err
}

func (r *RabbitMQ) StartConsumer(clientID string) error {
    queueName := fmt.Sprintf("%s.process", clientID)
    stopChan := make(chan bool)
    r.consumers[clientID] = stopChan

    go func() {
        defer func() {
            if err := recover(); err != nil {
                r.logger.Errorf("Consumer panic recovered: %v", err)
            }
        }()

        msgs, err := r.ch.Consume(
            queueName,
            "",    
            true,  
            false, 
            false,
            false, 
            nil,  
        )
        if err != nil {
            r.logger.Errorf("Failed to start consumer: %v", err)
            return
        }

        for {
            select {
            case msg := <-msgs:
                r.logger.Infof("Received message for tenant %s: %s", clientID, string(msg.Body))
            case <-stopChan:
                r.logger.Infof("Stopping consumer for tenant %s", clientID)
                return
            }
        }
    }()

    return nil
}

func (r *RabbitMQ) StopConsumer(clientID string) {
    if stopChan, exists := r.consumers[clientID]; exists {
        close(stopChan)
        delete(r.consumers, clientID)
    }
}

func (r *RabbitMQ) PublishMessage(clientID string, payload []byte) error {
    queueName := fmt.Sprintf("%s.process", clientID)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return r.ch.PublishWithContext(ctx,
        "",       // exchange
        queueName, // routing key
        false,    // mandatory
        false,    // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:       payload,
        },
    )
} 