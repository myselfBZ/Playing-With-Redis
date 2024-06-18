package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Book struct{
    name string
    id string
    author string 
}


var ctx = context.Background()

var Client *redis.Client

func init(){
    Client  = InitRedis()
}

func main(){
    books  := []Book{
        {name: "war and peice", id: "1", author: "no one"},
        {name: "war and peice", id: "5",author: "no one"},
        {name: "it ends with us ", id: "4", author: "Some one"},
    }

    Preload(books)

    book, err := GetBook("5")
    if err != nil{
        panic(err)
    } 
    fmt.Println(book)

    
}


func InitRedis() *redis.Client {
    Client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB: 0,
        Password: "",
    })
    _, err := Client.Ping(ctx).Result()
    if err != nil{
        panic("You suck at initializing Redis client")
    }
    return Client 
}


func serializer(book Book) map[string]string{
    return map[string]string{
        "id":book.id,
        "name":book.name,
        "author":book.author,
    }
}


func Preload(books []Book){
    for _, b := range books{
        serialized := serializer(b)
        err := Client.HSet(ctx, b.id, serialized).Err()

        if err != nil{
            fmt.Println("error loading a book")
        }


    }
}

func deserializer(book map[string]string) *Book{
    return &Book{
        name: book["name"],
        author: book["author"],
        id: book["id"],
    }
}

func GetBook(id string) (*Book, error){
    cachedBook, err := Client.HGetAll(ctx,id).Result()
    if err != nil{
        fmt.Println("error fetching that mf book")
        return nil, err 
    }
    return deserializer(cachedBook), nil
}















