## istio_grpc

```
proxy -> cache -> store
```

**service api design**

```
service CallProxy
{
    rpc GetName(NameReq) returns(NameResp);
    rpc GetDataFromCache(DataModel) returns(DataModel);
    rpc GetDataFromStore(DataModel) returns(DataModel);
    rpc SetDataToCache(DataModel) returns(DataModel);
    rpc SetDataToStore(DataModel) returns(DataModel);

    // get from cache, if cache null, continue get store
    rpc GetData(DataModel) returns(DataModel);

    // delete cache, update store
    rpc SetData(DataModel) returns(DataModel);
}

service CallCache
{
    rpc GetName(NameReq) returns(NameResp);
    rpc GetData(DataModel) returns(DataModel);
    rpc GetDataTrace(DataModel) returns(DataModel);
    rpc SetData(DataModel) returns(DataModel);
    rpc DeleteData(DataModel) returns(DataModel);
}

service CallStore
{
    rpc GetName(NameReq) returns(NameResp);
    rpc GetData(DataModel) returns(DataModel);
    rpc SetData(DataModel) returns(DataModel);
    rpc DeleteData(DataModel) returns(DataModel);
}
```

[call.proto](https://github.com/rfyiamcool/istio_grpc/blob/master/proto/call.proto)


### Dep

```
go get github.com/mattn/goreman
```

### Local Usage

1. protoc

```
make gen
```

2. build golang file

```
make build
```

3. run proxy/cache/store

```
make run
```

4. use client to test

```
./bin/client -debug=true
```

### istio Usage

1. docker build

```
make docker-build
```

2. apply istio

```
to do
```

3. client test

```
to do
```