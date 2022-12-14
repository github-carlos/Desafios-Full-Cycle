import { Sequelize } from "sequelize-typescript";
import Address from "../../domain/entity/address";
import Customer from "../../domain/entity/customer";
import Order from "../../domain/entity/order";
import OrderItem from "../../domain/entity/order_item";
import Product from "../../domain/entity/product";
import EventDispatcher from "../../domain/event/@shared/event-dispatcher";
import CustomerModel from "../db/sequelize/models/customer.model";
import { OrderItemModel, OrderModel } from "../db/sequelize/models/order.model";
import ProductModel from "../db/sequelize/models/product";
import CustomerRepository from "./customer.repository";
import OrderRepository from "./order.repository";
import ProductRepository from "./product.repository";

async function makeOrder(orderItems: Array<OrderItem>): Promise<Order> {
  const eventDispatcher = new EventDispatcher();
  const customerRepository = new CustomerRepository(eventDispatcher);
  const customer = new Customer("123", "Carlos");
  const address = new Address("Street 1", 1, "Zipcode", "City");
  customer.changeAddress(address);
  await customerRepository.create(customer);

  const productRepository = new ProductRepository();
  const product = new Product("123", "Product Name", 10);
  await productRepository.create(product);

  return new Order("123", customer.id, orderItems);
}

describe("Order repository test", () => {
  let sequelize: Sequelize;

  beforeEach(async () => {
    sequelize = new Sequelize({
      dialect: "sqlite",
      storage: ":memory:",
      logging: false,
      sync: { force: true },
    });

    await sequelize.addModels([
      CustomerModel,
      OrderModel,
      OrderItemModel,
      ProductModel,
    ]);
    await sequelize.sync();
  });

  afterEach(async () => {
    await sequelize.drop();
    await sequelize.close();
  });

  it("should create a new order", async () => {
    const orderItem = new OrderItem('123', 'item 1', 10, '123', 10);
    const order = await makeOrder([orderItem]);
    const orderRepository = new OrderRepository();
    await orderRepository.create(order);

    const orderModel = await OrderModel.findOne({
      where: { id: order.id },
      include: ["items"],
    });

    expect(orderModel.toJSON()).toStrictEqual({
      id: "123",
      customer_id: "123",
      total: order.total(),
      items: order.items.map((orderItem) => ({
        id: orderItem.id,
        name: orderItem.name,
        price: orderItem.price,
        quantity: orderItem.quantity,
        product_id: "123",
        order_id: "123",
      })),
    });
  });

  it("should find an order", async () => {
    const orderItem = new OrderItem('123', 'item 1', 10, '123', 10);
    const order = await makeOrder([orderItem]);
    const orderRepository = new OrderRepository();

    await orderRepository.create(order);

    const orderResult = await orderRepository.find(order.id);

    expect(orderResult).toStrictEqual(order);
  });

  it("should throw error when order not found", async () => {
    const orderRepository = new OrderRepository();

    expect(async () => {
      await orderRepository.find("123invalid");
    }).rejects.toThrow("Order not found");
  });

  it("should update an order", async () => {
    const orderItem = new OrderItem('123', 'item 1', 10, '123', 10);
    const order = await makeOrder([orderItem]);
    const orderRepository = new OrderRepository();
    await orderRepository.create(order);

    // creating new orderItem
    const productRepository = new ProductRepository();
    const product = new Product("456", "Product Name", 10);
    await productRepository.create(product);   
    const newOrderITem = new OrderItem("456", "item 2", 20, product.id, 2);

    order.addItem(newOrderITem);
    await orderRepository.update(order);

    const updatedOrder = await OrderModel.findOne({
      where: { id: order.id },
      include: ["items"],
    });

    expect(updatedOrder.toJSON().items).toStrictEqual([
      {id: '123', name: 'item 1', price: 10, product_id: '123', quantity: 10, order_id: order.id},
      {id: '456', name: 'item 2', price: 20, product_id: '456', quantity: 2, order_id: order.id},
    ]);
  });

  it("should find all orders", async() => {
  const eventDispatcher = new EventDispatcher();
  const customerRepository = new CustomerRepository(eventDispatcher);
  const customer = new Customer("123", "Carlos");
  const address = new Address("Street 1", 1, "Zipcode", "City");
  customer.changeAddress(address);
  await customerRepository.create(customer);

  const productRepository = new ProductRepository();
  const product = new Product("123", "Product Name", 10);
  await productRepository.create(product);

  const orderItem1 = new OrderItem('123', 'item 1', 10, '123', 10);
  const orderItem2 = new OrderItem('456', 'item 2', 10, '123', 10);

  const order1 = new Order("123", customer.id, [orderItem1]);
  const order2 = new Order("345", customer.id, [orderItem2]);

  const orderRepository = new OrderRepository();
  await orderRepository.create(order1);
  await orderRepository.create(order2);


  const orders = await orderRepository.findAll();

  expect(orders.length).toBe(2);
  });
});
