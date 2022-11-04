import OrderItem from "./order_item";

export default class Order {

  private _total: number;
  constructor(private id: string, private customerId: string, private _items: OrderItem[]) {
    this.validate();
  }

  validate() {
    if (this.id.length === 0) {
      throw new Error("id is required");
    }
    if (this.customerId.length === 0) {
      throw new Error("customerId is required");
    }
    if (this._items.length === 0) {
      throw new Error("Order should have at least one item");
    }
  }

  total(): number {
    return this._items.reduce((acc, item) => acc + item.orderItemTotal(), 0);
  }
}