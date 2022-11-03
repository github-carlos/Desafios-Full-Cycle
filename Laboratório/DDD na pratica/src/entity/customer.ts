import Address from "./address";

export default class Customer {
  constructor(private _id: string, private _name: string, private _address?: Address, private _active?: boolean) {
    this.validate();
  } 

  get name() {
    return this._name;
  }

  get active() {
    return this._active;
  }

  changeName(name: string) {
    this._name = name;
    this.validate();
  }

  activate() {
    if (this._address === undefined) {
      throw new Error("Address is needed to activate a Customer");
    }
    this._active = true;
  }

  deactivate() {
    this._active = false;
  }

  validate() {
    if (this._name.length === 0) {
      throw new Error("Customer name can not be empty");
    }

    if (this._id.length === 0) {
      throw new Error("Customer id can not be empty");
    }
  }
}