<h1>Test commands</h1>

<h2>Adding Records:</h2>
curl -X POST -d '{"customer_id":"'$(uuidgen)'","line_items":[{"item_id":"'$(uuidgen)'","quantity":5,"price":200}]}' localhost:3000/orders

<h2>Retrieving records:</h2>
curl -sS localhost:3000/orders | jq

<h2>Find Record on redis by running:</h2>
-redis-cli
-GET "order:[order_id]"

<h2>View all orders by running:</h2>
-redis-cli
-SMEMBERS orders
