
export default class OrderItem {
  constructor(public id: string,
     public name: string, private _price: number, private _productId: string, private _quantity: number) {}

  get price(): number {
    return this._price;
  }

  orderItemTotal(): number {
    return this._price * this._quantity;
  }
}