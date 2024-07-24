import axios from "axios";

const API_URL = "http://localhost:8000";

const api = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

export const fetchUsers = () => api.get("/users");
export const fetchUser = (userId: string) => api.get(`/users/${userId}`);
export const signUp = (data: any) => api.post("/users/signup", data);
export const login = (data: any) => api.post("/users/login", data);

export const fetchFoods = () => api.get("/foods");
export const fetchFood = (foodId: string) => api.get(`/foods/${foodId}`);
export const createFood = (data: any) => api.post("/foods", data);
export const updateFood = (foodId: string, data: any) =>
  api.patch(`/foods/${foodId}`, data);

export const fetchMenus = () => api.get("/menus");
export const fetchMenu = (menuId: string) => api.get(`/menus/${menuId}`);
export const createMenu = (data: any) => api.post("/menus", data);
export const updateMenu = (menuId: string, data: any) =>
  api.patch(`/menus/${menuId}`, data);

export const fetchTables = () => api.get("/tables");
export const fetchTable = (tableId: string) => api.get(`/tables/${tableId}`);
export const createTable = (data: any) => api.post("/tables", data);
export const updateTable = (tableId: string, data: any) =>
  api.patch(`/tables/${tableId}`, data);

export const fetchOrders = () => api.get("/orders");
export const fetchOrder = (orderId: string) => api.get(`/orders/${orderId}`);
export const createOrder = (data: any) => api.post("/orders", data);
export const updateOrder = (orderId: string, data: any) =>
  api.patch(`/orders/${orderId}`, data);

export const fetchOrderItems = () => api.get("/orderItems");
export const fetchOrderItem = (orderItemId: string) =>
  api.get(`/orderItems/${orderItemId}`);
export const fetchOrderItemsByOrder = (orderId: string) =>
  api.get(`/orderItems-order/${orderId}`);
export const createOrderItem = (data: any) => api.post("/orderItems", data);
export const updateOrderItem = (orderItemId: string, data: any) =>
  api.patch(`/orderItems/${orderItemId}`, data);

export const fetchInvoices = () => api.get("/invoices");
export const fetchInvoice = (invoiceId: string) =>
  api.get(`/invoices/${invoiceId}`);
export const createInvoice = (data: any) => api.post("/invoices", data);
export const updateInvoice = (invoiceId: string, data: any) =>
  api.patch(`/invoices/${invoiceId}`, data);

export default api;
