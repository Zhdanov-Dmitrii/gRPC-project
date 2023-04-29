package main

import (
	context "context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net"

	pbf "gRPC-project/api"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type Service struct {
	pbf.UnimplementedKeyValueServiceServer
	db *sql.DB
}

func (s Service) FindById(ctx context.Context, request *pbf.GetKeyValueRequest) (*pbf.KeyValue, error) {
	query, err := s.db.Query("SELECT * FROM KeyValue WHERE id = ?", request.Id)
	defer query.Close()
	if err != nil {
		return nil, err
	}

	var id int32
	var val string
	if query.Next() {
		query.Scan(&id, &val)
		return &pbf.KeyValue{Id: id, Val: val}, nil
	}

	return nil, errors.New("no id")
}

func (s Service) Put(ctx context.Context, request *pbf.PutKeyValueRequest) (*pbf.KeyValue, error) {
	_, err := s.db.Exec("INSERT INTO KeyValue(id, val) VALUES(?,?)", request.Id, request.Val)
	if err != nil {
		return nil, err
	}

	return &pbf.KeyValue{Id: request.Id, Val: request.Val}, nil
}

func (s Service) Delete(ctx context.Context, value *pbf.DeleteKeyValue) (*pbf.KeyValue, error) {
	query, err := s.db.Query("DELETE FROM KeyValue WHERE id = ?", value.Id)
	defer query.Close()
	if err != nil {
		return nil, err
	}

	var id int32
	var val string

	query.Next()
	query.Scan(&id, &val)

	return &pbf.KeyValue{Id: id, Val: val}, nil
}

func (s Service) ManyKeyValues(ctx context.Context, request *pbf.PagingRequest) (*pbf.KeyValues, error) {
	query, err := s.db.Query("SELECT id, val FROM KeyValue LIMIT(?) OFFSET (?)",
		request.PageLength, request.PageLength*(request.PageNumber-1))
	defer query.Close()
	if err != nil {
		return nil, err
	}

	var id int32
	var val string
	res := pbf.KeyValues{}

	res.KeyValues = make([]*pbf.KeyValue, 0, 10)
	for query.Next() {
		query.Scan(&id, &val)
		res.KeyValues = append(res.KeyValues, &pbf.KeyValue{Id: id, Val: val})
	}

	return &res, nil
}

func (s Service) mustEmbedUnimplementedKeyValueServiceServer() {
	//TODO implement me
	panic("implement me")
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	db, err := sql.Open("sqlite3", "db.sqlite")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS KeyValue (id INTEGER PRIMARY KEY, val varchar(128));")
	if err != nil {
		panic(err)
	}
	service := Service{db: db}
	pbf.RegisterKeyValueServiceServer(s, &service)

	// Регистрация службы ответов на сервере gRPC.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
}
