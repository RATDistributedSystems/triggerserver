package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"bufio"
)

func quoteRequest(userID string, stockSymbol string, transactionNum int) int {
	//logUserEvent("TS1", transactionNum, "QUOTE", userID, stockSymbol, "")

	// Make Quote Request
	conn := GetQuoteServerConnection()
	fmt.Fprintf(conn, "%s,%s\n", stockSymbol, userID)

	// Get response
	message, _ := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
	messageArray := strings.Split(message, ",")
	logQuoteEvent(serverName, transactionNum, messageArray[0], messageArray[1], messageArray[2], messageArray[3], messageArray[4])
	return stringToCents(messageArray[0])
}

func getUsableCash(userId string) int {
	var usableCash int
	if err := sessionGlobalTS.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
		panic(fmt.Sprintf("problem creating session", err))
	}
	return usableCash
}

func stringToCents(x string) int {
	result := strings.Split(x, ".")
	dollars, err := strconv.Atoi(strings.TrimSpace(result[0]))
	if err != nil {
		log.Printf("Couldn't convert %s to int", result[0])
		return 0
	}

	cents, err := strconv.Atoi(strings.TrimSpace(result[1]))
	if err != nil {
		log.Printf("Couldn't convert %s to int", result[1])
		return 0
	}

	return (dollars * 100) + cents
}

//chekcs if the command can be executed
//ie check if buy set before commit etc
func checkDependency(command string, userId string, stock string) bool {

	var count int
	var err error = nil
	switch command {
	case "COMMIT_BUY":
		err = sessionGlobalTS.Query("SELECT count(*) FROM buypendingtransactions WHERE userId='" + userId + "'").Scan(&count)
	case "COMMIT_SELL":
		err = sessionGlobalTS.Query("SELECT count(*) FROM sellpendingtransactions WHERE userId='" + userId + "'").Scan(&count)
	case "CANCEL_BUY":
		sessionGlobalTS.Query("SELECT count(*) FROM buypendingtransactions WHERE userId='" + userId + "'").Scan(&count)
	case "CANCEL_SELL":
		err = sessionGlobalTS.Query("SELECT count(*) FROM sellpendingtransactions WHERE userId='" + userId + "'").Scan(&count)
	case "CANCEL_SET_BUY":
		err = sessionGlobalTS.Query("SELECT count(*) FROM buyTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count)
	case "CANCEL_SET_SELL":
		err = sessionGlobalTS.Query("SELECT count(*) FROM sellTriggers WHERE userid='" + userId + "' AND stock='" + stock + "'").Scan(&count)
	}

	if err != nil {
		panic(err)
	}
	return count != 0
}

func addFunds(userId string, addCashAmount int) {
	usableCash := getUsableCash(userId)
	totalCash := usableCash + addCashAmount
	totalCashString := strconv.FormatInt(int64(totalCash), 10)

	//return add funds to user
	if err := sessionGlobalTS.Query("UPDATE users SET usableCash =" + totalCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
		panic(err)
	}

}

//check if the trigger hasn't been cancelled
func checkTriggerExists(userId string, stock string, isBuyOperation bool) bool {

	var count int

	if isBuyOperation == true {
		if err := sessionGlobalTR.Query("SELECT count(*) FROM buyTriggers WHERE pending=FALSE AND userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil {
			panic(err)
		}
	} else {
		if err := sessionGlobalTR.Query("SELECT count(*) FROM sellTriggers WHERE pending=FALSE AND userid='" + userId + "' AND stock='" + stock + "'").Scan(&count); err != nil {
			panic(err)
		}
	}

	return count == 1
}

func checkStockOwnership(userId string, stock string) (int, string) {
	var ownedstockname string
	var ownedstockamount int
	var usid string
	//var hasStock bool

	iter := sessionGlobalTS.Query("SELECT usid, stock, stockamount FROM userstocks WHERE userid='" + userId + "'").Iter()
	for iter.Scan(&usid, &ownedstockname, &ownedstockamount) {
		if ownedstockname == stock {
			//hasStock = true
			break
		}
	}
	if err := iter.Close(); err != nil {
		panic(err)
	}

	//returns 0 if stock is empty
	return ownedstockamount, usid

}
