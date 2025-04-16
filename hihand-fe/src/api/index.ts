import http from './http';
import { Order } from '../types/Order';

export const fetchItems = async (): Promise<Order[]> => {
    const res = await http.get('/orders');
    return res.data.orders;
};

export const createOrder = async (order: Order): Promise<string> => {
    try {
        const orderData = {
            clientID: "client2",
            user_id: order.user_id,
            recipient_name: order.recipient_name,
            contact_phone: order.contact_phone,
            email: order.email,
            address: order.address,
            order_date: new Date().toISOString().slice(0, 19) + "Z",
            status: order.status,
            total: order.total,
            note: order.note,
            details: order.details,
        };

        const res = await http.post('/orders', orderData);
        return res.data?.message || 'Order created successfully!';
    } catch (error: any) {
        if (error.response && error.response.data && error.response.data.message) {
            throw new Error(error.response.data.message);
        }
        throw new Error('Đã xảy ra lỗi khi tạo đơn hàng.');
    }
};

