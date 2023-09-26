Test commands

Adding Records:
curl -X POST -d '{"customer_id":"'$(uuidgen)'","line_items":[{"item_id":"'$(uuidgen)'","quantity":5,"price":200}]}' localhost:3000/orders

Retrieving records:
curl -sS localhost:3000/orders | jq

Find Record on redis by running:
-redis-cli
-GET "order:[order_id]"

View all orders by running:
-redis-cli
-SMEMBERS orders
