Technical Assignment: HDFC

Please find the below two services and the operations that are allowed.

Product Service: provides information about the product like availability, price, category

Order service: provides information about the order like orderValue, dispatchDate, orderStatus, prodQuantity

The user should be able to get the product catalogue and using that info should be able to place an order.

Once the order is placed for a particular product, the product catalogue should be updated accordingly.
(Max quantity of a particular product that can be ordered is 10)
If the order contains 3 premium different products, order value should be discounted by 10%

The Order service should be able to update the orderStatus for a particular order.
dispatchDate should be populated only when the orderStatus is 'Dispatched'.


product category values: Premium/Regular/Budget
order status values: Placed/Dispatched/Completed/Cancelled
