package store

import (
	"context"

	"github.com/filedrive-team/go-ds-cluster/core"
	ds "github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type server struct {
	ctx      context.Context
	host     host.Host
	protocol protocol.ID
	ds       ds.Datastore
}

func NewStoreServer(ctx context.Context, h host.Host, pid protocol.ID, ds ds.Datastore) core.DataNodeServer {
	return &server{
		ctx:      ctx,
		host:     h,
		protocol: pid,
		ds:       ds,
	}
}

func (sv *server) Close() error {
	return sv.host.Close()
}

func (sv *server) Serve() {
	logging.Info("data node server set stream handler")
	sv.host.SetStreamHandler(sv.protocol, sv.handleStream)
}

func (sv *server) handleStream(s network.Stream) {
	defer s.Close()
	logging.Info("server incoming stream")
	reqMsg := new(RequestMessage)

	if err := ReadRequestMsg(s, reqMsg); err != nil {
		logging.Error(err)
		return
	}

	logging.Infof("req action %v", reqMsg.Action)
	switch reqMsg.Action {
	case ActGet:
		sv.get(s, reqMsg)
	case ActGetSize:
		sv.getSize(s, reqMsg)
	case ActHas:
		sv.has(s, reqMsg)
	case ActPut:
		sv.put(s, reqMsg)
	case ActDelete:
		sv.delete(s, reqMsg)
	default:
		logging.Warnf("unhandled action: %v", reqMsg.Action)
	}
}

func (sv *server) put(s network.Stream, req *RequestMessage) {
	logging.Infof("put %s, value size: %d", req.Key, len(req.Value))
	res := &ReplyMessage{}
	if err := sv.ds.Put(ds.NewKey(req.Key), req.Value); err != nil {
		res.Code = ErrOthers
		res.Msg = err.Error()
	}
	//res.Msg = "ok"
	if err := WriteReplyMsg(s, res); err != nil {
		logging.Error(err)
	}
}

func (sv *server) has(s network.Stream, req *RequestMessage) {
	res := &ReplyMessage{}
	exists, err := sv.ds.Has(ds.NewKey(req.Key))
	if err != nil {
		if err == ds.ErrNotFound {
			res.Code = ErrNotFound
		} else {
			res.Code = ErrOthers
		}
		res.Msg = err.Error()
	} else {
		res.Exists = exists
	}
	if err := WriteReplyMsg(s, res); err != nil {
		logging.Error(err)
	}
}

func (sv *server) getSize(s network.Stream, req *RequestMessage) {
	res := &ReplyMessage{}
	size, err := sv.ds.GetSize(ds.NewKey(req.Key))
	if err != nil {
		if err == ds.ErrNotFound {
			res.Code = ErrNotFound
		} else {
			res.Code = ErrOthers
		}
		res.Msg = err.Error()
	} else {
		res.Size = int64(size)
	}
	if err := WriteReplyMsg(s, res); err != nil {
		logging.Error(err)
	}
}

func (sv *server) get(s network.Stream, req *RequestMessage) {
	res := &ReplyMessage{}
	v, err := sv.ds.Get(ds.NewKey(req.Key))
	if err != nil {
		if err == ds.ErrNotFound {
			res.Code = ErrNotFound
		} else {
			res.Code = ErrOthers
		}
		res.Msg = err.Error()
	} else {
		res.Value = v
	}
	if err := WriteReplyMsg(s, res); err != nil {
		logging.Error(err)
	}
}

func (sv *server) delete(s network.Stream, req *RequestMessage) {
	res := &ReplyMessage{}
	err := sv.ds.Delete(ds.NewKey(req.Key))
	if err != nil {
		res.Code = ErrOthers
		res.Msg = err.Error()
	}
	if err := WriteReplyMsg(s, res); err != nil {
		logging.Error(err)
	}
}