# TODO

1. Create cli with Viper and Cobra.

2. Create server initialization that utilizes options functions and the
   following dependencies.
   > MySQL
   > Viper

3. Create MySQL shoppingcart DB schema. Insert multiple items that may be
   purchased by users.

4. Utilize the Chi router to write REST endpoints that access and manipulate
   the shopping cart resource.

   A user may access and modify there cart resource in the following ways...
   > add an item to a user's shopping cart (POST /cart/item/:id?count=:cnt)
   > delete an item from a user's shopping cart (DELETE /cart/item/:id)
   > update an item's count in the shopping cart (DELETE /cart/item/:id?count=:cnt)
   > get a user's shopping cart (GET /cart)

   A user may access the items resource in the following ways...
   > get a list of items (GET /items?max_results=:max_results&page_token=:page_token)

5. Write integration tests for between the server and DB layer.
