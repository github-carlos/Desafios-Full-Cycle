import { Sequelize } from "sequelize-typescript";
import Address from "../../domain/entity/address";
import Customer from "../../domain/entity/customer";
import EventDispatcher from "../../domain/event/@shared/event-dispatcher";
import CustomerAddressChangedEvent from "../../domain/event/customer/customer-address-changed.event";
import CustomerCreatedEvent from "../../domain/event/customer/customer-created.event";
import Log1MessageWhenCustomerIsCreatedHandler from "../../domain/event/customer/handlers/log-1-message-when-customer-is-created";
import Log1MessageWhenCustomerAddressChangedHandler from "../../domain/event/customer/handlers/log-1-messge-when-customer-address-changed";
import Log2MessageWhenCustomerIsCreatedHandler from "../../domain/event/customer/handlers/log-2-message-when-customer-is-created copy";
import CustomerModel from "../db/sequelize/models/customer.model";
import CustomerRepository from "./customer.repository";

describe("Customer repository test", () => {
  let sequelize: Sequelize;
  const eventDispatcher = new EventDispatcher();
  const customerCreatedHandler1 = new Log1MessageWhenCustomerIsCreatedHandler()
  const customerCreatedHandler2 = new Log2MessageWhenCustomerIsCreatedHandler()
  const customerCreatedHandler3 = new Log1MessageWhenCustomerAddressChangedHandler()

  beforeAll(() => {
    eventDispatcher.register('CustomerCreatedEvent', customerCreatedHandler1)
    eventDispatcher.register('CustomerCreatedEvent', customerCreatedHandler2)
    eventDispatcher.register('CustomerAddressChangedEvent', customerCreatedHandler3)
  });

  beforeEach(async () => {
    sequelize = new Sequelize({
      dialect: "sqlite",
      storage: ":memory:",
      logging: false,
      sync: { force: true },
    });

    await sequelize.addModels([CustomerModel]);
    await sequelize.sync();
  });

  afterEach(async () => {
    await sequelize.close();
  });

  it("should create a customer and dispatch CustomerCreated event", async () => {
    jest.useFakeTimers()
    .setSystemTime(new Date());

    const eventSpy = jest.spyOn(eventDispatcher, "notify")
    const handler1Spy = jest.spyOn(customerCreatedHandler1, 'handle');
    const handler2Spy = jest.spyOn(customerCreatedHandler2, 'handle');

    const customerRepository = new CustomerRepository(eventDispatcher);
    const customer = new Customer("123", "Customer 1");
    const address = new Address("Street 1", 1, "Zipcode 1", "City 1");
    customer.changeAddress(address);
    await customerRepository.create(customer);

    const customerModel = await CustomerModel.findOne({ where: { id: "123" } });

    expect(eventSpy).toHaveBeenCalledWith(new CustomerCreatedEvent(JSON.stringify(customer)));
    expect(handler1Spy).toHaveBeenCalled();
    expect(handler2Spy).toHaveBeenCalled();
    expect(customerModel.toJSON()).toStrictEqual({
      id: "123",
      name: customer.name,
      active: customer.isActive(),
      rewardPoints: customer.rewardPoints,
      street: address.street,
      number: address.number,
      zipcode: address.zip,
      city: address.city,
    });
  });

  it("should change customer address and dispatch event", async() => {
     jest.useFakeTimers()
    .setSystemTime(new Date());
    
    const eventSpy = jest.spyOn(eventDispatcher, "notify");
    const handleSpy = jest.spyOn(customerCreatedHandler3, 'handle');

    const customerRepository = new CustomerRepository(eventDispatcher);
    const customer = new Customer("123", "Customer 1");
    const address = new Address("Street 1", 1, "Zipcode 1", "City 1");
    customer.changeAddress(address);

    await customerRepository.create(customer);
    const newAddress = new Address("Another Street", 2, "Zipcode 2", "New City");
    customer.changeAddress(newAddress);
    await customerRepository.changeAddress(customer);
    const customerFromDb = await customerRepository.find(customer.id);

    expect(customerFromDb.Address).toEqual(newAddress);
    expect(handleSpy).toBeCalled();
    expect(eventSpy).toHaveBeenCalledWith(new CustomerAddressChangedEvent(customer));
  });

  it("should update a customer", async () => {
    const customerRepository = new CustomerRepository(eventDispatcher);
    const customer = new Customer("123", "Customer 1");
    const address = new Address("Street 1", 1, "Zipcode 1", "City 1");
    customer.changeAddress(address);
    await customerRepository.create(customer);

    customer.changeName("Customer 2");
    await customerRepository.update(customer);
    const customerModel = await CustomerModel.findOne({ where: { id: "123" } });

    expect(customerModel.toJSON()).toStrictEqual({
      id: "123",
      name: customer.name,
      active: customer.isActive(),
      rewardPoints: customer.rewardPoints,
      street: address.street,
      number: address.number,
      zipcode: address.zip,
      city: address.city,
    });
  });

  it("should find a customer", async () => {
    const customerRepository = new CustomerRepository(eventDispatcher);
    const customer = new Customer("123", "Customer 1");
    const address = new Address("Street 1", 1, "Zipcode 1", "City 1");
    customer.changeAddress(address);
    await customerRepository.create(customer);

    const customerResult = await customerRepository.find(customer.id);

    expect(customer).toStrictEqual(customerResult);
  });

  it("should throw an error when customer is not found", async () => {
    const customerRepository = new CustomerRepository(eventDispatcher);

    expect(async () => {
      await customerRepository.find("456ABC");
    }).rejects.toThrow("Customer not found");
  });

  it("should find all customers", async () => {
    const customerRepository = new CustomerRepository(eventDispatcher);
    const customer1 = new Customer("123", "Customer 1");
    const address1 = new Address("Street 1", 1, "Zipcode 1", "City 1");
    customer1.changeAddress(address1);
    customer1.addRewardPoints(10);
    customer1.activate();

    const customer2 = new Customer("456", "Customer 2");
    const address2 = new Address("Street 2", 2, "Zipcode 2", "City 2");
    customer2.changeAddress(address2);
    customer2.addRewardPoints(20);

    await customerRepository.create(customer1);
    await customerRepository.create(customer2);

    const customers = await customerRepository.findAll();

    expect(customers).toHaveLength(2);
    expect(customers).toContainEqual(customer1);
    expect(customers).toContainEqual(customer2);
  });
});