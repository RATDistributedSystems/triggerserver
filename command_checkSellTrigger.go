package main

import (
	"fmt"
	//"bufio"
	//"strings"
	//"strconv"
	//"net"
	//"github.com/gocql/gocql"
	"strconv"
	"time"
)


func checkSellTriggers() {
	for{

		fmt.Println("Checking Sell Triggers");
		timer1 := time.NewTimer(time.Millisecond * 500)
		<-timer1.C

		var userid string
		var pendingcash int
		var triggervalue int
		var stock string
		var transactionNum int

		//---------------------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------
		//-- probably need to store the transaction number in to the database so it can be used when proessing other requests
		//---------------------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------


		//check if user currently owns any of this stock
		iter := sessionGlobalTR.Query("SELECT userid, pendingcash, triggerValue, stock FROM sellTriggers WHERE pending=TRUE").Iter()
		for iter.Scan(&userid, &pendingcash, &triggervalue, &stock) {


			//delete record
			if err := sessionGlobalTR.Query("DELETE FROM sellTriggers WHERE pending=TRUE AND userid='" + userid + "' AND stock ='" + stock + "'").Exec(); err != nil {
				panic(fmt.Sprintf("Problem DELETING pending buy trigger", err))
			}
			//set record to not pending
			pendingcashstring := strconv.FormatInt(int64(pendingcash), 10)
			triggervaluestring := strconv.FormatInt(int64(triggervalue), 10)
			if err := sessionGlobalTR.Query("INSERT INTO sellTriggers (pending, userid, stock, pendingcash, triggervalue) VALUES (FALSE ,'" + userid + "','" + stock + "'," + pendingcashstring + "," + triggervaluestring + ")").Exec(); err != nil {
				panic(fmt.Sprintf("Problem INSERTING pending buy trigger", err))
			}

			//process the buy trigger
			go processSellTrigger(userid, stock, triggervalue, transactionNum)


		}
		if err := iter.Close(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

		fmt.Println("Done Checking Sell Triggers");
	}


}

func processSellTrigger(userId string, stock string, stockSellPriceCents int, transactionNum int){

	operation := false

	for {
		//check the quote server every 10 milliseconds
		timer1 := time.NewTimer(time.Millisecond * 500)
		<-timer1.C

		//if the trigger doesnt exist exit

		exists := checkTriggerExists(userId, stock, operation)
		if exists == false {
			return
		}

		//retrieve current stock price
		currentStockPrice := quoteRequest(userId, stock, transactionNum)

		if currentStockPrice > stockSellPriceCents {

			//sell the allocated stocks

			//get stocks allocated to sell
			var pendingStocks int
			if err := sessionGlobalTR.Query("SELECT stockAmount FROM sellTriggers WHERE pending=FALSE AND userid='" + userId + "' AND stock='" + stock + "' ").Scan(&pendingStocks); err != nil {
				//panic(fmt.Sprintf("Problem inputting to Triggers Table", err))
				return
			}

			sellProfits := pendingStocks * currentStockPrice

			//delete pending transaction
			if err := sessionGlobalTR.Query("DELETE FROM sellTriggers WHERE pending=FALSE AND userid='" + userId + "' AND stock='" + stock + "' ").Exec(); err != nil {
				//panic(err)
				return
			}

			//add profits from selling stock to account
			fmt.Println("Sell Trigger Sucessful, profits added to account")
			addFunds(userId, sellProfits)
			return
		}

	}

}
