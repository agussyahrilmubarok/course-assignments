package com.example.order.service;

import com.example.order.model.OrderDTO;

import java.util.List;

public interface OrderService {

    OrderDTO.StartOrderResponse startOrder(OrderDTO.StartOrderRequest payload);

    OrderDTO.FinishOrderRequest finishOrder(OrderDTO.FinishOrderRequest payload);

    OrderDTO.ProductDetail findById(String orderId);

    List<OrderDTO.ProductOrder> findAllByUser(String userId);
}