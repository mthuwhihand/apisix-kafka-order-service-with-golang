import http from './http';
import { Order, OrderRequest } from '../types/Order';
import { CLIENT_ID } from '../constants/config';
import { HttpStatusCode } from 'axios';

interface FetchOrdersResponse {
    orders: Order[];
    totalCount: number;
    currentPage: number;
    totalPages: number;
    hasNext: boolean;
}

export const fetchItems = async (
    page: number,
    pageSize: number
): Promise<FetchOrdersResponse> => {
    const skip = (page - 1) * pageSize;
    const limit = pageSize;

    const res = await http.get('/orders', {
        params: { skip, limit }
    });

    const data = res.data.data;

    return {
        orders: data.data,
        totalCount: data.total_records,
        currentPage: data.current_page,
        totalPages: data.total_pages,
        hasNext: data.has_next,
    };
};

export const createOrder = async (order: OrderRequest): Promise<string> => {
    try {
        const orderData = {
            clientID: CLIENT_ID,
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

export const deleteOrder = async (orderId: string): Promise<string> => {
    const res = await http.delete(`/orders/${orderId}`);
    if (res.status != HttpStatusCode.Ok) {
        throw new Error(res.data?.message || 'Xóa đơn hàng thất bại!');
    }
    return res.data?.message || 'Order created successfully!';
};
