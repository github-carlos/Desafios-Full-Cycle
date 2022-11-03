

export default class Product {

  get name() {
    return this._name;
  }

  get price() {
    return this._price;
  }

  constructor(private id: string, private _name: string, private _price: number) {
    this.validate();
  }

  validate() {
    if (this.id.length === 0) {
      throw new Error("id is required");
    }
    if (this._name.length === 0) {
      throw new Error("name is required");
    }
    if (!this._price || this._price < 0) {
      throw new Error("Price is not valid");
    }
  }

  changeName(newName: string) {
    this._name = newName;
    this.validate();
  }

  changePrice(newPrice: number) {
    this._price = newPrice;
    this.validate();
  }
}