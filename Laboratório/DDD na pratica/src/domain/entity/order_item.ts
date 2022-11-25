
export default class OrderItem {
  constructor(public id: string,
     public name: string, private _price: number, private _productId: string, private _quantity: number) {}

  get price(): number {
    return this._price;
  }

  get quantity(): number {
    return this._quantity;
  }

  get productId(): string {
    return this._productId;
  }

  orderItemTotal(): number {
    return this._price * this._quantity;
  }
}