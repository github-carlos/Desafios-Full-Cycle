import OrderItem from "./order_item";

export default class Order {

  get id() {
    return this._id;
  }

  get customerId(): string {
    return this._customerId;
  }

  get items(): Array<OrderItem> {
    return this._items;
  }

  constructor(private _id: string, private _customerId: string, private _items: OrderItem[]) {
    this.validate();
  }

  validate() {
    if (this._id.length === 0) {
      throw new Error("id is required");
    }
    if (this._customerId.length === 0) {
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