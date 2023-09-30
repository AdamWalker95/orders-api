<h1>Requirements</h1>
<h3>A copy of Redis is required for running software. If you need to download it you can do so by running the following on a linux terminal:</h3>
<p>sudo apt-get install redis</p>

<h1>REST calls for testing system</h1>
<h2>Adding Records:</h2>
curl -X POST -d '{"customer_id":"'$(uuidgen)'","line_items":[{"item_id":"'$(uuidgen)'","quantity":5,"price":200}]}' localhost:3000/orders

<h2>Retrieving records:</h2>
curl -sS localhost:3000/orders | jq

<h2>Retrieve Specific Orders</h2>
curl -sS localhost:3000/orders/[order_id] | jq

<h2>Updating Records</h2>
<h3>Examples below, adject as you see fit</h3>
<p>curl -X PUT -d '{"status":"shipped"}' -sS "localhost:3000/orders/[order_id]" | jq</p>
<p>curl -X PUT -d '{"status":"completed"}' -sS "localhost:3000/orders/[order_id]" | jq</p>

<h2>Delete Record</h2>
curl -X DELETE localhost:3000/orders/[order_id]

<h1>Commands for using Redis</h1>
Useful Troubleshooting if REST commands are failing

<h2>Find Record on redis by running:</h2>
<p>-redis-cli</p>
<p>-GET "order:[order_id]"</p>

<h2>View all orders by running:</h2>
<p>-redis-cli</p>
<p>-SMEMBERS orders</p>
