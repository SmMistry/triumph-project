# triumph-project
Crypto Sample Project for Triumph

## Requirements
You muust have a machine with a fairly recent version of Git and Go installed, and some basic understanding of command line opperation

## Get Source and Start running the server

Open a new command window, navigate to the directory where you would like to checkout the project.

Run the following command to download the source to your machine 
*(Note: the clone command will create a new directory called triumph-project in whatever directory you run the clone command from)*

	git clone https://github.com/SmMistry/triumph-project.git

this should retrieve the project code, enter the newly retreived directory:

	cd triumph-project

to compile and run the server simply run:

	go run .

## Calling the server

You may access the server by either opening a browser or using curl on the command line:

### Browser Method

**buy endpoint:**
navigate to: http://localhost:4000/buy?amount=1&symbol=BTC

**Sample response:**
> {"amount":1,"coin":"BTC","exchange":["coinbase"],"usdAmount":76526.31}

**sell endpoint:**
navigate to: http://localhost:4000/sell?amount=0.5&symbol=ETH

**Sample Response:**
> {"amount":0.5,"coin":"ETH","exchange":["coinbase"],"usdAmount":1476.065}

### curl Method

**buy endpoint:**
	curl 'http://localhost:4000/buy?amount=1&symbol=BTC'

**Sample response:**
>{"amount":1,"coin":"BTC","exchange":["coinbase"],"usdAmount":76526.31}

**sell endpoint:**
	curl 'http://localhost:4000/sell?amount=0.5&symbol=ETH'

**Sample Response:**
>{"amount":0.5,"coin":"ETH","exchange":["coinbase"],"usdAmount":1476.065}

### Supported Parameters

**amount:** supports any 64 bit float value

**symbol:** supports any tradeable token available on either coinbase or kraken

**Example symbols:**
BTC
ETH
DOGE
SOL
SHIB

## Running Tests

If you still have the server running you can use (ctrl)+C to terminate the running server.

From the same command prompt in the project directory (triumph-project) runn the following command:

	go test

You should see the following output if all tests are passing:

>2024/11/08 15:05:28 failed to get price from exchange: coinbase error
>2024/11/08 15:05:28 failed to get price from exchange: coinbase error
>2024/11/08 15:05:28 failed to get price from exchange: kraken error
>2024/11/08 15:05:28 failed to get price from exchange: coinbase error
>2024/11/08 15:05:28 failed to get price from exchange: coinbase error
>2024/11/08 15:05:28 failed to get price from exchange: kraken error
>PASS
>ok  	github.com/SmMistry/triumph-project	0.327s


