import Address from "./address";
import Customer from "./customer";

describe('#Customer', () => {
  it('should throw error when id is empty', () => {
    expect(() => new Customer('', 'Carlos', new Address('123', 10, '1234', 'city'), true)).toThrowError();
  });
  it('should throw error when name is empty', () => {
    expect(() => new Customer('123', '', new Address('123', 10, '1234', 'city'), true)).toThrowError();
  });

  it("should change name", () => {
    const customer = new Customer("123", "Carlos 1", new Address('123', 10, '1234', 'city'), true);
    const newName = "Carlos 2";
    customer.changeName(newName);
    expect(customer.name).toBe(newName);
  });

  it("should throw error when changed name is invalid", () => {
    const customer = new Customer("123", "Carlos 1", new Address('123', 10, '1234', 'city'), true);
    expect(() => customer.changeName('')).toThrowError();
  });

  it("should active customer", () => {
    const customer = new Customer("123", "Carlos 1", new Address('123', 10, '1234', 'city'), false);
    customer.activate();
    expect(customer.active).toBeTruthy();
  });

  it("should throw error when trying to activate customer without Address", () => {
    const customer = new Customer("123", "Carlos 1", undefined, true);
    expect(() => customer.activate()).toThrowError();
  });

  it("should deactivate user", () => {
    const customer = new Customer("123", "Carlos 1", new Address('123', 10, '1234', 'city'), true);
    customer.deactivate();
    expect(customer.active).toBeFalsy();
  })
});