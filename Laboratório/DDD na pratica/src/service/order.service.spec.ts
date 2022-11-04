import Order from "../entity/order";
import OrderItem from "../entity/order_item";
import OrderService from "./order.service";

describe("#OrderService", () => {

  it("should return total of all orders", () => {
    const orderItem1 = new OrderItem('id1', 'item1', 13, 'product1', 2);
    const orderItem2 = new OrderItem('id2', 'item2', 39, 'product2', 4);

    const order1 = new Order('id1', 'customer1', [orderItem1]);
    const order2 = new Order('id2', 'customer1', [orderItem1, orderItem2]);

    const total = OrderService.total([order1, order2]);

    expect(total).toBe(208);
  });
})