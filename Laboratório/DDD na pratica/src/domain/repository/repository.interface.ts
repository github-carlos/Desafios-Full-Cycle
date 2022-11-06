
export default interface Repository<T> {
  create(data: T): Promise<void>;
  find(id: string): Promise<T>;
  findAll(): Promise<Array<T>>;
}