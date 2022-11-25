import Address from "./address";

export default class Customer {
  private _rewardPoints: number = 0;
  constructor(private _id: string, private _name: string, private _address?: Address, private _active?: boolean) {
    this.validate();
  } 

  get id() {
    return this._id;
  }

  get name() {
    return this._name;
  }

  get active() {
    return this._active;
  }

  get rewardPoints(): number {
    return this._rewardPoints;
  }

  get Address(): Address {
    return this._address;
  }

  set Address(address: Address) {
    this._address = address;
  }

  changeName(name: string) {
    this._name = name;
    this.validate();
  }

  changeAddress(address: Address): void {
    this._address = address;
  }

  activate() {
    if (this._address === undefined) {
      throw new Error("Address is needed to activate a Customer");
    }
    this._active = true;
  }

  isActive(): boolean {
    return this._active;
  }

  addRewardPoints(points: number) {
    this._rewardPoints += points;
  }

  deactivate() {
    this._active = false;
  }

  validate() {

    if(this._active === undefined || this._active === null) {
      this._active = false;
    }

    if (this._name.length === 0) {
      throw new Error("Customer name can not be empty");
    }

    if (this._id.length === 0) {
      throw new Error("Customer id can not be empty");
    }
  }
}