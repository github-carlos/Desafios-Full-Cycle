import Order from "../../domain/entity/order";
import OrderItem from "../../domain/entity/order_item";
import { OrderRepositoryInterface } from "../../domain/repository/order.repository";
import { OrderItemModel, OrderModel } from "../db/sequelize/models/order.model";

export default class OrderRepository implements OrderRepositoryInterface {
  async find(id: string): Promise<Order> {
    try {
      const orderModel = await OrderModel.findOne({
        where: { id },
        include: [{ model: OrderItemModel }],
        rejectOnEmpty: true,
      });
      return new Order(
        orderModel.id,
        orderModel.customer_id,
        orderModel.items.map(
          (item) =>
            new OrderItem(
              item.id,
              item.name,
              item.price,
              item.product_id,
              item.quantity
            )
        )
      );
    } catch (err) {
      throw new Error("Order not found");
    }
  }
  async create(entity: Order): Promise<void> {
    await OrderModel.create(
      {
        id: entity.id,
        customer_id: entity.customerId,
        total: entity.total(),
        items: entity.items.map((item) => ({
          id: item.id,
          name: item.name,
          price: item.price,
          product_id: item.productId,
          quantity: item.quantity,
        })),
      },
      { include: [{ model: OrderItemModel }] }
    );
  }

  async update(entity: Order): Promise<void> {
    try {
      const order = await OrderModel.findOne({
        where: { id: entity.id },
        include: ["items"],
      });

      for (const item of entity.items) {
        if (!order.items.map((i) => i.id).includes(item.id)) {
          await (order as any).createItem({
            id: item.id,
            name: item.name,
            price: item.price,
            product_id: item.productId,
            quantity: item.quantity,
          });
        }
      }

      await OrderModel.update(
        {
          customer_id: entity.customerId,
          total: entity.total(),
        },
        {
          where: {
            id: entity.id,
          },
        }
      );
    } catch (err) {
      throw new Error("Update Order Error");
    }
  }

  async findAll(): Promise<Order[]> {
    const orderModels = await OrderModel.findAll({ include: [OrderItemModel] });

    const orders = orderModels.map((orderModel) => {
      const order = new Order(
        orderModel.id,
        orderModel.customer_id,
        orderModel.items.map(
          (item) =>
            new OrderItem(
              item.id,
              item.name,
              item.price,
              item.product_id,
              item.quantity
            )
        )
      );
      return order;
    });

    return orders;
  }
}
