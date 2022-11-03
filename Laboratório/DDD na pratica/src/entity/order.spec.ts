import Order from "./order";
import OrderItem from "./order_item";

describe("#Order", () => {

  it("should throw error when id is empty", () => {
    expect(() => new Order("", "123", [])).toThrowError();
  });

  it("should throw error when customerId is empty", () => {
    expect(() => new Order("123", "", [])).toThrowError();
  });

  it("should throw error when item list is zero", () => {
    expect(() => new Order("123", "123", [])).toThrowError();
  });

  it("should calculate total correctly when total called with 1 item", () => {
    const item = new OrderItem("123", "Item 1", 100);
    const order = new Order("123", "123", [item]);
    const total = order.total();
    expect(total).toBe(item.price);
  });
});