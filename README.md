<h1>Summary of Project</h1>
The software on this repository is used for database management. The software accesses a local database and is used to add, retrieve, update, and delete records via REST API. At the moment, the database holds two types of records; orders, and users. The intention is that once complete this will be used to pair customers with orders, and these records can be accessed via a frontend. For extended details of changes/commits added to this respository, please check the 'commit_notes' directory.

<h1>Requirements</h1>
<p>A copy of both MySQL and Redis are required for running software. If you need to download both packages please see below:</p>
<h3>MySQL</h3>
<p>Please run the following on a linux terminal:</p>
<p>-`sudo apt update`</p>
<p>-`sudo apt install mysql-server`</p>
<p>To view status of MySQL (i.e. if it is running) use `sudo systemctl status mysql`</p>
<p>After this, access MySQL by running `mysql -u root -p` (entering your password after `-p`), and then create a new user for MySQL by running `CREATE USER '[user]'@'localhost' IDENTIFIED BY '[password]';` (note "user" and "password" is what's set in the config file of application/config.go for this project)</p>
<p>Note that issues with running MySQL on wsl related to systemd are commonly caused by running an older version of wsl, and can be resolved with updating wsl.</p>
<h3>Redis</h3>
<p>Please run `sudo apt-get install redis` to install</p>
<p>Depending on system settings, you may need to start Redis everytime before use by running `redis-server`</p>

<h1>Running software</h1>
<p></p>For running the program, run `go run main.go` on a terminal</p>
<p>Then access database data by running any of the below curl commands on a second terminal:</p>

<h2>Curl Commands for Users</h2>
<p>Add user:</p>
<p>`curl -X POST -d '{"email":"email@email.com","password":"password123"}' localhost:3000/user`</p>

<p>Retrieve user data:</p>
<p>`curl -sS localhost:3000/user/[email] | jq`</p>

<p>Update user's Password:</p>
<p>`curl -X PUT -d '{"password":"password456"}' -sS "localhost:3000/user/email@email.com" | jq`</p>

<p>Delete user</p>
<p>`curl -X DELETE localhost:3000/user/email@email.com`</p>

<h2>Curl Commands for Orders</h2>
<p>Add Orders:</p>
<p>`curl -X POST -d '{"customer_id":[{"username":"bob"}],"line_items":[{"item_id":"'$(uuidgen)'","quantity":5,"price":200}]}' localhost:3000/orders`</p>

<p>Retrieving all orders:</p>
<p>`curl -sS localhost:3000/orders | jq`</p>

<p>Retrieve specific order:</p>
<p>`curl -sS localhost:3000/orders/[order_id] | jq`</p>

<p>Updating order's status</p>
<p>`curl -X PUT -d '{"status":"shipped"}' -sS "localhost:3000/orders/[order_id]" | jq`</p>
<p>`curl -X PUT -d '{"status":"completed"}' -sS "localhost:3000/orders/[order_id]" | jq`</p>

<p>Delete order</p>
<p>`curl -X DELETE localhost:3000/orders/[order_id]`</p>

<h1>Troubleshooting</h1>
Here are some useful tips if there are issues running software

<p>Find Record on redis by running:</p>
<p>Run `redis-cli`</p>
<p>Then run `GET "order:[order_id]"`</p>
<p>Alternatively, get all orders by running `SMEMBERS orders (for orders)`</p>
