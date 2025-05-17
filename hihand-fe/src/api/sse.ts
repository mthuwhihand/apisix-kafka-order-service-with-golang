import { CLIENT_ID } from "../constants/config";

// src/api/sse.ts
export const connectSSE = (onMessage: (msg: any) => void): EventSource => {
    const source = new EventSource(`http://localhost:9080/events/order_created?clientId=${CLIENT_ID}`);

    source.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log("Type of event.data, ", typeof (event.data))
        if (data.status_code === '200') {
            console.log('Order created successfully:', data.value);
            onMessage({
                success: true,
                message: 'Order created successfully!',
                order: data.value // Đối tượng order
            });
        } else {
            console.log('Error occurred:', data.message);
            onMessage({
                success: false,
                message: 'Error: ' + data.message,
                error: data.value // Mô tả lỗi
            });
        }
    };

    source.addEventListener("ping", (_) => {
        console.log("Received ping at", new Date().toISOString());
    });

    source.onerror = (err) => {
        console.log('SSE Error:', err);
        source.close();
    };

    return source;
};
