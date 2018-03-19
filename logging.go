package main


import(
	"fmt"
	"github.com/RATDistributedSystems/utilities"
	"net"
	"log"
)


func sendMsgToAuditServer(msg string) {
	conn := auditPool.getConnection()
	fmt.Fprintln(conn, msg)
	auditPool.returnConnection(conn)
}

func logQuoteEvent(server string, transactionNum int, price string, stockSymbol string, userid string, quoteservertime string, cryptokey string) {
	msg := fmt.Sprintf("Quote,%s,%s,%d,%s,%s,%s,%s,%s", utilities.GetTimestamp(), server, transactionNum, price, stockSymbol, userid, quoteservertime, cryptokey)
	sendMsgToAuditServer(msg)
}

func GetQuoteServerConnection() net.Conn {
	addr, protocol := configurationServer.GetServerDetails("quote")
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		log.Printf("Encountered error when trying to connect to quote server\n%s", err.Error())
	}
	return conn
}