<h1>Summary of Project</h1>
The software on this repository is used for database management. The software accesses a local database (redis) and is used to add, retrieve, update, and delete records via REST API. At the moment, the database holds two types of records; orders, and users.

<h1>Requirements</h1>
<h3>A copy of Redis is required for running software. If you need to download it you can do so by running the following on a linux terminal:</h3>
<p>sudo apt-get install redis</p>
<p>After this you will need three terminals:</p>
<p>(1). For running the database [run 'redis-server']</p>
<p>(2). For running the program [run 'go run main.go']</p>
<p>(3). For running the REST commands, below are many examples that can be used from</p>

<h1>REST calls for testing system</h1>
<h2>Add User</h2>
curl -X POST -d '{"email":"email@email.com","password":"password123"}' localhost:3000/user

<h2>Retrieve User</h2>
curl -sS localhost:3000/user/[email] | jq

<h2>Update User's Password</h2>
curl -X PUT -d '{"password":"password456"}' -sS "localhost:3000/user/email@email.com" | jq

<h2>Delete User</h2>
curl -X DELETE localhost:3000/user/email@email.com

<h2>Add Orders:</h2>
curl -X POST -d '{"customer_id":[{"username":"bob"}],"line_items":[{"item_id":"'$(uuidgen)'","quantity":5,"price":200}]}' localhost:3000/orders

<h2>Retrieving Orders:</h2>
curl -sS localhost:3000/orders | jq

<h2>Retrieve Specific Orders</h2>
curl -sS localhost:3000/orders/[order_id] | jq

<h2>Updating Orders</h2>
<h3>Examples below</h3>
<p>curl -X PUT -d '{"status":"shipped"}' -sS "localhost:3000/orders/[order_id]" | jq</p>
<p>curl -X PUT -d '{"status":"completed"}' -sS "localhost:3000/orders/[order_id]" | jq</p>

<h2>Delete Order</h2>
curl -X DELETE localhost:3000/orders/[order_id]

<h1>Commands for using Redis</h1>
Useful Troubleshooting if REST commands are failing

<h2>Find Record on redis by running:</h2>
<p>-redis-cli</p>
<p>-GET "order:[order_id]"</p>

<h2>View all records by running:</h2>
<p>-redis-cli</p>
<p>-SMEMBERS orders (for orders)</p>
<p>-SMEMBERS users (for users)</p>
