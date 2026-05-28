package server

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sort"
)

type Server struct {
	Port int
}

type DataPoint []int // [t, p]
type Series []DataPoint

func NewDataPoint(timestamp, price int32) DataPoint {
	return DataPoint{int(timestamp), int(price)}
}

func (ser *Series) Insert(d DataPoint) {
	t := d[0]
	i := sort.Search(len(*ser), func(i int) bool {
		return (*ser)[i][0] >= t
	})
	*ser = append(*ser, nil)
	copy((*ser)[i+1:], (*ser)[i:])
	(*ser)[i] = d
}

func (ser *Series) Query(mintime, maxtime int32) Series {
	var out Series
	for _, p := range *ser {
		if p[0] >= int(mintime) && p[0] <= int(maxtime) {
			out = append(out, p)
		}
	}
	return out
}
func (ser Series) Mean() int32 {
	var all int
	for _,point := range ser {
		all += point[1]
	}
	return int32(all / len(ser))
}

type Session struct {
	conn   net.Conn
	id     int
	series *Series
}

func New(port int) *Server {
	return &Server{
		Port: port,
	}
}

func (s *Session) String() string {
	out := fmt.Sprintf("Session:\nID: %d\n", s.id)
	out += fmt.Sprintf(
		"Connection:\n  Local:  %s\n  Remote: %s\n",
		s.conn.LocalAddr(),
		s.conn.RemoteAddr(),
	)
	out += fmt.Sprintf("Series points: %d\n", len(*s.series))
	for i, p := range *s.series {
		out += fmt.Sprintf("  [%d] t=%d p=%d\n", i, p[0], p[1])
	}
	return out
}

func handleInsert(msg []byte, session *Session) {
	timestamp_bytes := msg[1:5]
	timestamp := int32(binary.BigEndian.Uint32(timestamp_bytes))

	price_bytes := msg[5:]
	price := int32(binary.BigEndian.Uint32(price_bytes))

	session.series.Insert(NewDataPoint(timestamp, price))

}

func handleQuery(msg []byte, session *Session)  {

	mintime_bytes := msg[1:5]
	mintime := int32(binary.BigEndian.Uint32(mintime_bytes))

	maxtime_bytes := msg[5:]
	maxtime := int32(binary.BigEndian.Uint32(maxtime_bytes))


	answer := []byte{0,0,0,0}
	if mintime > maxtime {
	 _,_ = session.conn.Write([]byte(answer))
	 return
	}
	queryResult := session.series.Query(mintime, maxtime)

	if len(queryResult) == 0 {
	 _,_ = session.conn.Write([]byte(answer))
	 return
	}

	binary.BigEndian.PutUint32(answer, uint32(queryResult.Mean()))
	_,_ = session.conn.Write([]byte(answer))
	return
}

func handleClient(conn net.Conn, id int) {
	session := Session{
		conn:   conn,
		id:     id,
		series: &Series{},
	}
	fmt.Println("Started Session: ", id)
	defer fmt.Println("Ended Session: ", id)
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		var msg []byte = make([]byte, 9)
		_, err := io.ReadFull(r, msg)
		if err != nil {
			return
		}

		kind := msg[:1]
		switch kind[0] {
		case 'I':
			handleInsert(msg, &session)
		case 'Q':
			handleQuery(msg, &session)
		default:
			return
		}
	}
}

func acceptLoop(l net.Listener) error {
	id := 1
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go handleClient(conn, id)
		id += 1
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.Port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	err = acceptLoop(l)
	if err != nil {
		return err
	}

	return nil
}
